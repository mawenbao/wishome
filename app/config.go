package app

import (
    "time"
    "fmt"
    "net"
    "strings"
    "runtime"
    "path/filepath"
    "io/ioutil"
    "strconv"
    "github.com/robfig/revel"
)

const (
    DFT_KEY_LIFE = time.Hour
    DFT_KEY_LEN = 32
    DFT_CACHE_LIFE = time.Hour
    DFT_CACHE_SIGNIN_ERROR_LIFE = time.Hour
    DFT_SESSION_LIFE = time.Hour
    DFT_CAPTCHA_LEN = 6
    DFT_CAPTCHA_WIDTH = 100
    DFT_CAPTCHA_HEIGHT = 45
    DFT_SIGNIN_USE_CAPTCHA = 5
    DFT_SIGNIN_ERR_LIMIT = 20
    DFT_SIGNIN_BAN_TIME = time.Hour
    DFT_REDIS_POOL_MAXIDLE = 20
    DFT_REDIS_POOL_IDLE_TIMEOUT = 3 * time.Minute
)

// global custom config
var MyGlobal = new(MyGlobalConfig)

// store global custom configuration in app.conf
type MyGlobalConfig struct {
    AdminIPList []string
    AdminTimer bool

    AppUrl string
    AppCpuNum int
    AdminIpList []string

    ResetPassKeyLen int
    ResetPassKeyLife time.Duration
    SignupKeyLen int
    SignupKeyLife time.Duration

    MailServerAddr string
    MailServerHost string
    MailServerPort int
    MailSender string

    CaptchaLength,
    CaptchaImgWidth,
    CaptchaImgHeight int

    DbDriver,
    DbSpec string

    SessionLife time.Duration
    CacheLife time.Duration

    SigninErrCacheLife time.Duration
    SigninUseCaptchaErrorNum int
    SigninErrLimit int
    SigninBanTime time.Duration

    TemplateConfirmMail,
    TemplateResetPassMail string

    RedisServerAddr string
    RedisPoolMaxIdle int
    RedisPoolIdleTimeout time.Duration
}

func parseOnOff(str string) bool {
    if STR_ON == strings.ToLower(strings.TrimSpace(str)) {
        return true
    } else if STR_OFF == strings.ToLower(strings.TrimSpace(str)) {
        return false
    } else {
        revel.ERROR.Panicf("unrecogonized on/off string %s", str)
    }
    return false
}

func (gconf *MyGlobalConfig) IsAdminIP(host string) bool {
    for _, adminIP := range(gconf.AdminIPList) {
        if adminIP == host {
            return true
        }
    }
    return false
}

// panic
func configError(entry, err string) {
    revel.ERROR.Panicf("app.conf error, config entry %s: %s", entry, err)
}

// panic
func configNotFound(entry string) {
    configError(entry, "config entry not found")
}

func configNotFoundWarn(entry string) {
    revel.WARN.Println("app.conf warning, config entry %s not found, will use default value")
}

