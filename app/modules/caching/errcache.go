package caching

import (
    "github.com/robfig/revel"
    "github.com/robfig/revel/cache"
    "github.com/mawenbao/wishome/app"
)

type SigninErrorCache struct {
    Name string
    ErrorCount int
    CaptchaRequired bool
    Banned bool
}

func GetSigninErrorKeyName(name string) string {
    return name + "." + app.CACHE_SIGNIN_ERROR
}

func GetSigninError(name string) *SigninErrorCache {
    sess := new(SigninErrorCache)
    err := cache.Get(GetSigninErrorKeyName(name), sess)
    if nil != err {
        revel.INFO.Printf("error get signin error from cache for %s: %s", name, err)
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
func NewSigninError(name string) *SigninErrorCache {
    sess := GetSigninError(name)
    if nil == sess {
        // save new temp session in cache
        sess = &SigninErrorCache{
            Name: name,
            ErrorCount: 1,
            Banned: false,
        }
    } else {
        // update cache
        sess.ErrorCount += 1
        if sess.ErrorCount >= app.MyGlobal.SigninUseCaptchaErrorNum {
            sess.CaptchaRequired = true
        }
        if sess.ErrorCount >= app.MyGlobal.SigninErrLimit {
            sess.Banned = true
        }
    }

    go cache.Set(GetSigninErrorKeyName(name), *sess, app.MyGlobal.SigninErrCacheLife)
    return sess
}

