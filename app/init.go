package app

import (
    "log"
    "time"
    "net"
    "path/filepath"
    "io/ioutil"
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
    var found bool
    // parse app url, app.url must be set in app.conf file
    MyGlobal[CONFIG_APP_URL], found = revel.Config.String(CONFIG_APP_URL)
    if !found {
        revel.ERROR.Panicf("%s not set in app.conf. This is required configuration which denotes the host/ip address of your wishome app.", CONFIG_APP_URL)
    }

    // parse session lifetime
    sessLife, found := revel.Config.String(CONFIG_SESSION_LIFE)
    if !found {
        MyGlobal[CONFIG_SESSION_LIFE] = 1 * time.Hour
    } else {
        MyGlobal[CONFIG_SESSION_LIFE], err = time.ParseDuration(sessLife)
        if nil != err {
            revel.ERROR.Panicf("failed to parse %s in app.conf, value is %s", CONFIG_SESSION_LIFE, sessLife)
        }
    }

    // parse reset password key length
    MyGlobal[CONFIG_RESETPASS_KEY_LEN], found = revel.Config.Int(CONFIG_RESETPASS_KEY_LEN)
    if !found {
        MyGlobal[CONFIG_RESETPASS_KEY_LEN] = 32
    }

    // parse reset password key cache expire time
    rstKeyLife, found := revel.Config.String(CONFIG_RESETPASS_KEY_LIFE)
    if !found {
        MyGlobal[CONFIG_RESETPASS_KEY_LIFE] = time.Duration(30 * time.Minute)
    } else {
        MyGlobal[CONFIG_RESETPASS_KEY_LIFE], err = time.ParseDuration(rstKeyLife)
        if nil != err {
            revel.ERROR.Panicf("failed to parse %s in app.conf %s", CONFIG_RESETPASS_KEY_LIFE, rstKeyLife)
        }
    }

    // parse user signup email confirmation key length
    MyGlobal[CONFIG_SIGNUP_KEY_LEN], found = revel.Config.Int(CONFIG_SIGNUP_KEY_LEN)
    if !found {
        MyGlobal[CONFIG_SIGNUP_KEY_LEN] = 32
    }

    // parse user signup email confirmation key expire time
    cfmKeyLife, found := revel.Config.String(CONFIG_SIGNUP_KEY_LIFE)
    if !found {
        MyGlobal[CONFIG_SIGNUP_KEY_LIFE] = time.Duration(30 * time.Minute)
    } else {
        MyGlobal[CONFIG_SIGNUP_KEY_LIFE], err = time.ParseDuration(cfmKeyLife)
        if nil != err {
            revel.ERROR.Panicf("failed to parse %s in app.conf %s", CONFIG_SIGNUP_KEY_LIFE, cfmKeyLife)
        }
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
    MyGlobal[CONFIG_MAIL_SENDER], found = revel.Config.String(CONFIG_MAIL_SENDER)
    if !found {
        MyGlobal[CONFIG_MAIL_SENDER] = "noreply@atime.me"
    }

    // parse captcha
    MyGlobal[CONFIG_CAPTCHA_LENGTH], found = revel.Config.Int(CONFIG_CAPTCHA_LENGTH)
    if !found {
        MyGlobal[CONFIG_CAPTCHA_LENGTH] = 6
    }
    MyGlobal[CONFIG_CAPTCHA_WIDTH], found = revel.Config.Int(CONFIG_CAPTCHA_WIDTH)
    if !found {
        MyGlobal[CONFIG_CAPTCHA_WIDTH] = 100
    }
    MyGlobal[CONFIG_CAPTCHA_HEIGHT], found = revel.Config.Int(CONFIG_CAPTCHA_HEIGHT)
    if !found {
        MyGlobal[CONFIG_CAPTCHA_HEIGHT] = 40
    }

    // server side session cache
    sessLifeSignin, found := revel.Config.String(CONFIG_SIGNIN_CACHE_LIFE)
    if !found {
        MyGlobal[CONFIG_SIGNIN_CACHE_LIFE] = time.Duration(1 * time.Hour)
    } else {
        MyGlobal[CONFIG_SIGNIN_CACHE_LIFE], err = time.ParseDuration(sessLifeSignin)
        if nil != err {
            revel.ERROR.Panicf("failed to parse %s as duration, value is %s", CONFIG_SIGNIN_CACHE_LIFE, sessLifeSignin)
        }
    }
    MyGlobal[CONFIG_SIGNIN_ERROR_LIMIT], found = revel.Config.Int(CONFIG_SIGNIN_ERROR_LIMIT)
    if !found {
        MyGlobal[CONFIG_SIGNIN_ERROR_LIMIT] = 30
    }

    // after some signin errors, user should enter captcha to continue signing in
    MyGlobal[CONFIG_SIGNIN_USECAPTCHA], found = revel.Config.Int(CONFIG_SIGNIN_USECAPTCHA)
    if !found {
        MyGlobal[CONFIG_SIGNIN_USECAPTCHA] = 5
    }

    // lock user signin after too many errors in this time
    signinBanTime, found := revel.Config.String(CONFIG_SIGNIN_BAN_TIME)
    if !found {
        MyGlobal[CONFIG_SIGNIN_BAN_TIME] = time.Duration(1 * time.Hour)
    } else {
        MyGlobal[CONFIG_SIGNIN_BAN_TIME], err = time.ParseDuration(signinBanTime)
        if nil != err {
            revel.ERROR.Panicf("failed to parse user signin ban time %s, value is %s", CONFIG_SIGNIN_BAN_TIME, signinBanTime)
        }
    }

    // parse custom template paths and save the file content
    MyGlobal[CONFIG_TEMPLATE_CONFIRM_EMAIL], found = revel.Config.String(CONFIG_TEMPLATE_CONFIRM_EMAIL)
    if !found {
        MyGlobal[CONFIG_TEMPLATE_CONFIRM_EMAIL] = "data/ConfirmationEmail.html"
    }
    MyGlobal[CONFIG_TEMPLATE_CONFIRM_EMAIL] = filepath.Join(revel.BasePath, MyGlobal.String(CONFIG_TEMPLATE_CONFIRM_EMAIL))
    confirmTemplData, err := ioutil.ReadFile(MyGlobal.String(CONFIG_TEMPLATE_CONFIRM_EMAIL))
    if nil != err {
        revel.ERROR.Panicf("failed to read confirmation template from %s: %s", MyGlobal.String(CONFIG_TEMPLATE_CONFIRM_EMAIL), err)
    }
    MyGlobal[CONFIG_TEMPLATE_CONFIRM_EMAIL] = string(confirmTemplData)

    MyGlobal[CONFIG_TEMPLATE_RESETPASS_EMAIL], found = revel.Config.String(CONFIG_TEMPLATE_RESETPASS_EMAIL)
    if !found {
        MyGlobal[CONFIG_TEMPLATE_RESETPASS_EMAIL] = "data/ResetPassEmail.html"
    }
    MyGlobal[CONFIG_TEMPLATE_RESETPASS_EMAIL] = filepath.Join(revel.BasePath, MyGlobal.String(CONFIG_TEMPLATE_RESETPASS_EMAIL))
    resetPassTemplData, err := ioutil.ReadFile(MyGlobal.String(CONFIG_TEMPLATE_RESETPASS_EMAIL))
    if nil != err {
        revel.ERROR.Panicf("failed to read resetpass template from %s: %s", MyGlobal.String(CONFIG_TEMPLATE_RESETPASS_EMAIL), err)
    }
    MyGlobal[CONFIG_TEMPLATE_RESETPASS_EMAIL] = string(resetPassTemplData)
}

