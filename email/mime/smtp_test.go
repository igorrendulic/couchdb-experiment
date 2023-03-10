package mime

import (
	"testing"

	"github.com/igorrendulic/couchdb-experiment/email/mime/parser"
)

type fakeParser struct{}

func (p *fakeParser) Parse([]byte) (*parser.MailReceived, error) {
	return &parser.MailReceived{
		NotificationType: "fake",
	}, nil
}

func TestRegisterParser(t *testing.T) {
	Register("fake1", &fakeParser{})
	if len(Parsers()) != 1 {
		t.Error("expected 1 parser, got", len(Parsers()))
	}
}

func TestCallParser(t *testing.T) {
	fp := &fakeParser{}
	Register("fake2", fp)
	p, err := GetParser("fake2")
	if err != nil {
		t.Error("expected no error, got", err)
	}
	em, emErr := p.Parse([]byte("fake"))
	if emErr != nil {
		t.Error("expected no error, got", emErr)
	}
	if string(em.NotificationType) != "fake" {
		t.Error("expected fake, got", string(em.NotificationType))
	}
}

func TestParserNotRegistered(t *testing.T) {
	fp := &fakeParser{}
	Register("fake3", fp)
	_, err := GetParser("fake41")
	if err == nil {
		t.Error("expected error, got fake41")
	}
}

func TestListParsers(t *testing.T) {
	fp := &fakeParser{}
	Register("fake4", fp)
	all := Parsers()
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

func TestUnregisterAllParsers(t *testing.T) {
	fp := &fakeParser{}
	Register("fake5", fp)
	unregisterAllParsers()
	if len(Parsers()) != 0 {
		t.Error("expected 0 parsers, got", len(Parsers()))
	}
}
