package mail

import (
    "fmt"
    "net/smtp"
    "crypto/tls"
    "github.com/robfig/revel"
    "github.com/mawenbao/wishome/app"
    "github.com/mawenbao/wishome/app/modules/common"
)

const (
    NEWLINE = "\r\n"
    MIME = "MIME-Version: 1.0" + NEWLINE + "Content-Type: %s; charset=\"UTF-8\""
    BASE64_ENCODING = "Content-Transfer-Encoding: base64"
    SUBJECT = "Subject: %s"
    TO = "To: %s"
    MAIL_HEADER = MIME + NEWLINE + SUBJECT + NEWLINE + TO
    MAIL_FORMAT = MAIL_HEADER + NEWLINE + NEWLINE + "%s" // content-type + subject + to + body
    MAIL_FORMAT_BASE64 = BASE64_ENCODING + NEWLINE + MAIL_FORMAT
)

func SendTextMail(to, subject string, content []byte) bool {
    mail := fmt.Sprintf(
        MAIL_FORMAT,
        "text/plain",
        subject,
        to,
        content,
    )
    return SendRawMail(to, []byte(mail))
}

func SendTextMailBase64(to, subject string, content []byte) bool {
    mail := fmt.Sprintf(
        MAIL_FORMAT_BASE64,
        "text/plain",
        subject,
        to,
        common.EncodeBase64(content),
    )
    return SendRawMail(to, []byte(mail))
}

func SendHtmlMail(to, subject string, content []byte) bool {
    mail := fmt.Sprintf(
        MAIL_FORMAT,
        "text/html",
        subject,
        to,
        content,
    )
    return SendRawMail(to, []byte(mail))
}

func SendHtmlMailBase64(to, subject string, content []byte) bool {
    mail := fmt.Sprintf(
        MAIL_FORMAT_BASE64,
        "text/html",
        subject,
        to,
        common.EncodeBase64(content),
    )
    return SendRawMail(to, []byte(mail))
}

func SendRawMail(to string, mailData []byte) bool {
    mailClient, err := smtp.Dial(app.MyGlobal.String(app.CONFIG_MAIL_SMTP_ADDR))
    if nil != err {
        revel.ERROR.Printf("failed to connect to smtp server %s: %s", app.MyGlobal.String(app.CONFIG_MAIL_SMTP_ADDR), err)
        return false
    }
    defer mailClient.Close()

    // check STARTTLS support
    if ok, _ := mailClient.Extension("STARTTLS"); ok {
        // avoid error: x509: certificate signed by unknown authority
        tlc := &tls.Config{
            InsecureSkipVerify: true,
            ServerName: app.MyGlobal.String(app.CONFIG_MAIL_SMTP_HOST),
        }
        if err = mailClient.StartTLS(tlc); nil != err {
            revel.ERROR.Printf("failed to start tls: %s", err)
            return false
        }
    }

    mailClient.Mail(app.MyGlobal.String(app.CONFIG_MAIL_SENDER))
    mailClient.Rcpt(to)

    wc, err := mailClient.Data()
    if nil != err {
        revel.ERROR.Printf("failed to open mail client writer: %s", err)
        return false
    }
    defer wc.Close()

    _, err = wc.Write(mailData)
    if nil != err {
        revel.ERROR.Printf("failed to write mail content: %s", err)
        return false
    }

    return true
}

