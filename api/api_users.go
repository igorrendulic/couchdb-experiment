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

type UserApi struct {
	userService services.CouchDB
}

func NewUserAPI(userService services.CouchDB) *UserApi {
	return &UserApi{
		userService: userService,
	}
}

// Regiser new user godoc
// @Summary REgister user
// @Description Regiter user
// @Tags Login and Registration
// @Param user body models.RegisterInput true "Name and Password required"
// @Success 200 {object} models.User
// @Failure 401 {object} w3kitgonic.JSONError "login failed"
// @Failure 403 {object} w3kitgonic.JSONError "login forbidden"
// @Failure 400 {object} w3kitgonic.JSONError "can't login (no partner association)"
// @Accept json
// @Produce json
// @Router /api/v1/register [post]
func (ua *UserApi) RegisterUser(c *gin.Context) {
	var user models.RegisterInput
	if err := c.ShouldBindWith(&user, binding.JSON); err != nil {
		global.Logger.Log("missing required fields", err)
		w3kitgonic.AbortWithError(c, http.StatusBadRequest, "invalid json format or missing required field")
		return
	}

	out, err := ua.userService.RegisterUser(user.Email, user.Password)
	if err != nil {
		global.Logger.Log("failed to register user", err)
		w3kitgonic.AbortWithError(c, http.StatusInternalServerError, "failed to register user")
		return
	}
	c.JSON(http.StatusOK, out)
}
