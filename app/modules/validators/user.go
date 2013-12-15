package validators

import (
    "regexp"
    "github.com/robfig/revel"
    "github.com/mawenbao/wishome/app"
    "github.com/mawenbao/wishome/app/models"
    "github.com/mawenbao/wishome/app/modules/database"
    "github.com/mawenbao/wishome/app/modules/common"
)

var (
    USER_NAME_REGEX = regexp.MustCompile(`^[a-zA-Z][-._0-9a-zA-Z]*[0-9a-zA-Z]$`)
    USER_EMAIL_REGEX = regexp.MustCompile(`^[0-9a-zA-Z]([-._0-9a-zA-Z]*[0-9a-zA-Z])?[@[-_0-9a-zA-Z]+\.[-._0-9a-zA-Z]+`)
)

func ValidateSignup(c *revel.Controller, name, email, password string) (result *revel.ValidationResult) {
    result = ValidateName(c, name)
    if !result.Ok {
        return
    }

    result = ValidateEmail(c, email)
    if !result.Ok {
        return
    }

    result = ValidatePassword(c, password)
    if !result.Ok {
        return
    }

    // check name and email in db
    result = ValidateDbNameNotExists(c, name)
    if !result.Ok {
        return
    }

    result = ValidateDbEmailNotExists(c, email)
    return
}

// allow user who have not verified their email address to sigin in
func ValidateSignin(c *revel.Controller, name, password string) (result *revel.ValidationResult, u *models.User) {
    result = ValidateName(c, name)
    if !result.Ok {
        return
    }

    u = database.FindUserByName(name)
    if nil == u || !u.IsValid() {
        c.Validation.Error(c.Message("user.signin.error.user")).Key(app.STR_NAME)
        return
    }

    result = ValidatePassword(c, password)
    if !result.Ok {
        return result, nil
    }

    result = ValidateDbPassword(c, password, u)
    return
}

func ValidateResetPassNameEmail(c *revel.Controller, name, email string) (result *revel.ValidationResult) {
    result = ValidateName(c, name)
    if !result.Ok {
        return
    }

    result = ValidateEmail(c, email)
    if !result.Ok {
        return
    }

    // make sure user has verified email address
    result = ValidateDbEmailVerified(c, name)
    if !result.Ok {
        return
    }

    result = ValidateDbNameEmail(c, name, email)
    return
}

func ValidateName(c *revel.Controller, name string) *revel.ValidationResult {
    return c.Validation.Check(name, revel.Required{}, revel.MaxSize{15}, revel.MinSize{4}, revel.Match{USER_NAME_REGEX})
}

func ValidateDbNameNotExists(c *revel.Controller, name string) *revel.ValidationResult {
    // name should not exists in db
    if database.IsNameExists(name) {
        return c.Validation.Error(c.Message("user.signup.failed.exist.name")).Key(app.STR_NAME)
    }
    return &revel.ValidationResult{Ok: true}
}

func ValidateDbNameExists(c *revel.Controller, name string) *revel.ValidationResult {
    // name should exists in db
    if !database.IsNameExists(name) {
        return c.Validation.Error(c.Message("user.signin.error.user")).Key(app.STR_NAME)
    }
    return &revel.ValidationResult{Ok: true}
}

func ValidateEmail(c *revel.Controller, email string) *revel.ValidationResult {
    return c.Validation.Check(email, revel.Required{}, revel.MaxSize{50}, revel.MinSize{6}, revel.Match{USER_EMAIL_REGEX})
}

func ValidateEmailVerified(c *revel.Controller, emailVerified bool) *revel.ValidationResult {
    if !emailVerified {
        return c.Validation.Error(c.Message("user.signin.failed.mail.notverified")).Key("email verify error")
    }
    return &revel.ValidationResult{Ok: true}
}

func ValidateDbEmailVerified(c *revel.Controller, name string) *revel.ValidationResult {
    if !database.IsEmailVerified(name) {
        return c.Validation.Error(c.Message("user.signin.failed.mail.notverified")).Key("email verify error")
    }
    return &revel.ValidationResult{Ok: true}
}

func ValidateDbEmailNotExists(c *revel.Controller, email string) *revel.ValidationResult {
    // email should not exists in db
    if database.IsEmailExists(email) {
        return c.Validation.Error(c.Message("user.signup.failed.exist.email")).Key(app.STR_EMAIL)
    }
    return &revel.ValidationResult{Ok: true}
}

func ValidatePassword(c *revel.Controller, password string) *revel.ValidationResult {
    return c.Validation.Check(password, revel.Required{}, revel.MinSize{4}, revel.MaxSize{15})
}

// encrypt password and compare it with saved pass in db
func ValidateDbPassword(c *revel.Controller, password string, u *models.User) *revel.ValidationResult {
    gotPass := common.MD5Sum(password + u.PassSalt)
    if gotPass != u.Password {
        return c.Validation.Error(c.Message("user.signin.error.password")).Key(app.STR_PASSWORD)
    }
    return &revel.ValidationResult{Ok: true}
}

func ValidateDbNameEmail(c *revel.Controller, name, email string) *revel.ValidationResult {
    if !database.IsNameEmailExists(name, email) {
        return c.Validation.Error(c.Message("user.resetpass.error.mismatch")).Key("reset password error")
    }
    return &revel.ValidationResult{Ok: true}
}

