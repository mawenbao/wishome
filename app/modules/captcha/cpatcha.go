package captcha

import (
    "io"
    "github.com/robfig/revel"
    "github.com/dchest/captcha"
    "github.com/mawenbao/wishome/app"
)

func GenerateCaptchaImage(id string, out io.Writer) bool {
    if "" == id {
        revel.ERROR.Println("cannot generate captcha image for empty id")
        return false
    }
    err := captcha.WriteImage(out, id, app.MyGlobal.Int(app.CONFIG_CAPTCHA_WIDTH), app.MyGlobal.Int(app.CONFIG_CAPTCHA_HEIGHT))
    if nil != err {
        revel.ERROR.Printf("failed to generate captcha image for id %s", id)
        return false
    }
    return true
}

func NewCaptcha() string {
    return captcha.NewLen(app.MyGlobal.Int(app.CONFIG_CAPTCHA_LENGTH))
}

// realod captcha and save it
func ReloadCaptcha(id string) bool {
    if "" == id {
        revel.ERROR.Printf("try to reload a captcha with empty id")
        return false
    }
    if !captcha.Reload(id) {
        revel.ERROR.Printf("captcha reload failed, id %s does no exist", id)
        return false
    }
    return true
}

// get a new captcha if id is "" or reload
// the captcha identified by id
// return captcha id
func GetCaptchaWithImage(id string, out io.Writer) string {
    if "" == id {
        id = NewCaptcha()
    } else if !ReloadCaptcha(id) {
        // reload captcha failed, generate a new one
        id = NewCaptcha()
    }

    if !GenerateCaptchaImage(id, out) {
        return ""
    }
    return id
}

// captcha is removed immediately after verification, no matter passed or failed
func VerifyCaptchaString(id, value string) bool {
    if "" == id || "" == value {
        return false
    }
    return captcha.VerifyString(id, value)
}

