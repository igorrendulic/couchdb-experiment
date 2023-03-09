package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/mail"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/go-resty/resty/v2"
	"github.com/igorrendulic/couchdb-experiment/global"
	"github.com/igorrendulic/couchdb-experiment/models"
	"github.com/igorrendulic/couchdb-experiment/services"
	"github.com/igorrendulic/couchdb-experiment/utils"
)

// couchDB interface implementation
type couchDB struct {
	restClient    *resty.Client
	snowflakeNode *snowflake.Node
}

func NewCouchDB() services.CouchDB {
	client := resty.New()
	client.SetHostURL(fmt.Sprintf("%s://%s:%d", global.Conf.CouchDB.Scheme, global.Conf.CouchDB.Host, global.Conf.CouchDB.Port))
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Accept", "application/json")
	client.SetHeader("User-Agent", "go-web3-kit/1.0.0")
	client.SetBasicAuth(global.Conf.CouchDB.Username, global.Conf.CouchDB.Password)

	snfl, snErr := snowflake.NewNode(1)
	if snErr != nil {
		log.Fatalf("failed to create snowflake node: %v", snErr)
	}
	return &couchDB{
		restClient:    client,
		snowflakeNode: snfl,
	}
}

// Get returns the value for the given key.
func (c *couchDB) Get(database, key string) ([]byte, error) {
	resp, err := c.restClient.R().Get(fmt.Sprintf("/%s/%s", database, key))
	return resp.Body(), err
}

// Set sets the value for the given key.
func (c *couchDB) Set(user, database, key, value string) error {
	return nil
}

// Delete deletes the value for the given key.
func (c *couchDB) Delete(user, database, key string) error {
	return nil
}

// Register new user
func (c *couchDB) RegisterUser(email, password string) (*models.User, error) {
	emailAddr, err := mail.ParseAddress(email)
	if err != nil {
		global.Logger.Log(err, "Invalid email address")
		return nil, err
	}
	user := strings.Split(emailAddr.Address, "@")[0]
	resp, err := c.restClient.R().SetBody(map[string]interface{}{"name": user, "password": password, "roles": []string{}, "type": "user"}).Put(fmt.Sprintf("/_users/org.couchdb.user:%s", user))
	if err != nil {
		global.Logger.Log(err, "Failed to register user")
		return nil, err
	}

	hErr := handleError(resp.Body())
	if hErr != nil {
		return nil, hErr
	}

	// create index on default user database by folder field

	hexUser := "userdb-" + utils.HexEncodeToString(user)

	// wait for database to be created
	for i := 1; i < 5; i++ {
		resp, err = c.restClient.SetRetryCount(5).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(5 * time.Second).R().Get(fmt.Sprintf("/%s", hexUser))
		hErr = handleError(resp.Body())
		if hErr != nil {
			if hErr.Error() == "not_found" {
				backoff := int(100 * math.Pow(2, float64(i)))
				fmt.Printf("backoff: %d", backoff)
				time.Sleep(time.Duration(backoff) * time.Millisecond)
				continue
			} else {
				return nil, hErr
			}
		}
	}

	// create index on database
	folderIndex := map[string]interface{}{
		"index": map[string]interface{}{
			"fields": []map[string]interface{}{{"folder": "desc"}, {"created": "desc"}},
		},
		"name": "folder-index",
		"type": "json",
		"ddoc": "folder-index",
	}
	resp, err = c.restClient.R().SetBody(folderIndex).Post(fmt.Sprintf("/%s/_index", hexUser))
	hErr = handleError(resp.Body())
	if hErr != nil {
		return nil, hErr
	}

	// registration automatically creates databased userdb-{f(username) = hex}
	u := &models.User{
		Email:    emailAddr.Address,
		Username: user,
	}

	return u, err
}

// AddEmail adds new email to database
func (c *couchDB) AddEmail(username string, email *models.Email) (*models.Email, error) {

	hexUser := "userdb-" + utils.HexEncodeToString(username)
	fmt.Printf("hex user: %s", hexUser)

	id := c.snowflakeNode.Generate().String()
	email.Id = id
	email.Created = utils.GetTimestamp()

	resp, err := c.restClient.R().SetBody(email).Put(fmt.Sprintf("/%s/%s", hexUser, email.Id))
	if err != nil {
		return nil, err
	}
	var response map[string]interface{}
	uErr := json.Unmarshal(resp.Body(), &response)
	if uErr != nil {
		global.Logger.Log(uErr, "Failed to unmarshal response")
		return nil, uErr
	}
	if response["error"] != nil {
		global.Logger.Log(response["error"])
		return nil, errors.New(response["error"].(string))
	}

	return email, nil
}

// func (c *couchDB) ListEmails(username string, folder string, limit int) ([]*models.Email, error) {
// 	hexUser := "userdb-" + utils.HexEncodeToString(username)
// 	query := map[string]interface{}{
// 		"selector": map[string]interface{}{
// 			"folder": map[string]interface{}{
// 				"$eq": folder,
// 			},
// 			"limit": 21,
// 			"skip":  0,
// 		},
// 	}
// 	resp, err := c.restClient.R().SetBody(query).Post(fmt.Sprintf("/%s/_find", hexUser))
// 	if err != nil {
// 		return nil, err
// 	}
// 	hErr := handleError(resp.Body())
// 	if hErr != nil {
// 		return nil, hErr
// 	}
// 	output := resp.Body()

// 	var emails []*models.Email
// 	uErr := json.Unmarshal(output, &emails)
// 	if uErr != nil {
// 		return nil, uErr
// 	}
// 	return emails, nil
// }
