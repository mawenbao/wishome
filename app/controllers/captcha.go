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
func (c Captcha) GetCaptcha(captchaid string) revel.Result {
    revel.ERROR.Println(captchaid)
    capResult := new(CaptchaResult)
    if "" == captchaid {
        capResult.ID = captcha.NewCaptcha()
    } else {
        capResult.ID = captchaid
        if !captcha.ReloadCaptcha(captchaid) {
            capResult.ID = captcha.NewCaptcha()
        }
    }

    capResult.ImageURL = "/captcha/getcaptchaimage?id=" + capResult.ID
    return c.RenderJson(capResult)
}

