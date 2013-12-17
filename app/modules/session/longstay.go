package session

import (
    "github.com/robfig/revel"
    "github.com/mawenbao/wishome/app"
    "github.com/mawenbao/wishome/app/modules/common"
)

type LongStaySession struct {
    LastUser string
    Encrypted bool
}

func (s *LongStaySession) Encrypt() bool {
    s.LastUser = common.EncodeToHexString([]byte(s.LastUser))
    s.Encrypted = true
    return true
}

func (s *LongStaySession) Decrypt() bool {
    lastUserSL, _ := common.DecodeHexString(s.LastUser)
    if nil == lastUserSL {
        revel.ERROR.Printf("failed to decode last user hex string %s", s.LastUser)
        return false
    }

    s.LastUser = string(lastUserSL)
    s.Encrypted = false
    return true
}

func (s *LongStaySession) IsEncrypted() bool {
    return s.Encrypted
}

func (s *LongStaySession) IsExpired() bool {
    return false
}

func (s *LongStaySession) Load(session map[string]string) bool {
    s.LastUser = session[app.STR_LASTUSER]
    return s.Decrypt()
}

func (s *LongStaySession) Save(session map[string]string) bool {
    if !s.IsEncrypted() && !s.IsEncrypted() {
        revel.ERROR.Println("failed to encrypt long stay session in cookie")
        return false
    }

    session[app.STR_LASTUSER] = s.LastUser
    return true
}

