package controllers

import (
    "github.com/robfig/revel"
    "github.com/mawenbao/wishome/app/models"
    _ "github.com/mawenbao/wishome/app/modules/common"
    "github.com/mawenbao/wishome/app/modules/database"
    "github.com/mawenbao/wishome/app/modules/validators"
    "github.com/mawenbao/wishome/app/routes"
)

type User struct {
    *revel.Controller
}

func (c User) setUserSession(u models.User) {
    /*
    aesKey := common.NewRandomString(32) // aes-256
    c.Session["key"] = aesKey

    c.Session["user"] = common.AesEncrypt([]byte(u.Name), []byte(aesKey))
    c.Session["email"] = common.AesEncrypt([]byte(u.Email), []byte(aesKey))
    c.Session["pass"] = common.AesEncrypt([]byte(u.Password), []byte(aesKey))
    */
    c.Session["key"] = "abc"
    c.Session["user"] = u.Name
    c.Session["email"] = u.Email
    c.Session["pass"] = u.Password

    c.RenderArgs["user"] = &u
}

func (c User) getUserSession() (u *models.User) {
    if "" == c.Session["key"] || "" == c.Session["user"] || "" == c.Session["email"] || "" == c.Session["pass"] {
        return nil
    }
    u = new(models.User)

    /*
    key := []byte(c.Session["key"])
    u.Name = common.AesDecrypt([]byte(c.Session["user"]), key)
    u.Email = common.AesDecrypt([]byte(c.Session["email"]), key)
    u.Password = common.AesDecrypt([]byte(c.Session["pass"]), key)
    */
    u.Name = c.Session["user"]
    u.Email = c.Session["email"]
    u.Password = c.Session["pass"]

    if "" == u.Name || "" == u.Email || "" == u.Password {
        return nil
    }
    return u
}

func (c User) emptyUserSession() {
    c.RenderArgs["user"] = nil

    for k := range c.Session {
        if "user" != k {
            delete(c.Session, k)
        }
    }
}

func (c User) connected() *models.User {
    if nil != c.RenderArgs["user"] {
        return c.RenderArgs["user"].(*models.User)
    }

    // check session, with no ID and PassSalt field
    return c.getUserSession()
}

func (c User) DoSignin(name, rawPass string) revel.Result {
    if cu := c.connected(); nil != cu {
        c.Flash.Success("welcome back, %s", cu.Name)
        return c.Redirect(routes.User.Home())
    }

    dbmgr := database.NewDbManager()
    if nil == dbmgr {
        panic("db error: failed to init")
    }
    defer dbmgr.Close()

    _, user := validators.ValidateSignin(c.Validation, dbmgr.DbMap, name, rawPass)
    if c.Validation.HasErrors() {
        c.Validation.Keep()
        c.FlashParams()
        return c.Redirect(routes.User.Signin())
    }

    c.setUserSession(*user)
    c.Flash.Success("welcome back %s", name)
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
        c.Flash.Error("failed to save user in database")
        return c.Redirect(routes.User.Signup())
    } else {
        // set session
        c.setUserSession(user)
        c.Flash.Success("welcome back, %s", user.Name)
        return c.Redirect(routes.User.Home())
    }
}

func (c User) DoSignout() revel.Result {
    c.emptyUserSession()
    return c.Redirect(routes.User.Signin())
}

func (c User) Signin() revel.Result {
    return c.Render()
}

func (c User) Signup() revel.Result {
    c.emptyUserSession()
    return c.Render()
}

func (c User) Home() revel.Result {
    return c.Render()
}

