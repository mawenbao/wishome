package app

import (
    "log"
    "time"
    "net"
    "strings"
    "strconv"
    "github.com/robfig/revel"
)

// store global custom configuration in app.conf
type MyGlobalConfig map[string]interface{}

func (my MyGlobalConfig) String(key string) string {
    return my[key].(string)
}

func (my MyGlobalConfig) Int(key string) int {
    return my[key].(int)
}

func (my MyGlobalConfig) Duration(key string) time.Duration {
    return my[key].(time.Duration)
}

var MyGlobal = make(MyGlobalConfig)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.ActionInvoker,           // Invoke the action.
	}

    // set log flags
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

    // init my global
    revel.OnAppStart(parseDbConfig)
    revel.OnAppStart(parseCustomConfig)
}

// parse db related config
func parseDbConfig() {
    if dbdriver, found := revel.Config.String("db.driver"); !found {
        revel.ERROR.Panic("app.conf error: db.driver not defined")
    } else {
        MyGlobal[CONFIG_DB_DRIVER] = dbdriver
    }
    if dbspec, found := revel.Config.String("db.spec"); !found {
        revel.ERROR.Panic("app.conf error: db.spec not defined")
    } else {
        MyGlobal[CONFIG_DB_SPEC] = dbspec
    }
}

// parse custom settings in app.conf and save in MyGlobal
// parse failure will result in panic
func parseCustomConfig() {
    var err error
    // parse app url, app.url must be set in app.conf file
    appURL, found := revel.Config.String(CONFIG_APP_URL)
    if !found {
        revel.ERROR.Panicf("%s not set in app.conf. This is required configuration which denotes the host/ip address of your wishome app.", CONFIG_APP_URL)
    } else {
        MyGlobal[CONFIG_APP_URL] = strings.TrimRight(appURL, "/")
    }

    // parse session lifetime
    sessLife, found := revel.Config.String(CONFIG_SESSION_LIFE)
    if !found {
        MyGlobal[CONFIG_SESSION_LIFE] = 1 * time.Hour
    } else {
        sessLifeTime, err := time.ParseDuration(sessLife)
        if nil != err {
            revel.ERROR.Panicf("failed to parse %s in app.conf, value is %s", CONFIG_SESSION_LIFE, sessLife)
        }
        MyGlobal[CONFIG_SESSION_LIFE] = sessLifeTime
    }

    // parse reset password key length
    rstKeyLen, found := revel.Config.Int(CONFIG_RESETPASS_KEY_LEN)
    if !found {
        MyGlobal[CONFIG_RESETPASS_KEY_LEN] = 32
    } else {
        MyGlobal[CONFIG_RESETPASS_KEY_LEN] = rstKeyLen
    }

    // parse reset password key cache expire time
    rstKeyLife, found := revel.Config.String(CONFIG_RESETPASS_KEY_LIFE)
    if !found {
        MyGlobal[CONFIG_RESETPASS_KEY_LIFE] = time.Duration(30 * time.Minute)
    } else {
        keyLifeTime, err := time.ParseDuration(rstKeyLife)
        if nil != err {
            revel.ERROR.Panicf("failed to parse %s in app.conf %s", CONFIG_RESETPASS_KEY_LIFE, rstKeyLife)
        }
        MyGlobal[CONFIG_RESETPASS_KEY_LIFE] = keyLifeTime
    }

    // parse user signup email confirmation key length
    cfmKeyLen, found := revel.Config.Int(CONFIG_SIGNUP_KEY_LEN)
    if !found {
        MyGlobal[CONFIG_SIGNUP_KEY_LEN] = 32
    } else {
        MyGlobal[CONFIG_SIGNUP_KEY_LEN] = cfmKeyLen
    }

    // parse user signup email confirmation key expire time
    cfmKeyLife, found := revel.Config.String(CONFIG_SIGNUP_KEY_LIFE)
    if !found {
        MyGlobal[CONFIG_SIGNUP_KEY_LIFE] = time.Duration(30 * time.Minute)
    } else {
        keyLifeTime, err := time.ParseDuration(cfmKeyLife)
        if nil != err {
            revel.ERROR.Panicf("failed to parse %s in app.conf %s", CONFIG_SIGNUP_KEY_LIFE, cfmKeyLife)
        }
        MyGlobal[CONFIG_SIGNUP_KEY_LIFE] = keyLifeTime
    }

    // parse mail smtp server address
    smtpServer, found := revel.Config.String(CONFIG_MAIL_SMTP_ADDR)
    if !found {
        MyGlobal[CONFIG_MAIL_SMTP_ADDR] = "localhost:25"
        MyGlobal[CONFIG_MAIL_SMTP_HOST] = "localhost"
        MyGlobal[CONFIG_MAIL_SMTP_PORT] = 25
    } else {
        MyGlobal[CONFIG_MAIL_SMTP_ADDR] = smtpServer
        host, port, err := net.SplitHostPort(smtpServer)
        if nil != err {
            revel.ERROR.Panicf("failed to split host:port from %s, value is %s", CONFIG_MAIL_SMTP_ADDR, smtpServer)
        }
        MyGlobal[CONFIG_MAIL_SMTP_HOST] = host
        MyGlobal[CONFIG_MAIL_SMTP_PORT], err = strconv.Atoi(port)
        if nil != err {
            revel.ERROR.Panicf("failed to convert port string %s to int", port)
        }
    }

    // parse mail sender
    mailSender, found := revel.Config.String(CONFIG_MAIL_SENDER)
    if !found {
        MyGlobal[CONFIG_MAIL_SENDER] = "noreply@atime.me"
    } else {
        MyGlobal[CONFIG_MAIL_SENDER] = mailSender
    }

    // parse captcha
    captchaLen, found := revel.Config.String(CONFIG_CAPTCHA_LENGTH)
    if !found {
        MyGlobal[CONFIG_CAPTCHA_LENGTH] = 6
    } else {
        MyGlobal[CONFIG_CAPTCHA_LENGTH], err = strconv.Atoi(captchaLen)
        if nil != err {
            revel.ERROR.Panicf("failed to convert captcha length %s to int", captchaLen)
        }
    }
    captchaWidth, found := revel.Config.String(CONFIG_CAPTCHA_WIDTH)
    if !found {
        MyGlobal[CONFIG_CAPTCHA_WIDTH] = 100
    } else {
        MyGlobal[CONFIG_CAPTCHA_WIDTH], err = strconv.Atoi(captchaWidth)
        if nil != err {
            revel.ERROR.Panicf("failed to convert captcha width %s to int", captchaWidth)
        }
    }
    captchaHeight, found := revel.Config.String(CONFIG_CAPTCHA_HEIGHT)
    if !found {
        MyGlobal[CONFIG_CAPTCHA_HEIGHT] = 40
    } else {
        MyGlobal[CONFIG_CAPTCHA_HEIGHT], err = strconv.Atoi(captchaHeight)
        if nil != err {
            revel.ERROR.Panicf("failed to convert captcha height %s to int", captchaHeight)
        }
    }
}

