package controllers

import (
    "time"
    "fmt"
    "strings"
    "github.com/robfig/revel"
    "github.com/robfig/revel/cache"
    "github.com/mawenbao/wishome/app"
    "github.com/mawenbao/wishome/app/models"
    "github.com/mawenbao/wishome/app/modules/common"
    "github.com/mawenbao/wishome/app/modules/database"
    "github.com/mawenbao/wishome/app/modules/validators"
    "github.com/mawenbao/wishome/app/routes"
)

type User struct {
    *revel.Controller
}

func (c User) setUserSession(u models.User) {
    aesKey := common.NewRawRandom(32) // aes-256
    c.Session[app.STR_KEY] = aesKey

    c.Session[app.STR_NAME] = u.Name
    c.Session[app.STR_EMAIL] = common.EncodeToHexString(common.AesEncrypt([]byte(u.Email), []byte(aesKey)))
    /*
    c.Session[app.STR_KEY] = "abc"
    c.Session[app.STR_EMAIL] = u.Email
    */

    c.RenderArgs[app.STR_USER] = &u
}

func (c User) getUserSession() (u *models.User) {
    if "" == c.Session[app.STR_KEY] || "" == c.Session[app.STR_NAME] || "" == c.Session[app.STR_EMAIL] {
        return nil
    }
    u = new(models.User)

    u.Name = c.Session[app.STR_NAME]
    key := []byte(c.Session[app.STR_KEY])
    u.Email = string(common.AesDecrypt(common.DecodeHexString(c.Session[app.STR_EMAIL]), key))
    /*
    u.Email = c.Session[app.STR_EMAIL]
    */

    if "" == u.Name || "" == u.Email {
        return nil
    }
    return u
}

func (c User) emptyUserSession() {
    c.RenderArgs[app.STR_USER] = nil

    for k := range c.Session {
        if app.STR_NAME != k {
            delete(c.Session, k)
        }
    }
}

func (c User) connected() *models.User {
    if nil != c.RenderArgs[app.STR_USER] {
        return c.RenderArgs[app.STR_USER].(*models.User)
    }

    // check session, with no ID and PassSalt field
    return c.getUserSession()
}

func (c User) DoSignin(name, password string) revel.Result {
    if nil !=  c.connected() {
        c.Flash.Success(c.Message("user.signout.succeeded"))
        c.emptyUserSession()
    }

    dbmgr := database.NewDbManager()
    if nil == dbmgr {
        panic("db error: failed to init")
    }
    defer dbmgr.Close()

    _, user := validators.ValidateSignin(c.Validation, dbmgr.DbMap, name, password)
    if c.Validation.HasErrors() {
        c.Validation.Keep()
        c.FlashParams()
        return c.Redirect(routes.User.Signin())
    }

    c.setUserSession(*user)
    c.Flash.Success(c.Message("user.greeting.old", name))
    return c.Redirect(routes.User.Home())
}

func (c User) DoSignup(name, email, password string) revel.Result {
    // check if user has signed in already, if so, sign him out
    if nil != c.connected() {
        c.emptyUserSession()
    }

    user := models.User{
        Name: name,
        Email: email,
        Password: password,
    }
    // init db connection
    dbmgr := database.NewDbManager()
    if nil == dbmgr {
        panic("db error: failed to init")
    }
    defer dbmgr.Close()

    // validate input 
    validators.ValidateSignup(c.Validation, dbmgr.DbMap, user.Name, user.Email, user.Password)
    if c.Validation.HasErrors() {
        c.Validation.Keep()
        c.FlashParams()
        return c.Redirect(routes.User.Signup())
    }

    // generate password salt and encrypt user password
    user.EncryptPass()

    if !database.SaveUser(dbmgr.DbMap, user) {
        c.Flash.Error(c.Message("user.save.error.db"))
        return c.Redirect(routes.User.Signup())
    } else {
        // set session
        c.setUserSession(user)
        c.Flash.Success(c.Message("user.greeting.new", user.Name))
        return c.Redirect(routes.User.Home())
    }
}

func (c User) DoSignout() revel.Result {
    c.emptyUserSession()
    c.Flash.Success(c.Message("user.signout.succeeded"))
    return c.Redirect(routes.User.Signin())
}

func (c User) Signin() revel.Result {
    // check if user has signed in
    if uc := c.connected(); nil != uc {
        //c.Flash.Success(c.Message("user.greeting.old", uc.Name))
        return c.Redirect(routes.User.Home())
    }
    return c.Render()
}

func (c User) Signup() revel.Result {
    if nil != c.connected() {
        c.Flash.Success(c.Message("user.signout.succeeded"))
        c.emptyUserSession()
    }

    moreNavbarLinks := []models.NavbarLink{
    }
    return c.Render(moreNavbarLinks)
}

