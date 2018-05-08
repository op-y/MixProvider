package email

import (
    "log"
    "net/smtp"
    "strconv"
    "strings"
)

// Send email through SMTP
func SendEmail(username, password, address, tos, subject, body, mailType string) error {
    endpoint := strings.Split(address, ":")
    auth := smtp.PlainAuth("", username, password, endpoint[0])

    var contentType string
    if mailType == "html" {
        contentType = "Content-Type: text/" + mailType + "; charset=UTF-8"
    } else {
        contentType = "Content-Type: text/plain" + "; charset=UTF-8"
    }

    msg := []byte("To: " + tos + "\r\nFrom: " + username + ">\r\nSubject: " + "\r\n" + contentType + "\r\n\r\n" + body)
    sendTo := strings.Split(tos, ";")
    if err := smtp.SendMail(address, auth, username, sendTo, msg); err != nil {
        log.Printf("fail to send email: %s", err.Error())
        return err
    }
    log.Printf("send email<%s> to <%s> successfully", tos, body)
    return nil
}

// The entrance of email
func EmailGo(username string, password string, ip string, port int, tos string, subject string, content string) error {
    address := ip + ":" + strconv.Itoa(port)
    body := "<html><body><h3>" + content + "</h3></body></html>"
     if err := SendEmail(username, password, address, tos, subject, body, "html"); err != nil {
         log.Println("fail to send email")
         return err
     }
     return nil
}

