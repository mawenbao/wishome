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

    CACHE_RESETPASS = "reset_pass_"
    CACHE_SIGNUP_CONFIRM = "signup_confirm_"

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
    CONFIG_SESSION_LIFE = "session.life"
    CONFIG_DB_DRIVER = "db.driver"
    CONFIG_DB_SPEC = "db.spec"
)

// default settings
var (
    DEFAULT_TIME_FORMAT = time.ANSIC
)

