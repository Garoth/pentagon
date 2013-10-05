package mail

import (
    "log"
    "net/smtp"
    "encoding/json"


    "github.com/Garoth/pentagon-model"
)

func Start() chan string {
    comm := make(chan string, 10)

    go func() {
        for {
            message := <-comm

            mail := &pentagonmodel.MailComponentMessage{}
            if err := json.Unmarshal([]byte(message), &mail); err != nil {
                log.Println("Decoding Message:", message, "Error:", err)
                continue
            }

            body := "To: " + mail.To + "\r\nSubject: " +
               mail.Subject + "\r\n\r\n" + mail.Message
            conf := pentagonmodel.GetConfig()

            auth := smtp.PlainAuth("",
                conf.GmailAddress, conf.GmailPassword, "smtp.gmail.com")

            err := smtp.SendMail("smtp.gmail.com:587", auth, conf.GmailAddress,
               []string{mail.To}, []byte(body))
            if err != nil {
               log.Println("Couldn't send mail:", err)
            }

            log.Println("Sent email successfully")
        }
    }()

    return comm
}

