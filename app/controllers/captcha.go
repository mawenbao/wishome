package controllers

import (
    "fmt"
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

func (c Captcha) GetCaptchaImage(id, v string) revel.Result {
    captchaImageBuff := new(bytes.Buffer)
    if !captcha.GenerateCaptchaImage(id, captchaImageBuff) {
        revel.ERROR.Printf("failed to get catpcha, id %s version %s, return 404", id, v)
        return c.NotFound("captcha not found") // 404
    }

    return results.ImagePngResult(captchaImageBuff.Bytes())
}

// reload captcha by id and return its url
func (c Captcha) GetCaptcha(captchaid string) revel.Result {
    capResult := new(CaptchaResult)
    if "" == captchaid {
        capResult.ID = captcha.NewCaptcha()
    } else {
        capResult.ID = captchaid
        if !captcha.ReloadCaptcha(captchaid) {
            revel.ERROR.Printf("failed to reload captcha id %s", captchaid)
            capResult.ID = captcha.NewCaptcha()
        }
    }

    capResult.ImageURL = fmt.Sprintf(
        "/captcha/getcaptchaimage?id=%s",
        capResult.ID,
    )
    return c.RenderJson(capResult)
}

