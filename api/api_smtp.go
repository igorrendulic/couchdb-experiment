package api

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	awssmtphandler "github.com/igorrendulic/couchdb-email-aws-parse"
	"github.com/igorrendulic/couchdb-experiment/email/mime"
	"github.com/igorrendulic/couchdb-experiment/email/mime/parser"
	w3kitgonic "github.com/mailio/go-web3-kit/gingonic"
)

type ApiSmtp struct {
	smtpHandler parser.Parser
}

func NewApiSmtp() *ApiSmtp {
	mime.Register("aws", awssmtphandler.NewAwsParser())
	p, err := mime.GetSmtpHandler("aws")
	if err != nil {
		panic(err.Error())
	}
	return &ApiSmtp{
		smtpHandler: p,
	}
}

// Get - a test link
func (ah *ApiSmtp) Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusNotFound,
		"message": "not found",
	})
}

// Receive webhook - entry point for all incoming emails
// @Summary Smtp Receive webhook
// @Description Receive webhook - entry point for all incoming emails
// @Description check how to handle notifications: https://aws.amazon.com/premiumsupport/knowledge-center/ses-bounce-notifications-sns/
// @Description additional docs: https://docs.aws.amazon.com/ses/latest/dg/notification-contents.html#bounce-object
// @Description https://docs.aws.amazon.com/ses/latest/dg/receiving-email-concepts.html
// @Description https://docs.aws.amazon.com/ses/latest/dg/receiving-email-setting-up.html
// @Tags Emails
// @Param payload body sns.Payload true "receiving email"
// @Success 200 {object} parser.MailReceived
// @Failure 401 {object} w3kitgonic.JSONError "login failed"
// @Failure 400 {object} w3kitgonic.JSONError "can't login (no partner association)"
// @Accept json
// @Produce json
// @Router /api/smtp [post]
func (ah *ApiSmtp) SmtpWebhook(c *gin.Context) {

	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		w3kitgonic.AbortWithError(c, http.StatusBadRequest, "webhook payload is not valid")
		return
	}

	fmt.Printf("json: %s", string(jsonData))
	parsedPayload, pErr := ah.smtpHandler.Parse(jsonData)
	if pErr != nil {
		w3kitgonic.AbortWithError(c, http.StatusBadRequest, "SNS aws payload can't be parsed")
		return
	}
	c.JSON(http.StatusOK, parsedPayload)
}
