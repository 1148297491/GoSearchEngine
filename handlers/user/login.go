package user

import (
	"GoSearchEngine/common"
	"GoSearchEngine/dto"
	"GoSearchEngine/middleware"
	"GoSearchEngine/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Login(c *gin.Context) {
	var r dto.LoginRequest
	if err := c.ShouldBind(&r); err == nil {
		var u models.User
		//db := *(c.MustGet("db").(**gorm.db))
		result := middleware.Db.Unscoped().Where("userphone= ?", r.Userphone).First(&u)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, dto.LoginResponse{Code: common.WrongPassword})
			return
		}
		var du models.User
		result = middleware.Db.Where("userphone= ?", r.Userphone).First(&du)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, dto.LoginResponse{
				Code: common.WrongPassword,
				Data: struct {
					UserID string
				}{UserID: strconv.FormatInt(u.UserID, 10)},
			})
			return
		}
		if du.Password != r.Password {
			c.JSON(http.StatusOK, dto.LoginResponse{
				Code: common.WrongPassword,
				Data: struct {
					UserID string
				}{UserID: strconv.FormatInt(du.UserID, 10)},
			})
			return
		}
		c.SetCookie("camp-session", strconv.FormatInt(du.UserID, 10), 3600, "/", "127.0.0.1", false, false)
		c.JSON(http.StatusOK, dto.LoginResponse{
			Code: common.OK,
			Data: struct {
				UserID string
			}{UserID: strconv.FormatInt(du.UserID, 10)},
		})
	} else {
		c.JSON(http.StatusOK, dto.LoginResponse{Code: common.ParamInvalid})
	}
}
