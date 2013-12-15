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

func ValidateSignup(v *revel.Validation, name, email, password string) (result *revel.ValidationResult) {
    result = ValidateName(v, name)
    if !result.Ok {
        return
    }

    result = ValidateEmail(v, email)
    if !result.Ok {
        return
    }

    result = ValidatePassword(v, password)
    if !result.Ok {
        return
    }

    // check name and email in db
    result = ValidateDbNameNotExists(v, name)
    if !result.Ok {
        return
    }

    result = ValidateDbEmailNotExists(v, email)
    return
}

// allow user who have not verified their email address to sigin in
func ValidateSignin(v *revel.Validation, name, password string) (result *revel.ValidationResult, u *models.User) {
    result = ValidateName(v, name)
    if !result.Ok {
        return
    }

    u = database.FindUserByName(name)
    if nil == u || !u.IsValid() {
        v.Error("user not found").Key(app.STR_NAME)
        return
    }

    result = ValidatePassword(v, password)
    if !result.Ok {
        return result, nil
    }

    result = ValidateDbPassword(v, password, u)
    return
}

func ValidateResetPassNameEmail(v *revel.Validation, name, email string) (result *revel.ValidationResult) {
    result = ValidateName(v, name)
    if !result.Ok {
        return
    }

    result = ValidateEmail(v, email)
    if !result.Ok {
        return
    }

    // make sure user has verified email address
    result = ValidateDbEmailVerified(v, name)
    if !result.Ok {
        return
    }

    result = ValidateDbNameEmail(v, name, email)
    return
}

func ValidateName(v *revel.Validation, name string) *revel.ValidationResult {
    return v.Check(name, revel.Required{}, revel.MaxSize{15}, revel.MinSize{4}, revel.Match{USER_NAME_REGEX})
}

func ValidateDbNameNotExists(v *revel.Validation, name string) *revel.ValidationResult {
    // name should not exists in db
    if database.IsNameExists(name) {
        return v.Error("name already exists").Key(app.STR_NAME)
    }
    return &revel.ValidationResult{Ok: true}
}

func ValidateDbNameExists(v *revel.Validation, name string) *revel.ValidationResult {
    // name should exists in db
    if !database.IsNameExists(name) {
        return v.Error("name not exists").Key(app.STR_NAME)
    }
    return &revel.ValidationResult{Ok: true}
}

func ValidateEmail(v *revel.Validation, email string) *revel.ValidationResult {
    return v.Check(email, revel.Required{}, revel.MaxSize{50}, revel.MinSize{6}, revel.Match{USER_EMAIL_REGEX})
}

func ValidateEmailVerified(v *revel.Validation, emailVerified bool) *revel.ValidationResult {
    if !emailVerified {
        return v.Error("your email address has not been verified").Key("email verify error")
    }
    return &revel.ValidationResult{Ok: true}
}

func ValidateDbEmailVerified(v *revel.Validation, name string) *revel.ValidationResult {
    if !database.IsEmailVerified(name) {
        return v.Error("you have to verify your email address first").Key("email verify error")
    }
    return &revel.ValidationResult{Ok: true}
}

func ValidateDbEmailNotExists(v *revel.Validation, email string) *revel.ValidationResult {
    // email should not exists in db
    if database.IsEmailExists(email) {
        return v.Error("email already exists").Key(app.STR_EMAIL)
    }
    return &revel.ValidationResult{Ok: true}
}

func ValidatePassword(v *revel.Validation, password string) *revel.ValidationResult {
    return v.Check(password, revel.Required{}, revel.MinSize{4}, revel.MaxSize{15})
}

// encrypt password and compare it with saved pass in db
func ValidateDbPassword(v *revel.Validation, password string, u *models.User) *revel.ValidationResult {
    gotPass := common.MD5Sum(password + u.PassSalt)
    if gotPass != u.Password {
        return v.Error("wrong password").Key(app.STR_PASSWORD)
    }
    return &revel.ValidationResult{Ok: true}
}

func ValidateDbNameEmail(v *revel.Validation, name, email string) *revel.ValidationResult {
    if !database.IsNameEmailExists(name, email) {
        return v.Error("name and email do not match").Key("reset password error")
    }
    return &revel.ValidationResult{Ok: true}
}