// parse custom settings in app.conf and save in MyGlobal
// parse failure will result in panic
func parseCustomConfig() {
    var err error
    var found bool
    // parse app url, app.url must be set in app.conf file
    MyGlobal.AppUrl, found = revel.Config.String(CONFIG_APP_URL)
    if !found {
        configNotFound(CONFIG_APP_URL)
    }

    // parse admin ip list
    ipListStr, found := revel.Config.String(CONFIG_ADMIN_IP_LIST)
    if !found {
        MyGlobal.AdminIPList = []string{"127.0.0.1"}
    } else {
        for _, ip := range strings.Split(ipListStr, ":") {
            MyGlobal.AdminIPList = append(MyGlobal.AdminIPList, strings.TrimSpace(ip))
        }
    }

    // parse admin timer on/off
    adminTimerOnOff, found := revel.Config.String(CONFIG_ADMIN_TIMER)
    if !found {
        MyGlobal.AdminTimer = false // default off
    }
    MyGlobal.AdminTimer = parseOnOff(adminTimerOnOff)

    // parse db related config
    if MyGlobal.DbDriver, found = revel.Config.String(CONFIG_DB_DRIVER); !found {
        configNotFound(CONFIG_DB_DRIVER)
    }
    if MyGlobal.DbSpec, found = revel.Config.String(CONFIG_DB_SPEC); !found {
        configNotFound(CONFIG_DB_SPEC)
    }

    // parse cpu number
    MyGlobal.AppCpuNum, found = revel.Config.Int(CONFIG_APP_CPU_NUM)
    if !found {
        configNotFoundWarn(CONFIG_APP_CPU_NUM)
        MyGlobal.AppCpuNum = runtime.NumCPU()
    } else {
        if 1 > MyGlobal.AppCpuNum {
            revel.WARN.Printf("app.conf warning, config entry %s invalid, will use default value", CONFIG_APP_CPU_NUM)
            MyGlobal.AppCpuNum = runtime.NumCPU()
        }
    }
    runtime.GOMAXPROCS(MyGlobal.AppCpuNum)

    // parse session lifetime
    sessLife, found := revel.Config.String(CONFIG_SESSION_LIFE)
    if !found {
        MyGlobal.SessionLife = DFT_SESSION_LIFE
    } else {
        MyGlobal.SessionLife, err = time.ParseDuration(sessLife)
        if nil != err {
            configError(CONFIG_SESSION_LIFE, fmt.Sprintf("failed to parse %s", sessLife))
        }
    }

    // parse reset password key length
    MyGlobal.ResetPassKeyLen, found = revel.Config.Int(CONFIG_RESETPASS_KEY_LEN)
    if !found {
        MyGlobal.ResetPassKeyLen = DFT_KEY_LEN
    }

    // parse reset password key cache expire time
    rstKeyLife, found := revel.Config.String(CONFIG_RESETPASS_KEY_LIFE)
    if !found {
        MyGlobal.ResetPassKeyLife = DFT_KEY_LIFE
    } else {
        MyGlobal.ResetPassKeyLife, err = time.ParseDuration(rstKeyLife)
        if nil != err {
            configError(CONFIG_RESETPASS_KEY_LIFE, fmt.Sprintf("failed to parse %s", rstKeyLife))
        }
    }

    // parse user signup email confirmation key length
    MyGlobal.SignupKeyLen, found = revel.Config.Int(CONFIG_SIGNUP_KEY_LEN)
    if !found {
        MyGlobal.SignupKeyLen = DFT_KEY_LEN
    }

    // parse user signup email confirmation key expire time
    cfmKeyLife, found := revel.Config.String(CONFIG_SIGNUP_KEY_LIFE)
    if !found {
        MyGlobal.SignupKeyLife = DFT_KEY_LIFE
    } else {
        MyGlobal.SignupKeyLife, err = time.ParseDuration(cfmKeyLife)
        if nil != err {
            configError(CONFIG_SIGNUP_KEY_LIFE, fmt.Sprintf("failed to parse %s", cfmKeyLife))
        }
    }

    // parse mail smtp server address
    MyGlobal.MailServerAddr, found = revel.Config.String(CONFIG_MAIL_SMTP_ADDR)
    if !found {
        configNotFound(CONFIG_MAIL_SMTP_ADDR)
    } else {
        host, port, err := net.SplitHostPort(MyGlobal.MailServerAddr)
        if nil != err {
            configError(CONFIG_MAIL_SMTP_ADDR, fmt.Sprintf("failed to split host:port %s", MyGlobal.MailServerAddr))
        }
        MyGlobal.MailServerHost = host
        MyGlobal.MailServerPort, err = strconv.Atoi(port)
        if nil != err {
            configError(CONFIG_MAIL_SMTP_ADDR, fmt.Sprintf("failed to convert port string %s to int", port))
        }
    }

    // parse mail sender
    MyGlobal.MailSender, found = revel.Config.String(CONFIG_MAIL_SENDER)
    if !found {
        configNotFound(CONFIG_MAIL_SENDER)
    }

    // parse captcha
    MyGlobal.CaptchaLength, found = revel.Config.Int(CONFIG_CAPTCHA_LENGTH)
    if !found {
        MyGlobal.CaptchaLength = DFT_CAPTCHA_LEN
    }
    MyGlobal.CaptchaImgWidth, found = revel.Config.Int(CONFIG_CAPTCHA_WIDTH)
    if !found {
        MyGlobal.CaptchaImgWidth = DFT_CAPTCHA_WIDTH
    }
    MyGlobal.CaptchaImgHeight, found = revel.Config.Int(CONFIG_CAPTCHA_HEIGHT)
    if !found {
        MyGlobal.CaptchaImgHeight = DFT_CAPTCHA_HEIGHT
    }

    // server side session cache
    sessLifeSignin, found := revel.Config.String(CONFIG_SIGNIN_CACHE_LIFE)
    if !found {
        MyGlobal.SigninErrCacheLife = DFT_CACHE_SIGNIN_ERROR_LIFE
    } else {
        MyGlobal.SigninErrCacheLife, err = time.ParseDuration(sessLifeSignin)
        if nil != err {
            configError(CONFIG_SIGNIN_CACHE_LIFE, fmt.Sprintf("failed to parse %s as duration", sessLifeSignin))
        }
    }
    MyGlobal.SigninErrLimit, found = revel.Config.Int(CONFIG_SIGNIN_ERROR_LIMIT)
    if !found {
        MyGlobal.SigninErrLimit = DFT_SIGNIN_ERR_LIMIT
    }

    // after some signin errors, user should enter captcha to continue signing in
    MyGlobal.SigninUseCaptchaErrorNum, found = revel.Config.Int(CONFIG_SIGNIN_USECAPTCHA)
    if !found {
        MyGlobal.SigninUseCaptchaErrorNum = DFT_SIGNIN_USE_CAPTCHA
    }

    // lock user signin after too many errors in this time
    signinBanTime, found := revel.Config.String(CONFIG_SIGNIN_BAN_TIME)
    if !found {
        MyGlobal.SigninBanTime = DFT_SIGNIN_BAN_TIME
    } else {
        MyGlobal.SigninBanTime, err = time.ParseDuration(signinBanTime)
        if nil != err {
            configError(CONFIG_SIGNIN_BAN_TIME, fmt.Sprintf("failed to parse user signin ban time %s", signinBanTime))
        }
    }

    // parse custom template paths and save the file content
    MyGlobal.TemplateConfirmMail, found = revel.Config.String(CONFIG_TEMPLATE_CONFIRM_EMAIL)
    if !found {
        configNotFound(CONFIG_TEMPLATE_CONFIRM_EMAIL)
    }
    MyGlobal.TemplateConfirmMail = filepath.Join(revel.BasePath, MyGlobal.TemplateConfirmMail)
    confirmTemplData, err := ioutil.ReadFile(MyGlobal.TemplateConfirmMail)
    if nil != err {
        configError(CONFIG_TEMPLATE_CONFIRM_EMAIL, fmt.Sprintf("failed to read confirmation template from %s: %s", MyGlobal.TemplateConfirmMail, err))
    }
    MyGlobal.TemplateConfirmMail = string(confirmTemplData)

    MyGlobal.TemplateResetPassMail, found = revel.Config.String(CONFIG_TEMPLATE_RESETPASS_EMAIL)
    if !found {
        configNotFound(CONFIG_TEMPLATE_RESETPASS_EMAIL)
    }
    MyGlobal.TemplateResetPassMail = filepath.Join(revel.BasePath, MyGlobal.TemplateResetPassMail)
    resetPassTemplData, err := ioutil.ReadFile(MyGlobal.TemplateResetPassMail)
    if nil != err {
        configError(CONFIG_TEMPLATE_RESETPASS_EMAIL, fmt.Sprintf("failed to read resetpass template from %s: %s", MyGlobal.TemplateResetPassMail, err))
    }
    MyGlobal.TemplateResetPassMail = string(resetPassTemplData)

    // parse redis config
    MyGlobal.RedisServerAddr, found = revel.Config.String(CONFIG_REDIS_SERVER_ADDR)
    if !found {
        configNotFound(CONFIG_REDIS_SERVER_ADDR)
    }
    MyGlobal.RedisPoolMaxIdle, found = revel.Config.Int(CONFIG_REDIS_POOL_MAXIDLE)
    if !found {
        MyGlobal.RedisPoolMaxIdle = 20
    }
    redisIdleTimeout, found := revel.Config.String(CONFIG_REDIS_IDLE_TIMEOUT)
    if !found {
        MyGlobal.RedisPoolIdleTimeout = DFT_REDIS_POOL_IDLE_TIMEOUT
    } else {
        MyGlobal.RedisPoolIdleTimeout, err = time.ParseDuration(redisIdleTimeout)
        if nil != err {
            configError(CONFIG_REDIS_IDLE_TIMEOUT, fmt.Sprintf("failed to parse %s: %s", redisIdleTimeout, err))
        }
    }
}

