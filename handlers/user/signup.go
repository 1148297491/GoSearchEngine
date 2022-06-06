package user

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gofound/common"
	"gofound/dto"
	"gofound/middleware"
	"gofound/models"
	"net/http"
	"strconv"
)

func Signup(c *gin.Context) {
	var r dto.SignUpRequest
	if err := c.ShouldBind(&r); err == nil {

		//参数校验通过，获取用户
		u := models.User{Userphone: r.Userphone, Password: r.Password}
		var uu models.User

		//查找是否有用户名重复的用户，若有则报错
		result := middleware.Db.Where("userphone = ?", u.Userphone).First(&uu)
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, dto.CreateUserResponse{
				Code: common.UserHasExisted,
				Data: struct {
					UserID string
				}{UserID: strconv.FormatInt(uu.UserID, 10)},
			})
			return
		}

		//创建新用户
		middleware.Db.Create(&u)
		//rd.Set(ctx, USER_REDIS_PREFIX+strconv.FormatInt(u.UserID, 10), 1, 0)
		middleware.Db.Model(&u).Update("user_id", int64(u.ID))
		c.SetCookie("camp-session", strconv.FormatInt(u.UserID, 10), 3600, "/", "127.0.0.1", false, true)
		c.JSON(http.StatusOK, dto.SignUpResponse{
			Code: common.OK,
			Data: struct {
				UserID string
			}{UserID: strconv.FormatInt(u.UserID, 10)},
		})
	} else {
		//参数校验失败的报错
		fmt.Println("bind err: ", err)
		c.JSON(http.StatusOK, dto.CreateUserResponse{Code: common.ParamInvalid})
	}
}
