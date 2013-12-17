package captcha

import (
    "io"
    "github.com/robfig/revel"
    "github.com/dchest/captcha"
)

func NewImageCaptcha(out io.Writer) string {
    id := captcha.NewLen(6)
    err := captcha.WriteImage(out, id, 100, 50)
    if nil != err {
        revel.ERROR.Printf("failed to generate captcha image")
        return ""
    }
    return id
}

// realod captcha and save it
func ReloadImageCaptcha(id string, out io.Writer) bool {
    if "" == id {
        revel.ERROR.Printf("try to reload a captcha with empty id")
        return false
    }
    if !captcha.Reload(id) {
        revel.ERROR.Printf("captcha reload failed, id %s does no exist", id)
        return false
    }

    err := captcha.WriteImage(out, id, 100, 50)
    if nil != err {
        revel.ERROR.Printf("failed to generate captcha image")
        return false
    }
    return true
}

// get a new captcha if id is "" or reload
// the captcha identified by id
// return captcha id
func GetImageCaptcha(id string, out io.Writer) string {
    if "" == id {
        return NewImageCaptcha(out)
    } else {
        if !ReloadImageCaptcha(id, out) {
            return ""
        }
        return id
    }
}

func VerifyCaptchaString(id, value string) bool {
    return captcha.VerifyString(id, value)
}

