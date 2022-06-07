package dir

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"GoSearchEngine/common"
	"GoSearchEngine/dto"
	"GoSearchEngine/middleware"
	"GoSearchEngine/models"
	"net/http"
	"strconv"
)

func GetDir(c *gin.Context) {
	u, is := c.Get("user")
	if !is {
		c.JSON(http.StatusOK, dto.GetDirResponse{
			Code: common.UserHasDeleted,
		})
		return
	}
	var userId = strconv.FormatInt(u.(models.User).UserID, 10)
	//判断文件夹是不是存在同时判断文件夹是不是用户自己的（只能操作属于自己的文件夹）
	var dirs []models.Dir
	result := middleware.Db.Where("user_id = ?", userId).Find(&dirs)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, dto.GetDirResponse{
				Code: common.UserHaveNoDir,
			})
			fmt.Println("GetUserDir err:", result.Error)
			return
		} else {
			c.JSON(http.StatusOK, dto.GetDirResponse{
				Code: common.UnknownError,
			})
			fmt.Println("GetUserDir err:", result.Error)
			return
		}
	}
	c.JSON(http.StatusOK, dto.GetDirResponse{
		Code: common.OK,
		Data: struct{ DirList []models.Dir }{DirList: dirs},
	})
	return
}
