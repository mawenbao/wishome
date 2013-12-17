package controllers

import (
    "bytes"
    "github.com/robfig/revel"
    "github.com/mawenbao/wishome/app/results"
    "github.com/mawenbao/wishome/app/modules/ext/captcha"
)

type CaptchaC struct {
    *revel.Controller
}

// reload captcha by id and return its url
func (c CaptchaC) GetImageCaptcha(captchaID string) revel.Result {
    captchaImageBuff := new(bytes.Buffer)
    if "" == captcha.GetImageCaptcha(captchaID, captchaImageBuff) {
        revel.ERROR.Printf("failed to get catpcha, id %s, return 404", captchaID)
        return c.NotFound("captcha not found") // 404
    }

    return results.ImagePngResult(captchaImageBuff.Bytes())
}

