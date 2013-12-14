package mail

import (
    "fmt"
    "net/smtp"
    "crypto/tls"
    "github.com/robfig/revel"
    "github.com/mawenbao/wishome/app"
)

func SendMail(to, subject string, content []byte) bool {
    mailClient, err := smtp.Dial(app.MyGlobal[app.CONFIG_MAIL_SMTP_ADDR].(string))
    if nil != err {
        revel.ERROR.Printf("failed to connect to smtp server %s: %s", app.MyGlobal[app.CONFIG_MAIL_SMTP_ADDR].(string), err)
        return false
    }
    defer mailClient.Close()

    // check STARTTLS support
    if ok, _ := mailClient.Extension("STARTTLS"); ok {
        // avoid error: x509: certificate signed by unknown authority
        tlc := &tls.Config{
            InsecureSkipVerify: true,
            ServerName: app.MyGlobal[app.CONFIG_MAIL_SMTP_HOST].(string),
        }
        if err = mailClient.StartTLS(tlc); nil != err {
            revel.ERROR.Printf("failed to start tls: %s", err)
            return false
        }
    }

    mailClient.Mail(app.MyGlobal[app.CONFIG_MAIL_SENDER].(string))
    mailClient.Rcpt(to)

    wc, err := mailClient.Data()
    if nil != err {
        revel.ERROR.Printf("failed to open mail client writer: %s", err)
        return false
    }
    defer wc.Close()

    mailBody := fmt.Sprintf(
        "To: %s\r\nSubject: %s\r\n\r\n%s",
        to,
        subject,
        content,
    )

    _, err = wc.Write([]byte(mailBody))
    if nil != err {
        revel.ERROR.Printf("failed to write mail content: %s", err)
        return false
    }

    return true
}

