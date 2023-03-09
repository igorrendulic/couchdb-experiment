package mime

type EmailParser interface {
	Parse(message []byte) (*Email, error)
}
