package controllers

import (
    "time"
    "fmt"
    "github.com/robfig/revel"
    "github.com/mawenbao/wishome/app"
    "github.com/mawenbao/wishome/app/models"
    "github.com/mawenbao/wishome/app/modules/common"
    "github.com/mawenbao/wishome/app/modules/database"
    "github.com/mawenbao/wishome/app/modules/validators"
    "github.com/mawenbao/wishome/app/modules/ext/caching"
    "github.com/mawenbao/wishome/app/modules/ext/mail"
    "github.com/mawenbao/wishome/app/routes"
)

type User struct {
    *revel.Controller
}

// return true if user session has expired
func (c User) isSessionExpired() bool {
    // check if session expires
    if expiresHex, ok := c.Session[app.STR_EXPIRE]; !ok {
        return true
    } else {
        if expires, ok := common.DecodeHexString(expiresHex); !ok {
            revel.ERROR.Printf("failed to decode hex expire time %s", expiresHex)
            return true
        } else {
            if expireTime, err:= time.Parse(app.DEFAULT_TIME_FORMAT, string(expires)); nil != err {
                revel.ERROR.Printf("failed to parse expire time %s with format %s", expires, app.DEFAULT_TIME_FORMAT)
                return true
            } else {
                if !expireTime.After(time.Now()) {
                    revel.TRACE.Printf("session expired")
                    return true
                }
            }
        }
    }

    return false
}

func (c User) setUserSession(u *models.User) bool {
    if nil == u || !u.IsSecured() {
        return false
    }

    expireTime := time.Now().Add(app.MyGlobal[app.CONFIG_SESSION_LIFE].(time.Duration))
    aesKey := []byte(common.NewRawRandom(32)) // aes-256
    cipherName, ok := common.AesEncrypt([]byte(u.Name), aesKey)
    if !ok {
        revel.ERROR.Printf("failed to encrypt user name %s", u.Name)
        return false
    }
    cipherPass, ok := common.AesEncrypt([]byte(u.Password), aesKey)
    if !ok {
        revel.ERROR.Printf("failed to encrypt password for user %s", u.Name)
        return false
    }

    // save user name in cookie
    c.Session[app.STR_KEY] = common.EncodeToHexString(aesKey)
    c.Session[app.STR_NAME] = common.EncodeToHexString(cipherName)
    c.Session[app.STR_PASSWORD] = common.EncodeToHexString(cipherPass)
    c.Session[app.STR_EXPIRE] = common.EncodeToHexString([]byte(expireTime.Format(app.DEFAULT_TIME_FORMAT)))

    // save entire user in server cache
    caching.SetUser(u)

    return true
}

func (c User) getUserSession() (u *models.User) {
    if "" == c.Session[app.STR_KEY] || "" == c.Session[app.STR_NAME] {
        return nil
    }

    // check if session expires
    if c.isSessionExpired() {
        return nil
    }

    key, ok := common.DecodeHexString(c.Session[app.STR_KEY])
    if !ok {
        revel.ERROR.Printf("failed to decode hex key string %s", c.Session[app.STR_KEY])
        return nil
    }
    nameDeHex, _ := common.DecodeHexString(c.Session[app.STR_NAME])
    passDeHex, _ := common.DecodeHexString(c.Session[app.STR_PASSWORD])
    if nil == nameDeHex || nil == passDeHex {
        revel.ERROR.Printf("failed to decode hex name %s or hex password %s", c.Session[app.STR_NAME], c.Session[app.STR_PASSWORD])
        return nil
    }

    nameSl, _ := common.AesDecrypt(nameDeHex, key)
    passSl, _ := common.AesDecrypt(passDeHex, key)
    if nil == nameSl || nil == passSl {
        revel.ERROR.Printf("failed to decode name %s or password %s", nameDeHex, passDeHex)
        return nil
    }

    name := string(nameSl)
    pass := string(passSl)

    // get user from cache
    u = caching.GetUser(name)
    if nil == u || !u.IsSecured() {
        // cache missed, session valid, try to reload user cache from database
        // aquire a new db connection
        dbmgr := database.NewDbManager()
        if nil == dbmgr {
            panic("db error: failed to init")
        }
        defer dbmgr.Close()
        u = caching.ReloadUser(dbmgr, name)
        if nil == u || !u.IsSecured() {
            revel.ERROR.Printf("failed to reload user cache for %s", name)
            return nil
        }
    }

    if name != u.Name || pass != u.Password {
        return nil
    }

    return u
}

