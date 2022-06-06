package collection

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

func Collect(c *gin.Context) {
	u, is := c.Get("user")
	if !is {
		c.JSON(http.StatusOK, dto.CollectResponse{
			Code: common.UserHasDeleted,
		})
		return
	}

	var userId = strconv.FormatInt(u.(models.User).UserID, 10)
	var s dto.CollectRequest
	if err := c.ShouldBind(&s); err == nil {
		var dir models.Dir
		result := middleware.Db.Where("dir_id = ?", s.DirId).First(&dir)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, dto.CollectResponse{
				Code: common.DirNotExisted,
			})
			return
		}
		if strconv.FormatInt(dir.UserId, 10) != userId {
			c.JSON(http.StatusOK, dto.CollectResponse{
				Code: common.PermDenied,
			})
			return
		}
		var collection models.Collection
		result2 := middleware.Db.Where("dir_id = ? and word = ? and url_name = ?", s.DirId, s.Word, s.UrlName).First(&collection)
		if !errors.Is(result2.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, dto.CollectResponse{
				Code: common.CollectionHasExisted,
			})
			return
		}

		collection = models.Collection{DirId: s.DirId, Word: s.Word, UrlName: s.UrlName}
		result = middleware.Db.Create(&collection)
		middleware.Db.Model(&collection).Update("collection_id", int64(collection.ID))
		if result.Error != nil {
			c.JSON(http.StatusOK, dto.CollectResponse{
				Code: common.UnknownError,
			})
			fmt.Println("CollectErr: ", result.Error)
			return
		}

		c.JSON(http.StatusOK, dto.CollectResponse{
			Code: common.OK,
			Data: struct{ CollectionID string }{CollectionID: strconv.FormatInt(collection.CollectionId, 10)},
		})
		return
	} else {
		//参数校验失败的报错
		fmt.Println("bind err: ", err)
		c.JSON(http.StatusOK, dto.CollectResponse{Code: common.ParamInvalid})
	}

}
