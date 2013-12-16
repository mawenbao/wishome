package caching

import (
    "time"
    "github.com/robfig/revel"
    "github.com/robfig/revel/cache"
    "github.com/mawenbao/wishome/app"
)

var (
    SIGNIN_SESSION_LIFE = time.Duration(1 * time.Hour)
    SIGNIN_CAPCHAR_LIMIT = 5
    SIGNIN_ERROR_LIMIT = 50
)

type SigninErrorSession struct {
    Name string
    ErrorCount int
    CaptchaRequired bool
    Banned bool
}

func GetSigninErrorKeyName(name string) string {
    return name + "." + app.CACHE_SIGNIN_ERROR
}

func GetSigninError(name string) *SigninErrorSession {
    sess := new(SigninErrorSession)
    err := cache.Get(GetSigninErrorKeyName(name), sess)
    if nil != err {
        revel.INFO.Printf("error get temp session from cache: %s", err)
        return nil
    }
    return sess
}

func IsSigninCaptchaRequired(name string) bool {
    sess := GetSigninError(name)
    if nil != sess {
        return sess.CaptchaRequired
    }
    return false
}

func IsSigninBanned(name string) bool {
    sess := GetSigninError(name)
    if nil != sess {
        return sess.Banned
    }
    return false
}

// new user signin error
func NewSigninError(name string) *SigninErrorSession {
    sess := GetSigninError(name)
    if nil == sess {
        // save new temp session in cache
        sess = &SigninErrorSession{
            Name: name,
            ErrorCount: 1,
            Banned: false,
        }
    } else {
        // update cache
        sess.ErrorCount += 1
        if sess.ErrorCount >= SIGNIN_CAPCHAR_LIMIT {
            sess.CaptchaRequired = true
        }
        if sess.ErrorCount >= SIGNIN_ERROR_LIMIT {
            sess.Banned = true
        }
    }

    go cache.Set(GetSigninErrorKeyName(name), *sess, SIGNIN_SESSION_LIFE)
    return sess
}