func (c User) emptyUserSession(name string) {
    // clear cookie
    for k := range c.Session {
        delete(c.Session, k)
    }
    // clear server cache entry
    caching.Remove(name)
}

func (c User) DoSignin(name, password string) revel.Result {
    if uc :=  c.getUserSession(); nil != uc {
        c.emptyUserSession(name)
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

    c.setUserSession(user)
    c.Flash.Success(c.Message("user.greeting.old", name))
    revel.INFO.Printf("user %s signed in", name)
    return c.Redirect(routes.User.Home())
}

func (c User) DoSignup(name, email, password string) revel.Result {
    // check if user has signed in already, if so, sign him out
    if nil != c.getUserSession() {
        c.emptyUserSession(name)
    }

    user := models.User{
        Name: name,
        Email: email,
        EmailVerified: false,
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

    // send user a confirmation mail
    go c.sendConfirmEmail(name, email)

    // save user in database
    if !database.SaveUser(dbmgr.DbMap, user) {
        c.Flash.Error(c.Message("user.save.error.db"))
        return c.Redirect(routes.User.Signup())
    } else {
        revel.INFO.Printf("new user %s signed up", name)
        c.Flash.Success(c.Message("misc.signup.notice.verify.email"))
        c.Flash.Data[app.STR_NAME] = name // pass name to signin page
        return c.Redirect(routes.User.Signin())
    }
}

func (c User) DoSignout() revel.Result {
    if uc := c.getUserSession(); nil != uc && uc.IsSecured() {
        revel.INFO.Printf("user %s signed out", uc.Name)
        c.emptyUserSession(uc.Name)
    }
    c.Flash.Success(c.Message("user.signout.succeeded"))
    return c.Redirect(routes.User.Signin())
}

func (c User) Signin() revel.Result {
    // check if user has signed in
    if uc := c.getUserSession(); nil != uc && uc.IsSecured() {
        return c.Redirect(routes.User.Home())
    }

    // get user name, may be empty
    name := c.Flash.Data[app.STR_NAME]
    return c.Render(name)
}

func (c User) Signup() revel.Result {
    if uc := c.getUserSession(); nil != uc && uc.IsSecured() {
        c.Flash.Success(c.Message("user.signout.succeeded"))
        c.emptyUserSession(uc.Name)
    }

    moreNavbarLinks := []models.NavbarLink{
    }
    return c.Render(moreNavbarLinks)
}

func (c User) ResetPass() revel.Result {
    return c.Render()
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

    // send resetpass mail
    go c.sendResetPassEmail(name, email)

    c.Flash.Success(c.Message("user.resetpass.mail.succeeded", email))
    return c.Redirect(routes.User.ResetPass())
}

func (c User) PostResetPass(name, key string) revel.Result {
    revel.TRACE.Printf("postresetpass got reset password key %s for name %s", key, name)
    return c.Render(name, key)
}

func (c User) DoResetPass(name, password, key string) revel.Result {
    revel.TRACE.Printf("doresetpass got key %s for name %s", key, name)
    // check key first
    if !caching.CheckCachedResetPassKey(name, key) {
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
    // reset user cache
    caching.SetUser(u)

    revel.INFO.Printf("successfully updated user password for %s", name)
    c.Flash.Success(c.Message("user.resetpass.succeeded"))
    return c.Redirect(routes.User.Signin())
}

func (c User) sendResetPassEmail(name, email string) bool {
    resetPassUrl := fmt.Sprintf(
        "%s/user/postresetpass?name=%s&key=%s",
        app.MyGlobal[app.CONFIG_APP_URL],
        name,
        caching.NewResetPassKey(name),
    )

    revel.INFO.Printf("try to send a reset password mail to %s", email)
    return mail.SendMail(
        email,
        c.Message("user.resetpass.mail.subject"),
        []byte(c.Message("user.resetpass.mail.content", resetPassUrl, app.MyGlobal[app.CONFIG_RESETPASS_KEY_LIFE].(time.Duration).Minutes())),
    )
}

func (c User) sendConfirmEmail(name, email string) bool {
    cfmURL := fmt.Sprintf(
        "%s/user/doverifyemail?name=%s&key=%s",
        app.MyGlobal[app.CONFIG_APP_URL],
        name,
        caching.NewSignupConfirmKey(name),
    )
    revel.INFO.Printf("try to send a signup confirmation email to %s", email)
    return mail.SendMail(
        email,
        c.Message("user.signup.mail.subject"),
        []byte(c.Message("user.signup.mail.content", cfmURL, app.MyGlobal[app.CONFIG_SIGNUP_KEY_LIFE].(time.Duration).Minutes())),
    )
}

// user need to sign in first in order to resend confirmation email
func (c User) ResendConfirmEmail() revel.Result {
    cu := c.getUserSession()
    if nil == cu || !cu.IsValid() {
        c.Flash.Error(c.Message("error.need.signin"))
        return c.Redirect(routes.User.Signin())
    }

    // check if email has been verified
    if cu.EmailVerified {
        c.Flash.Error(c.Message("user.signup.mail.verify.exist"))
        return c.Redirect(routes.User.Home())
    }

    if c.sendConfirmEmail(cu.Name, cu.Email) {
        c.Flash.Success(c.Message("user.signup.mail.succeeded", cu.Email))
    } else {
        c.Flash.Error(c.Message("user.signup.mail.failed"))
    }
    return c.Redirect(routes.User.Home())
}

func (c User) DoVerifyEmail(name, key string) revel.Result {
    if "" == name || "" == key {
        c.Flash.Error("empty name or key")
        return c.Redirect(routes.User.Signin())
    }

    // init db connection
    dbmgr := database.NewDbManager()
    if nil == dbmgr {
        panic("db error: failed to init")
    }
    defer dbmgr.Close()

    // check user name
    validators.ValidateName(c.Validation, name)
    validators.ValidateDbNameExists(c.Validation, dbmgr.DbMap, name)
    if c.Validation.HasErrors() {
        revel.ERROR.Printf("validation for name %s failed", name)
        c.Validation.Keep()
        c.FlashParams()
        return c.Redirect(routes.User.Signin())
    }

    // check confirmation key first
    if !caching.CheckCachedSignupConfirmKey(name, key) {
        c.Flash.Error(c.Message("user.signup.key.error"))
        return c.Redirect(routes.User.Signin())
    }

    // update user
    u := database.FindUserByName(dbmgr.DbMap, name)
    if nil == u || !u.IsSecured() {
        revel.ERROR.Printf("invalid user in database %s", name)
        c.Flash.Error(c.Message("error.internal"))
        return c.Redirect(routes.User.Signin())
    }
    if u.EmailVerified {
        revel.WARN.Printf("user %s's email address had been verified before", name)
        return c.Redirect(routes.User.Signin())
    }

    u.EmailVerified = true
    if !database.UpdateUser(dbmgr.DbMap, *u) {
        revel.ERROR.Printf("failed to update user email verify status for %s", name)
        c.Flash.Error(c.Message("error.database"))
        return c.Redirect(routes.User.Signin())
    }
    // reset user cache
    caching.SetUser(u)

    revel.INFO.Printf("user %s has verified email address successfully")
    c.Flash.Success(c.Message("misc.signup.notice.verify.succeeded"))
    return c.Redirect(routes.User.Signin())
}

func (c User) Home() revel.Result {
    // check if user has signed in
    if cu := c.getUserSession(); nil == cu || !cu.IsValid() {
        c.Flash.Error(c.Message("error.need.signin"))
        return c.Redirect(routes.User.Signin())
    }

    moreNavbarLinks := []models.NavbarLink{
        models.NavbarLink{"/user/dosignout", "sign out", "sign out", false},
    }
    return c.Render(moreNavbarLinks)
}

