package parser

type Parser interface {
	Parse([]byte) (*MailReceived, error)
}
