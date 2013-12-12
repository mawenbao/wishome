package validators

import (
    "regexp"
    "github.com/coopernurse/gorp"
    "github.com/robfig/revel"
    "github.com/mawenbao/wishome/app/models"
    "github.com/mawenbao/wishome/app/modules/database"
    "github.com/mawenbao/wishome/app/modules/common"
)

var (
    USER_NAME_REGEX = regexp.MustCompile(`^[a-zA-Z][-._0-9a-zA-Z]*[0-9a-zA-Z]$`)
    USER_EMAIL_REGEX = regexp.MustCompile(`^[0-9a-zA-Z]([-._0-9a-zA-Z]*[0-9a-zA-Z])?[@[-_0-9a-zA-Z]+\.[-._0-9a-zA-Z]+`)
)

func ValidateSignup(v *revel.Validation, dbmap *gorp.DbMap, name, email, rawPass string) (result *revel.ValidationResult) {
    result = ValidateName(v, name)
    if !result.Ok {
        return
    }

    result = ValidateEmail(v, email)
    if !result.Ok {
        return
    }

    result = ValidatePassword(v, rawPass)
    if !result.Ok {
        return
    }

    // check name and email in db
    result = ValidateDbName(v, name, dbmap)
    if !result.Ok {
        return
    }

    result = ValidateDbEmail(v, email, dbmap)
    return
}

func ValidateSignin(v *revel.Validation, dbmap *gorp.DbMap, name string, rawPass string) (result *revel.ValidationResult, u *models.User) {
    result = ValidateName(v, name)
    if !result.Ok {
        return
    }

    u = database.FindUserByName(dbmap, name)
    if nil == u {
        v.Error("user not found: %s", name)
        return
    }

    result = ValidatePassword(v, rawPass)
    if !result.Ok {
        return result, nil
    }

    result = ValidateDbPassword(v, rawPass, u)
    return
}

func ValidateName(v *revel.Validation, name string) *revel.ValidationResult {
    return v.Check(name, revel.Required{}, revel.MaxSize{15}, revel.MinSize{4}, revel.Match{USER_NAME_REGEX})
}

func ValidateDbName(v *revel.Validation, name string, dbmap *gorp.DbMap) *revel.ValidationResult {
    // check name in db
    if database.IsNameExists(dbmap, name) {
        return v.Error("name already exists")
    }
    return &revel.ValidationResult{Ok: true}
}

func ValidateEmail(v *revel.Validation, email string) *revel.ValidationResult {
    return v.Check(email, revel.Required{}, revel.MaxSize{50}, revel.MinSize{6}, revel.Match{USER_EMAIL_REGEX})
}

func ValidateDbEmail(v *revel.Validation, email string, dbmap *gorp.DbMap) *revel.ValidationResult {
    // check email in db
    if database.IsEmailExists(dbmap, email) {
        return v.Error("email already exists")
    }
    return &revel.ValidationResult{Ok: true}
}

func ValidatePassword(v *revel.Validation, rawPass string) *revel.ValidationResult {
    return v.Check(rawPass, revel.Required{}, revel.MinSize{4}, revel.MaxSize{15})
}

// encrypt rawPass and compare it with saved pass in db
func ValidateDbPassword(v *revel.Validation, rawPass string, u *models.User) *revel.ValidationResult {
    gotPass := common.MD5Sum(rawPass + u.PassSalt)
    if gotPass != u.Password {
        return v.Error("wrong password").Key("password error")
    }
    return &revel.ValidationResult{Ok: true}
}

