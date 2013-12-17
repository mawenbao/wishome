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
    "github.com/mawenbao/wishome/app/modules/caching"
    "github.com/mawenbao/wishome/app/modules/mail"
    "github.com/mawenbao/wishome/app/modules/captcha"
    "github.com/mawenbao/wishome/app/modules/session"
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

func (c User) getLastUser() string {
    if lastuser, ok := c.Session[app.STR_LASTUSER]; !ok {
        return ""
    } else {
        lastuserSL, _ := common.DecodeHexString(lastuser)
        if nil == lastuserSL {
            revel.ERROR.Printf("failed to decode LastUser from user session %s", lastuser)
            return ""
        }
        return string(lastuserSL)
    }
}

func (c User) setLastUser(lastuser string) {
    c.Session[app.STR_LASTUSER] = common.EncodeToHexString([]byte(lastuser))
}

func (c User) setUserSession(u *models.User) bool {
    if nil == u || !u.IsSecured() {
        return false
    }

    sess := session.NewUserSession(u)
    if nil == sess {
        return false
    }
    sess.Encrypt()

    c.Session[app.STR_NAME] = sess.UserName
    c.Session[app.STR_PASSWORD] = sess.Password
    c.Session[app.STR_KEY] = sess.AesKey
    c.Session[app.STR_LASTUSER] = sess.LastUser
    c.Session[app.STR_EXPIRE] = sess.Expire
    return true
}

func (c User) loadUser() (u *models.User) {
    if "" == c.Session[app.STR_KEY] || "" == c.Session[app.STR_NAME] {
        return nil
    }

    // check if session expires
    if c.isSessionExpired() {
        return nil
    }

    sess := &session.UserSession {
        AesKey: c.Session[app.STR_KEY],
        UserName: c.Session[app.STR_NAME],
        Password: c.Session[app.STR_PASSWORD],
        LastUser: c.Session[app.STR_LASTUSER],
        Expire: c.Session[app.STR_EXPIRE],
        Encrypted: true,
    }

    if !sess.Decrypt() {
        revel.ERROR.Printf("failed to decrypt user session for user %s", sess.UserName)
        return nil
    }

    // get user from cache
    u = caching.GetUser(sess.UserName)
    if nil == u || !u.IsSecured() {
        // cache missed, session valid, try to reload user cache from database
        u = caching.ReloadUser(sess.UserName)
        if nil == u || !u.IsSecured() {
            revel.ERROR.Printf("failed to reload user cache for %s", sess.UserName)
            return nil
        }
    }

    if sess.UserName != u.Name || sess.Password != u.Password {
        return nil
    }

    return u
}

func (c User) emptyUserSession(name string) {
    // clear cookie
    for k := range c.Session {
        // do not delete LastUser in cookie
        if app.STR_LASTUSER != k {
            delete(c.Session, k)
        }
    }
    // clear server cache entry
    caching.Remove(name)
}

func (c User) checkCaptcha(captchaid, captchaval string) bool {
    if "" == captchaid || "" == captchaval {
        c.Flash.Error(c.Message("error.require.captcha"))
        return false
    }
    if !captcha.VerifyCaptchaString(captchaid, captchaval) {
        c.Flash.Error(c.Message("error.captcha"))
        return false
    }
    return true
}

func (c User) DoSignin(name, password, captchaid, captchaval string) revel.Result {
    if uc := c.loadUser(); nil != uc {
        c.emptyUserSession(name)
    }

    // save the name in cookie
    c.setLastUser(name)

    // check user signin error type and times
    if caching.IsSigninBanned(name) {
        c.Flash.Error(c.Message("user.signin.error.banned", caching.SIGNIN_SESSION_LIFE))
        return c.Redirect(routes.User.Signin())
    }

    // check captcha
    if caching.IsSigninCaptchaRequired(name) {
        if !c.checkCaptcha(captchaid, captchaval) {
            return c.Redirect(routes.User.Signin())
        }
    }

    // check user, password and set signin error
    _, user := validators.ValidateSignin(c.Controller, name, password)
    if c.Validation.HasErrors() {
        c.Validation.Keep()
        c.FlashParams()
        return c.Redirect(routes.User.Signin())
    }

    // save user in session and server cache
    c.setUserSession(user)
    caching.SetUser(user)

    c.Flash.Success(c.Message("user.greeting.old", name))
    revel.INFO.Printf("user %s signed in", name)
    return c.Redirect(routes.User.Home())
}

func (c User) DoSignup(name, email, password, captchaid, captchaval string) revel.Result {
    // check if user has signed in already, if so, sign him out
    if nil != c.loadUser() {
        c.emptyUserSession(name)
    }

    // check captcha
    if !c.checkCaptcha(captchaid, captchaval) {
        return c.Redirect(routes.User.Signup())
    }

    user := models.User{
        Name: name,
        Email: email,
        EmailVerified: false,
        Password: password,
    }

    // validate input 
    validators.ValidateSignup(c.Controller, user.Name, user.Email, user.Password)
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
    if !database.SaveUser(user) {
        c.Flash.Error(c.Message("user.save.error.db"))
        return c.Redirect(routes.User.Signup())
    } else {
        revel.INFO.Printf("new user %s signed up", name)
        c.Flash.Success(c.Message("misc.signup.notice.verify.email"))
        return c.Redirect(routes.User.Signin())
    }
}

