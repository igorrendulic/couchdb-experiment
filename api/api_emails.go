package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/igorrendulic/couchdb-experiment/global"
	"github.com/igorrendulic/couchdb-experiment/models"
	"github.com/igorrendulic/couchdb-experiment/services"
	w3kitgonic "github.com/mailio/go-web3-kit/gingonic"
)

type EmailApi struct {
	emailService services.CouchDB
}

func NewEmailAPI(emailService services.CouchDB) *EmailApi {
	return &EmailApi{
		emailService: emailService,
	}
}

// new email godoc
// @Summary New email received
// @Description Nwe email received (eithe rSMTP or mailio)
// @Tags Emails
// @Param user body models.Email true "receiving email"
// @Success 200 {object} models.Email
// @Failure 401 {object} w3kitgonic.JSONError "login failed"
// @Failure 400 {object} w3kitgonic.JSONError "can't login (no partner association)"
// @Accept json
// @Produce json
// @Router /api/v1/email [post]
func (ea *EmailApi) AddEmail(c *gin.Context) {
	var email models.Email
	if err := c.ShouldBindWith(&email, binding.JSON); err != nil {
		global.Logger.Log("missing required fields", err)
		w3kitgonic.AbortWithError(c, http.StatusBadRequest, "invalid json format or missing required field")
		return
	}

	out, err := ea.emailService.AddEmail(email.Username, &email)
	if err != nil {
		global.Logger.Log("failed to add email", err)
		w3kitgonic.AbortWithError(c, http.StatusInternalServerError, "failed to add email")
		return
	}
	c.JSON(http.StatusOK, out)
}

// // get emails godoc
// // @Summary get emails
// // @Description get emails
// // @Tags Emails
// // @Param username path string true "username"
// // @Param folder path string true "folder"
// // @Param limit path int false "default = 20"
// // @Success 200 {array} models.Email
// // @Failure 401 {object} w3kitgonic.JSONError "login failed"
// // @Failure 400 {object} w3kitgonic.JSONError "can't login (no partner association)"
// // @Accept json
// // @Produce json
// // @Router /api/v1/email [get]
// func (ea *EmailApi) ListEmails(C *gin.Context) {
// 	var emails []*models.Email
// 	var err error
// 	username := C.Param("username")
// 	folder := C.Param("folder")
// 	limit := C.Param("limit")

// 	limitInt := 20
// 	if limit != "" {
// 		limitInt = utils.StringToInt(limit)
// 	}

// 	emails, err = ea.emailService.ListEmails(username, folder, limitInt)

// 	if err != nil {
// 		global.Logger.Log("failed to list emails", err)
// 		w3kitgonic.AbortWithError(C, http.StatusInternalServerError, "failed to list emails")
// 		return
// 	}
// 	C.JSON(http.StatusOK, emails)
// }
