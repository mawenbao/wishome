package template

import (
    "fmt"
    "github.com/mawenbao/wishome/app"
)

type ConfirmEmailArgs struct {
    User,
    Host,
    Link,
    LinkLife string
}

type ResetPassEmailArgs struct {
    User,
    Host,
    Link,
    LinkLife string
}

func LoadConfirmEmail(user, link string) []byte {
    args := ConfirmEmailArgs {
        User: user,
        Host: app.MyGlobal.AppUrl,
        Link: link,
        LinkLife: fmt.Sprintf("%.0f minutes", app.MyGlobal.SignupKeyLife.Minutes()),
    }

    return LoadTempate("Confirmation Email Template", app.MyGlobal.TemplateConfirmMail, args)
}

func LoadResetPassEmail(user, link string) []byte {
    args := ResetPassEmailArgs {
        User: user,
        Host: app.MyGlobal.AppUrl,
        Link: link,
        LinkLife: fmt.Sprintf("%.0f minutes", app.MyGlobal.ResetPassKeyLife.Minutes()),
    }

    return LoadTempate("Reset Password Email Template", app.MyGlobal.TemplateResetPassMail, args)
}

