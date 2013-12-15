package caching

import (
    "time"
    "github.com/robfig/revel"
    "github.com/robfig/revel/cache"
    "github.com/mawenbao/wishome/app"
    "github.com/mawenbao/wishome/app/models"
    "github.com/mawenbao/wishome/app/modules/database"
    "github.com/mawenbao/wishome/app/modules/common"
)

// get cache key of reset pass key
func GetResetPassKeyName(name string) string {
    return name + "." + app.CACHE_RESETPASS
}

// get cache key of signup confirmation key 
func GetSignupConfirmKeyName(name string) string {
    return name + "." + app.CACHE_SIGNUP_CONFIRM
}

// check name and resetpass key in cache
func CheckCachedResetPassKey(name, key string) bool {
    var cachedKey string
    err := cache.Get(GetResetPassKeyName(name), &cachedKey)
    if nil != err {
        revel.ERROR.Printf("failed to get reset password key from cache for user %s key %s: %s", name, key, err)
        return false
    }
    if key != cachedKey {
        revel.ERROR.Printf("got mismatched reset password key, expect %s, got %s", cachedKey, key)
        return false
    }

    // remove key at last
    Remove(GetResetPassKeyName(name))
    return true
}

// check name and confirmation key in cache
func CheckCachedSignupConfirmKey(name, key string) bool {
    var cachedKey string
    err := cache.Get(GetSignupConfirmKeyName(name), &cachedKey)
    if nil != err {
        revel.ERROR.Printf("failed to get signup confirmation key from cache for user %s key %s: %s", name, key, err)
        return false
    }
    if key != cachedKey {
        revel.ERROR.Printf("got mismatched signup confirmation key, expect %s, got %s", cachedKey, key)
        return false
    }

    // Remove key at last
    Remove(GetSignupConfirmKeyName(name))
    return true
}

// insert a random key into cache
// return the newly generated key
func NewReadableKey(keyName string, keyLen int, expires time.Duration) string {
    key := common.NewReadableRandom(keyLen)
    go cache.Set(keyName, key, expires)
    return key
}

func NewResetPassKey(name string) string {
    keyLen := app.MyGlobal[app.CONFIG_RESETPASS_KEY_LEN].(int)
    keyLife := app.MyGlobal[app.CONFIG_RESETPASS_KEY_LIFE].(time.Duration)
    return NewReadableKey(GetResetPassKeyName(name), keyLen, keyLife)
}

func NewSignupConfirmKey(name string) string {
    keyLen := app.MyGlobal[app.CONFIG_SIGNUP_KEY_LEN].(int)
    keyLife := app.MyGlobal[app.CONFIG_SIGNUP_KEY_LIFE].(time.Duration)
    return NewReadableKey(GetSignupConfirmKeyName(name), keyLen, keyLife)
}

// remove cache
func Remove(keyName string) bool {
    err := cache.Delete(keyName)
    if nil != err {
        revel.INFO.Printf("delete cache missed %s: %s", keyName, err)
        return false
    }
    return true
}

// set user in cache
func SetUser(u *models.User) {
    go cache.Set(u.Name, *u, app.MyGlobal[app.CONFIG_SESSION_LIFE].(time.Duration))
}

// get user from cache
func GetUser(name string) *models.User {
    var u models.User
    err := cache.Get(name, &u)
    if nil != err {
        revel.INFO.Printf("get cache for user %s missed: %s", name, err)
        return nil
    }
    return &u
}

// remove the old user from cache and get a new one from database
// then save it in cache
func ReloadUser(name string) *models.User {
    Remove(name)
    u := database.FindUserByName(name)
    if nil != u && u.IsSecured() {
        SetUser(u)
        return u
    }
    return nil
}