func (c User) ResetPass() revel.Result {
    return c.Render()
}

// get cache key of reset pass key
func getResetpassKeyName(name string) string {
    return app.CACHE_RESETPASS + name
}

// check name and key in cache
func checkCachedResetpassKey(name, key string) bool {
    var cachedKey string
    err := cache.Get(getResetpassKeyName(name), &cachedKey)
    if nil != err {
        revel.ERROR.Printf("failed to get reset password key from cache for user %s key %s: %s", name, key, err)
        return false
    }
    if key != cachedKey {
        revel.ERROR.Printf("got mismatched reset password key, expect %s, got %s", cachedKey, key)
        return false
    }
    return true
}

// validate name, email and send an email to user with a random key
// which is valid in half an hour
func (c User) PreResetPass(name, email string) revel.Result {
    // init db connection
    dbmgr := database.NewDbManager()
    if nil == dbmgr {
        panic("db error: failed to init")
    }
    defer dbmgr.Close()

    // validate name and email
    validators.ValidateResetPassNameEmail(c.Validation, dbmgr.DbMap, name, email)
    if c.Validation.HasErrors() {
        c.Validation.Keep()
        c.FlashParams()
        return c.Redirect(routes.User.ResetPass())
    }

    // generate a random key for user, expires in 30 minutes
    keyLen := app.MyGlobal[app.CONFIG_RESETPASS_KEY_LEN].(int)
    keyLife := app.MyGlobal[app.CONFIG_RESETPASS_KEY_LIFE].(time.Duration)
    key := common.NewReadableRandom(keyLen)
    revel.INFO.Printf("generate new reset password key %s for name %s", key, name)
    go cache.Set(getResetpassKeyName(name), key, keyLife)

    resetPassUrl, found := revel.Config.String(app.CONFIG_APP_URL)
    if !found {
        revel.ERROR.Printf("%s not set in app.conf", app.CONFIG_APP_URL)
        panic("internal error")
    }
    resetPassUrl = fmt.Sprintf(
        "%s/user/postresetpass?name=%s&key=%s",
        strings.TrimRight(resetPassUrl, "/"),
        name,
        key,
    )

    // mail the url link
    revel.ERROR.Printf(resetPassUrl)

    c.Flash.Success(c.Message("user.resetpass.mail.succeeded", email))
    return c.Redirect(routes.User.ResetPass())
}

func (c User) PostResetPass(name, key string) revel.Result {
    revel.INFO.Printf("post got reset password key %s for name %s", key, name)
    if !checkCachedResetpassKey(name, key) {
        c.Flash.Error(c.Message("user.resetpass.key.error"))
        return c.Redirect(routes.User.ResetPass())
    }
    return c.Render(name, key)
}

func (c User) DoResetPass(name, password, key string) revel.Result {
    revel.INFO.Printf("do got key %s for name %s", key, name)
    // check key first
    if !checkCachedResetpassKey(name, key) {
        c.Flash.Error(c.Message("user.resetpass.key.error"))
        return c.Redirect(routes.User.ResetPass())
    }

    // check password
    validators.ValidatePassword(c.Validation, password)
    if c.Validation.HasErrors() {
        c.Validation.Keep()
        c.FlashParams()
        return c.Redirect(routes.User.PostResetPass(name, key))
    }

    // init db connection
    dbmgr := database.NewDbManager()
    if nil == dbmgr {
        panic("db error: failed to init")
    }
    defer dbmgr.Close()

    u := database.FindUserByName(dbmgr.DbMap, name)
    if nil == u || !u.IsSecured() {
        c.Flash.Error(c.Message("common.user.none", name))
        return c.Redirect(routes.User.Signin())
    }

    // change password and password salt
    u.Password = password
    u.EncryptPass()
    if !database.UpdateUser(dbmgr.DbMap, *u) {
        revel.ERROR.Printf("failed to update user password for %s", u)
        c.Flash.Error(c.Message("user.resetpass.failed"))
        return c.Redirect(routes.User.ResetPass())
    }

    revel.INFO.Printf("successfully updated user password for %s", name)
    c.Flash.Success(c.Message("user.resetpass.succeeded"))
    return c.Redirect(routes.User.Signin())
}

func (c User) Home() revel.Result {
    // check if user has signed in
    if nil == c.connected() {
        c.Flash.Error(c.Message("user.signin.succeeded"))
    }

    moreNavbarLinks := []models.NavbarLink{
        models.NavbarLink{"/user/dosignout", "sign out", "sign out", false},
    }
    return c.Render(moreNavbarLinks)
}

