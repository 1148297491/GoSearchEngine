package dir

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

func UpdateName(c *gin.Context) {
	u, is := c.Get("user")
	if !is {
		c.JSON(http.StatusOK, dto.UpdateDirNameResponse{
			Code: common.UserHasDeleted,
		})
		return
	}
	var userId = strconv.FormatInt(u.(models.User).UserID, 10)
	var s dto.UpdateDirNameRequest
	if err := c.ShouldBind(&s); err == nil {
		var dir models.Dir
		//判断文件夹是不是存在同时判断文件夹是不是用户自己的（只能操作属于自己的文件夹）
		result := middleware.Db.Where("user_id = ? and dir_id = ?", userId, s.DirId).First(&dir)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, dto.UpdateDirNameResponse{
					Code: common.DirNotExisted,
				})
				fmt.Println("UpdateDirName err:", result.Error)
				return
			} else {
				c.JSON(http.StatusOK, dto.UpdateDirNameResponse{
					Code: common.UnknownError,
				})
				fmt.Println("UpdateDirName err:", result.Error)
				return
			}
		}
		dir.DirName = s.NewDirName
		result = middleware.Db.Save(&dir)
		if result.Error != nil {
			c.JSON(http.StatusOK, dto.UpdateDirNameResponse{
				Code: common.UnknownError,
			})
			fmt.Println("UpdateDirName err: ", err)
			return
		}
		c.JSON(http.StatusOK, dto.UpdateDirNameResponse{
			Code: common.OK,
		})
		return
	} else {
		//参数校验失败的报错
		fmt.Println("bind err: ", err)
		c.JSON(http.StatusOK, dto.UpdateDirNameResponse{Code: common.ParamInvalid})
	}
}
