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

func NewDir(c *gin.Context) {
	u, is := c.Get("user")
	if !is {
		c.JSON(http.StatusOK, dto.NewDirResponse{
			Code: common.UserHasDeleted,
		})
		return
	}
	var userId = strconv.FormatInt(u.(models.User).UserID, 10)
	var s dto.NewDirRequest
	if err := c.ShouldBind(&s); err == nil {
		var dir models.Dir
		result := middleware.Db.Where("dir_name = ? and user_id = ?", s.DirName, userId).First(&dir)
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, dto.NewDirResponse{
				Code: common.DirHasExisted,
				Data: struct{ DirId string }{DirId: strconv.FormatInt(dir.DirId, 10)},
			})
			return
		}

		dir = models.Dir{UserId: u.(models.User).UserID, DirName: s.DirName}
		result = middleware.Db.Create(&dir)
		middleware.Db.Model(&dir).Update("dir_id", int64(dir.ID))
		if result.Error != nil {
			c.JSON(http.StatusOK, dto.NewDirResponse{
				Code: common.UnknownError,
			})
			fmt.Println("NewDirErr: ", result.Error)
			return
		}

		c.JSON(http.StatusOK, dto.NewDirResponse{
			Code: common.OK,
			Data: struct{ DirId string }{DirId: strconv.FormatInt(dir.DirId, 10)},
		})
		return
	} else {
		//参数校验失败的报错
		fmt.Println("bind err: ", err)
		c.JSON(http.StatusOK, dto.NewDirResponse{Code: common.ParamInvalid})
	}
}
