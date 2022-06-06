package dir

import (
	"GoSearchEngine/common"
	"GoSearchEngine/dto"
	"GoSearchEngine/middleware"
	"GoSearchEngine/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func DeleteDir(c *gin.Context) {
	u, is := c.Get("user")
	if !is {
		c.JSON(http.StatusOK, dto.DeleteDirResponse{
			Code: common.UserHasDeleted,
		})
		return
	}
	var userId = strconv.FormatInt(u.(models.User).UserID, 10)
	var s dto.DeleteDirRequest
	if err := c.ShouldBind(&s); err == nil {
		var dir models.Dir
		//判断文件夹是不是存在同时判断文件夹是不是用户自己的（只能操作属于自己的文件夹
		result := middleware.Db.Where("user_id = ? and dir_id = ?", userId, s.DirId).First(&dir)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, dto.DeleteUserResponse{
					Code: common.DirNotExisted,
				})
				fmt.Println("DeleteDir err:", result.Error)
				return
			} else {
				c.JSON(http.StatusOK, dto.DeleteDirResponse{
					Code: common.UnknownError,
				})
				fmt.Println("DeleteDir err:", result.Error)
				return
			}
		}

		middleware.Db.Delete(&dir)
		c.JSON(http.StatusOK, dto.DeleteDirResponse{
			Code: common.OK,
		})
		return
	} else {
		//参数校验失败的报错
		fmt.Println("bind err: ", err)
		c.JSON(http.StatusOK, dto.DeleteDirResponse{Code: common.ParamInvalid})
	}
}
