package user

import (
	"GoSearchEngine/common"
	"GoSearchEngine/dto"
	"GoSearchEngine/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	u, is := c.Get("user")
	if !is {
		c.JSON(http.StatusOK, dto.LogoutResponse{Code: common.UserHasDeleted})
		return
	}
	c.SetCookie("camp-session", strconv.FormatInt(u.(models.User).UserID, 10), -1, "/", "127.0.0.1", false, true)
	c.JSON(http.StatusOK, dto.LogoutResponse{Code: common.OK})
}
