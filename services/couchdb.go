package services

import (
	"github.com/igorrendulic/couchdb-experiment/models"
)

type CouchDB interface {
	// Get returns the value for the given key.
	Get(database, key string) ([]byte, error)
	// Set sets the value for the given key.
	Set(user, database, key, value string) error
	// Delete deletes the value for the given key.
	Delete(user, database, key string) error

	RegisterUser(username, password string) (*models.User, error)

	AddEmail(username string, email *models.Email) (*models.Email, error)

	// ListEmails(username string, folder string, limit int) ([]*models.Email, error)
}
