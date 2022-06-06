package middleware

import (
	"GoSearchEngine/common"
	"GoSearchEngine/dto"
	"GoSearchEngine/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func CookieMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("camp-session")
		if err != nil {
			c.JSON(http.StatusOK, dto.LogoutResponse{Code: common.LoginRequired})
			c.Abort()
			return
		}
		var u models.User
		//db := *(c.MustGet("db").(**gorm.db))
		result := Db.Where("user_id = ?", cookie).First(&u)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.Set("deletedUser", u)
			return
		}
		c.Set("user", u)
	}
}
