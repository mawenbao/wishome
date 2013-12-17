package app

import (
    "time"
)

const (
    STR_USER = "user"
    STR_NAME = "name"
    STR_PASSWORD = "password"
    STR_EMAIL = "email"
    STR_KEY = "key"
    STR_EXPIRE = "expire"
    STR_LASTUSER = "lastuser"

    CACHE_RESETPASS = "reset_pass"
    CACHE_SIGNUP_CONFIRM = "signup_confirm"
    CACHE_SIGNIN_ERROR = "signin_err_sess"

    TABLE_USERS = "users"

    CONFIG_APP_URL = "app.url"
    CONFIG_RESETPASS_KEY_LEN = "user.resetpass.keylen"
    CONFIG_RESETPASS_KEY_LIFE = "user.resetpass.keylife"
    CONFIG_SIGNUP_KEY_LEN = "user.signup.keylen"
    CONFIG_SIGNUP_KEY_LIFE = "user.signup.keylife"
    CONFIG_MAIL_SMTP_ADDR = "mail.smtp.address"
    CONFIG_MAIL_SMTP_HOST = "mail.smtp.host"
    CONFIG_MAIL_SMTP_PORT = "mail.smtp.port"
    CONFIG_MAIL_SENDER = "mail.sender"
    CONFIG_DB_DRIVER = "db.driver"
    CONFIG_DB_SPEC = "db.spec"
    CONFIG_CAPTCHA_WIDTH = "captcha.width"
    CONFIG_CAPTCHA_HEIGHT = "captcha.height"
    CONFIG_CAPTCHA_LENGTH = "captcha.length"
    // client side cookie life
    CONFIG_SESSION_LIFE = "session.life"
    // server side session cache
    CONFIG_SIGNIN_CACHE_LIFE = "session.signin.life"
    CONFIG_SIGNIN_USECAPTCHA = "user.signin.usecaptcha"
    CONFIG_SIGNIN_ERROR_LIMIT = "session.error.limit"
)

// default settings
var (
    DEFAULT_TIME_FORMAT = time.ANSIC
)

