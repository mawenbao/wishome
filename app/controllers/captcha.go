package controllers

import (
    "bytes"
    "github.com/robfig/revel"
    "github.com/mawenbao/wishome/app/results"
    "github.com/mawenbao/wishome/app/modules/captcha"
)

type CaptchaResult struct {
    ID string `json:"id"`
    ImageURL string `json:"imageurl"`
}

type Captcha struct {
    *revel.Controller
}

func (c Captcha) GetCaptchaImage(id string) revel.Result {
    captchaImageBuff := new(bytes.Buffer)
    if !captcha.GenerateCaptchaImage(id, captchaImageBuff) {
        revel.ERROR.Printf("failed to get catpcha, id %s, return 404", id)
        return c.NotFound("captcha not found") // 404
    }

    return results.ImagePngResult(captchaImageBuff.Bytes())
}

// reload captcha by id and return its url
func (c Captcha) GetCaptcha(captchaID string) revel.Result {
    capResult := new (CaptchaResult)
    if "" == captchaID {
        capResult.ID = captcha.NewCaptcha()
    } else {
        capResult.ID = captchaID
        if !captcha.ReloadCaptcha(captchaID) {
            capResult.ID = captcha.NewCaptcha()
        }
    }

    capResult.ImageURL = "/captcha/getcaptchaimage?id=" + capResult.ID
    return c.RenderJson(capResult)
}

