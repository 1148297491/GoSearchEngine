package user

import (
	"github.com/gin-gonic/gin"
	"GoSearchEngine/common"
	"GoSearchEngine/dto"
	"GoSearchEngine/models"
	"net/http"
	"strconv"
)

func Whoami(c *gin.Context) {
	u, is := c.Get("user")
	if !is {
		du := c.MustGet("deletedUser").(models.User)
		c.JSON(http.StatusOK, dto.WhoAmIResponse{
			Code: common.UserHasDeleted,
			Data: models.TUser{
				UserID:    strconv.FormatInt(du.UserID, 10),
				Userphone: du.Userphone,
			},
		})
		return
	}
	c.JSON(http.StatusOK, dto.WhoAmIResponse{
		Code: common.OK,
		Data: models.TUser{
			UserID:    strconv.FormatInt(u.(models.User).UserID, 10),
			Userphone: u.(models.User).Userphone,
		},
	})
}
