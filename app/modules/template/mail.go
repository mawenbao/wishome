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
        Host: app.MyGlobal.String(app.CONFIG_APP_URL),
        Link: link,
        LinkLife: fmt.Sprintf("%.0f minutes", app.MyGlobal.Duration(app.CONFIG_SIGNUP_KEY_LIFE).Minutes()),
    }

    return LoadTempate("Confirmation Email Template", app.MyGlobal.String(app.CONFIG_TEMPLATE_CONFIRM_EMAIL), args)
}

func LoadResetPassEmail(user, link string) []byte {
    args := ResetPassEmailArgs {
        User: user,
        Host: app.MyGlobal.String(app.CONFIG_APP_URL),
        Link: link,
        LinkLife: fmt.Sprintf("%.0f minutes", app.MyGlobal.Duration(app.CONFIG_RESETPASS_KEY_LIFE).Minutes()),
    }

    return LoadTempate("Reset Password Email Template", app.MyGlobal.String(app.CONFIG_TEMPLATE_RESETPASS_EMAIL), args)
}

