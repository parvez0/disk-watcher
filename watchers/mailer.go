package watchers

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"os"
	"strings"
)

type MailTemplate struct {
	Subject  string
	HtmlBody string
	TextBody string
}

func getMailTemplate(template string, body string) MailTemplate {
	var mailTemplate MailTemplate
	switch template {
	case "whatsapp-disk":
		mailTemplate.Subject = "Whatsapp disk usage"
		bodyLine1 := body
		bodyLine2 := "to check the latest info of whatsapp disk"
		approvalLink := "http://localhost/disk/metrix"
		name := "Admins"
		mailTemplate.HtmlBody = "<html><head><link href= \"https://fonts.googleapis.com/css?family=Crimson+Text\" rel= \"stylesheet\"><style>.btn{text-decoration:none;border:1px solid #4384f5;padding:4px;border-radius:4px;color:#4384f5;font-size:large;font-family:'Helvetica Neue',Helvetica,Arial,sans-serif}.btn:hover{background:#4384f5;color:ghostwhite;font-family:'Helvetica Neue',Helvetica,Arial,sans-serif}</style></head><body style= \"font-family: 'Crimson Text', serif; text-decoration-line: none;\"><table cellpadding= \"0\" cellspacing= \"0\" border= \"0\" width= \"100%\" style= \"background: #f5f8fa; min-width: 350px; font-size: 1px; line-height: normal;\"><tr><td align= \"center\" valign= \"top\"><table cellpadding= \"0\" cellspacing= \"0\" border= \"0\" width= \"750\" class= \"table750\" style= \"width: 100%; max-width: 750px; min-width: 350px; background: #f5f8fa;\"><tr><td class= \"mob_pad\" width= \"25\" style= \"width: 25px; max-width: 25px; min-width: 25px;\">&nbsp;</td><td align= \"center\" valign= \"top\" style= \"background: #ffffff;\"><table cellpadding= \"0\" cellspacing= \"0\" border= \"0\" width= \"100%\" style= \"width: 100% !important; min-width: 100%; max-width: 100%; background: #f5f8fa;\"><tr><td align= \"right\" valign= \"top\"><div class= \"top_pad\" style= \"height: 25px; line-height: 25px; font-size: 23px;\">&nbsp;</div></td></tr></table><table cellpadding= \"0\" cellspacing= \"0\" border= \"0\" width= \"88%\" style= \"width: 88% !important; min-width: 88%; max-width: 88%;\"><tr><td align= \"center\" valign= \"top\"><div style= \"height: 40px; line-height: 40px; font-size: 38px;\">&nbsp;</div> <a href= \"https://yellowmessenger.com\" target= \"_blank\" style= \"display: block; max-width: 192px;\"> <img src= \"https://cdn.yellowmessenger.com/LyG9qYl4ZHFs1573044356873.jpg\" alt= \"Yellow Messenger\" width= \"192\" border= \"0\" style= \"display: block; width: 192px;\" /> </a> <span style= \"color: #797979; margin-top:0px;font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;font-size: 10px\">Conversational AI for 10x Enterprise.</span><div class= \"top_pad2\" style= \"height: 25px; line-height: 48px; font-size: 46px;\">&nbsp;</div></td></tr></table><table cellpadding= \"0\" cellspacing= \"0\" border= \"0\" width= \"88%\" style= \"width: 88% !important; min-width: 88%; max-width: 88%; text-align: left;align-items: center;\"><tr><td style= \"align-content: left;\"><div style= \"font-family: 'Crimson Text', serif;align-content:center;color: #5a5a5a;font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;font-size: small;text-decoration-line: none;margin-bottom: 30px; \"><p style= \"margin: 10px;\">Hi " + name + ",</p><p style= \"margin: 10px;\">" + bodyLine1 + "</p><p style= \"margin: 10px;\">Please <a href=" + approvalLink + " target= \"_blank\" style= \"text-decoration-line: none;\">click here</a> " + bodyLine2 + " </p></div></td></tr><tr><td>&nbsp;</td></tr><tr><td>&nbsp;</td></tr></table><table cellpadding= \"0\" cellspacing= \"0\" border= \"0\" width= \"90%\" style= \"width: 90% !important; margin-top: 20px;min-width: 90%; max-width: 90%; border-width: 1px; border-style: solid; border-color: #e8e8e8; border-bottom: none; border-left: none; border-right: none;text-align: center\"><tr><td align= \"left\" valign= \"top\"><div style= \"height: 28px; line-height: 28px; font-size: 26px;\">&nbsp;</div></td></tr></table><table cellpadding= \"0\" cellspacing= \"0\" border= \"0\" width= \"88%\" style= \"width: 88% !important; min-width: 88%; max-width: 88%;\"><tr><td align= \"left\" valign= \"top\" style= \"text-align: center\"><div style= \"height: 30px; line-height: 30px; font-size: 28px;\">&nbsp;</div></td></tr></table><table cellpadding= \"0\" cellspacing= \"0\" border= \"0\" width= \"100%\" style= \"width: 100% !important; min-width: 100%; max-width: 100%; background: #f5f8fa;\"><tbody><tr><td align= \"center\" valign= \"top\"><div style= \"margin-top:15px;text-align: center\"><p face= \"'Source Sans Pro', sans-serif\" color= \"#868686\"> <span style= \"font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; color: #868686; font-size: 11px; line-height: 20px;\"> Copyright &copy; Bitonic Technology Labs, Inc. Silverside road, Wilmington, Delaware 02139.</span></p></div></td></tr></tbody></table></td><td class= \"mob_pad\" width= \"25\" style= \"width: 25px; max-width: 25px; min-width: 25px;\">&nbsp;</td></tr></table><div style= \"height:25px;line-height:25px;font-size:23px\">&nbsp;</div></td></tr></table></body></html>"
		mailTemplate.TextBody = "Required by sendgrid"
	}
	return mailTemplate
}

func Send(sendTo string, mailTemplate MailTemplate) {
	from := mail.NewEmail("Whatsapp Disk", "alert@yellowmessenger.com")
	to := mail.NewEmail("", sendTo)
	message := mail.NewSingleEmail(from, mailTemplate.Subject, to, mailTemplate.TextBody, mailTemplate.HtmlBody)
	logger.Info("Sending mail for " + mailTemplate.Subject + " to " + sendTo)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		logger.Error("failed to send mail to - ", sendTo, " error -", err)
	} else {
		logger.Info("mail send to "+sendTo+" with status code -", response.StatusCode, " body - ", response.Body)
	}
}

func SendMail(template string, body string, sendTo []string) {
	mailTemplate := getMailTemplate(template, body)
	if sendTo == nil || len(sendTo) == 0 {
		sendTo = strings.Split(os.Getenv("INFRA_MAILS_IDS"), ",")
		if len(sendTo) == 0 {
			logger.Error("Mail ID's not provided, for mail notification")
			return
		}
	}
	for _, email := range sendTo {
		Send(email, mailTemplate)
	}
}
