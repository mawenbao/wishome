package session

import (
    "time"
    "github.com/robfig/revel"
    "github.com/mawenbao/wishome/app"
    "github.com/mawenbao/wishome/app/models"
    "github.com/mawenbao/wishome/app/modules/common"
)

type UserSession struct {
    AesKey, // aes-256
    UserName,
    Password,
    LastUser,
    Expire string
    Encrypted bool
}

func NewUserSession(u *models.User) *UserSession {
    if nil == u || !u.IsSecured() {
        return nil
    }
    return &UserSession {
        AesKey: common.NewRawRandom(32), // aes-256
        UserName: u.Name,
        Password: u.Password,
        LastUser: u.Name,
        Expire: time.Now().Format(app.DEFAULT_TIME_FORMAT),
        Encrypted: false,
    }
}

func (s *UserSession) IsExpired() bool {
    expire := []byte(s.Expire)
    if s.Encrypted {
        expire, _ = common.DecodeHexString(s.Expire)
        if nil == expire {
            return true
        }
    }

    if expireTime, err:= time.Parse(app.DEFAULT_TIME_FORMAT, string(expire)); nil != err {
        revel.ERROR.Printf("failed to parse expire time %s with format %s", expire, app.DEFAULT_TIME_FORMAT)
        return true
    } else {
        if !expireTime.After(time.Now()) {
            revel.TRACE.Printf("session expired")
            return true
        }
    }
    return false
}

func (s *UserSession) Encrypt() bool {
    if "" == s.AesKey || "" == s.Expire {
        revel.ERROR.Printf("failed to encrypt session, key or expire time is empty")
        return false
    }

    // encrypt name and password with aes
    cipherName, ok := common.AesEncrypt([]byte(s.UserName), []byte(s.AesKey))
    if !ok {
        revel.ERROR.Printf("failed to encrypt session user name %s", s.UserName)
        return false
    }
    cipherPass, ok := common.AesEncrypt([]byte(s.Password), []byte(s.AesKey))
    if !ok {
        revel.ERROR.Printf("failed to encrypt password for session user %s", s.UserName)
        return false
    }

    // encode all fields with hex
    s.AesKey = common.EncodeToHexString([]byte(s.AesKey))
    s.UserName = common.EncodeToHexString(cipherName)
    s.Password = common.EncodeToHexString(cipherPass)
    s.LastUser = common.EncodeToHexString([]byte(s.LastUser))
    s.Expire = common.EncodeToHexString([]byte(s.Expire))
    s.Encrypted = true
    return true
}

func (s *UserSession) Decrypt() bool {
    if !s.Encrypted {
        revel.ERROR.Printf("cannot decrypt an plain user session for %s", s.UserName)
        return false
    }

    // hex decode AesKey
    key, _ := common.DecodeHexString(s.AesKey)
    if nil == key {
        revel.ERROR.Printf("failed to decode hex key string %s", s.AesKey)
        return false
    }
    // hex decode username and password
    nameDeHex, _ := common.DecodeHexString(s.UserName)
    passDeHex, _ := common.DecodeHexString(s.Password)
    if nil == nameDeHex || nil == passDeHex {
        revel.ERROR.Printf("failed to decode hex name %s or hex password %s", s.UserName, s.Password)
        return false
    }
    // aes decrypt username and password
    nameSl, _ := common.AesDecrypt(nameDeHex, key)
    passSl, _ := common.AesDecrypt(passDeHex, key)
    if nil == nameSl || nil == passSl {
        revel.ERROR.Printf("failed to decode name %s or password %s", nameDeHex, passDeHex)
        return false
    }
    // hex decode expire
    expire, _ := common.DecodeHexString(s.Expire)
    if nil == expire {
        revel.ERROR.Printf("failed to decode hex expire string %s", s.Expire)
        return false
    }
    // hex decode lastuser
    lastuser, _ := common.DecodeHexString(s.LastUser)
    if nil == lastuser {
        revel.ERROR.Printf("failed to decode hex lastuser string %s", s.LastUser)
        return false
    }

    s.AesKey = string(key)
    s.UserName = string(nameSl)
    s.Password = string(passSl)
    s.LastUser = string(lastuser)
    s.Expire = string(expire)
    s.Encrypted = false
    return true
}

