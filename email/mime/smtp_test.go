package mime

import (
	"testing"

	"github.com/igorrendulic/couchdb-experiment/email/mime/handler"
)

type fakeHandler struct{}

func (p *fakeHandler) HandleSmtp([]byte) (*handler.MailReceived, error) {
	return &handler.MailReceived{
		NotificationType: "fake",
	}, nil
}

func TestRegisterSmtpHandler(t *testing.T) {
	Register("fake1", &fakeHandler{})
	if len(SmtpHandlers()) != 1 {
		t.Error("expected 1 parser, got", len(SmtpHandlers()))
	}
}

func TestCalSmtpHandler(t *testing.T) {
	fp := &fakeHandler{}
	Register("fake2", fp)
	p, err := GetSmtpHandler("fake2")
	if err != nil {
		t.Error("expected no error, got", err)
	}
	em, emErr := p.HandleSmtp([]byte("fake"))
	if emErr != nil {
		t.Error("expected no error, got", emErr)
	}
	if string(em.NotificationType) != "fake" {
		t.Error("expected fake, got", string(em.NotificationType))
	}
}

func TestSmtpHandlerNotRegistered(t *testing.T) {
	fp := &fakeHandler{}
	Register("fake3", fp)
	_, err := GetSmtpHandler("fake41")
	if err == nil {
		t.Error("expected error, got fake41")
	}
}

func TestListSmtpHandlers(t *testing.T) {
	fp := &fakeHandler{}
	Register("fake4", fp)
	all := SmtpHandlers()
	found := false
	for a := range all {
		if all[a] == "fake4" {
			found = true
		}
	}
	if !found {
		t.Error("expected to find fake4, got", all)
	}
}

func TestUnregisterAllSmtpHandlers(t *testing.T) {
	fp := &fakeHandler{}
	Register("fake5", fp)
	unregisterAllSmtpHandlers()
	if len(SmtpHandlers()) != 0 {
		t.Error("expected 0 handlers, got", len(SmtpHandlers()))
	}
}
