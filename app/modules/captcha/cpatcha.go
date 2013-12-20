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
    err := captcha.WriteImage(out, id, app.MyGlobal.CaptchaImgWidth, app.MyGlobal.CaptchaImgHeight)
    if nil != err {
        revel.ERROR.Printf("failed to generate captcha image for id %s", id)
        return false
    }
    return true
}

func NewCaptcha() string {
    return captcha.NewLen(app.MyGlobal.CaptchaLength)
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
    // disable captcha checking in debug mode
    if revel.DevMode {
        revel.WARN.Printf("catpcha verification will always return true in dev mode")
        return true
    }

    if "" == id || "" == value {
        return false
    }
    return captcha.VerifyString(id, value)
}