func (c User) DoSignout() revel.Result {
    if uc := c.loadUser(); nil != uc && uc.IsSecured() {
        revel.INFO.Printf("user %s signed out", uc.Name)
        c.emptyUserSession(uc.Name)
    }
    c.Flash.Success(c.Message("user.signout.succeeded"))
    return c.Redirect(routes.User.Signin())
}

func (c User) Signin() revel.Result {
    // check if user has signed in
    if uc := c.loadUser(); nil != uc && uc.IsSecured() {
        return c.Redirect(routes.User.Home())
    }

    // get user name, may be empty
    name := c.getLastUser()
    needCaptcha := false

    // generate captcha if requried
    if caching.IsSigninCaptchaRequired(name) {
        needCaptcha = true
    }

    return c.Render(name, needCaptcha)
}

func (c User) Signup() revel.Result {
    if uc := c.loadUser(); nil != uc && uc.IsSecured() {
        c.Flash.Success(c.Message("user.signout.succeeded"))
        c.emptyUserSession(uc.Name)
    }

    moreNavbarLinks := []models.NavbarLink{
    }
    return c.Render(moreNavbarLinks)
}

func (c User) ResetPass() revel.Result {
    name := c.getLastUser()
    return c.Render(name)
}

// validate name, email and send an email to user with a random key
// which is valid in half an hour
func (c User) PreResetPass(name, email, captchaid, captchaval string) revel.Result {
    // check captcha
    if !c.checkCaptcha(captchaid, captchaval) {
        return c.Redirect(routes.User.ResetPass())
    }

    // validate name and email
    validators.ValidateResetPassNameEmail(c.Controller, name, email)
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

func (c User) DoResetPass(name, password, key, captchaid, captchaval string) revel.Result {
    revel.TRACE.Printf("doresetpass got key %s for name %s", key, name)
    // check captcha
    if !c.checkCaptcha(captchaid, captchaval) {
        return c.Redirect(routes.User.PostResetPass(name, key))
    }

    // check key first
    if !caching.CheckCachedResetPassKey(name, key) {
        c.Flash.Error(c.Message("user.resetpass.key.error"))
        return c.Redirect(routes.User.ResetPass())
    }

    // check password
    validators.ValidatePassword(c.Controller, password)
    if c.Validation.HasErrors() {
        c.Validation.Keep()
        c.FlashParams()
        return c.Redirect(routes.User.PostResetPass(name, key))
    }

    u := database.FindUserByName(name)
    if nil == u || !u.IsSecured() {
        c.Flash.Error(c.Message("common.user.none", name))
        return c.Redirect(routes.User.Signin())
    }

    // change password and password salt
    u.Password = password
    u.EncryptPass()
    if !database.UpdateUser(*u) {
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
        []byte(c.Message("user.resetpass.mail.content", resetPassUrl, app.MyGlobal.Duration(app.CONFIG_RESETPASS_KEY_LIFE).Minutes())),
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
        []byte(c.Message("user.signup.mail.content", cfmURL, app.MyGlobal.Duration(app.CONFIG_SIGNUP_KEY_LIFE).Minutes())),
    )
}

// user need to sign in first in order to resend confirmation email
func (c User) ResendConfirmEmail() revel.Result {
    cu := c.loadUser()
    if nil == cu || !cu.IsValid() {
        c.Flash.Error(c.Message("error.require.signin"))
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

    // check if user signed in
    nextPage := routes.User.Signin()
    if u := c.loadUser(); nil != u && u.IsSecured() {
        nextPage = routes.User.Home()
    }

    // check user name
    validators.ValidateName(c.Controller, name)
    validators.ValidateDbNameExists(c.Controller, name)
    if c.Validation.HasErrors() {
        revel.ERROR.Printf("validation for name %s failed", name)
        c.Validation.Keep()
        c.FlashParams()
        return c.Redirect(nextPage)
    }

    // check confirmation key first
    if !caching.CheckCachedSignupConfirmKey(name, key) {
        c.Flash.Error(c.Message("user.signup.key.error"))
        return c.Redirect(nextPage)
    }

    // update user
    u := database.FindUserByName(name)
    if nil == u || !u.IsSecured() {
        revel.ERROR.Printf("invalid user in database %s", name)
        c.Flash.Error(c.Message("error.internal"))
        return c.Redirect(nextPage)
    }
    if u.EmailVerified {
        revel.WARN.Printf("user %s's email address had been verified before", name)
        return c.Redirect(nextPage)
    }

    u.EmailVerified = true
    if !database.UpdateUser(*u) {
        revel.ERROR.Printf("failed to update user email verify status for %s", name)
        c.Flash.Error(c.Message("error.database"))
        return c.Redirect(nextPage)
    }
    // reset user cache
    caching.SetUser(u)

    revel.INFO.Printf("user %s has verified email address successfully", name)
    c.Flash.Success(c.Message("misc.signup.notice.verify.succeeded"))
    return c.Redirect(nextPage)
}

func (c User) Home() revel.Result {
    // check if user has signed in
    if cu := c.loadUser(); nil == cu || !cu.IsValid() {
        c.Flash.Error(c.Message("error.require.signin"))
        return c.Redirect(routes.User.Signin())
    }

    moreNavbarLinks := []models.NavbarLink{
        models.NavbarLink{"/user/dosignout", "sign out", "sign out", false},
    }
    return c.Render(moreNavbarLinks)
}

