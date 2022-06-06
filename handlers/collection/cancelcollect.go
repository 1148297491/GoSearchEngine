package collection

import (
	"GoSearchEngine/common"
	"GoSearchEngine/dto"
	"GoSearchEngine/middleware"
	"GoSearchEngine/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func CancelCollect(c *gin.Context) {
	var s dto.CancelCollectRequest
	if err := c.ShouldBind(&s); err == nil {
		var collection models.Collection
		result := middleware.Db.Where("collection_id = ?", s.CollectionId).First(&collection)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, dto.CancelCollectResponse{
				Code: common.CollectionNotExisted,
			})
			return
		}

		middleware.Db.Delete(&collection)
		c.JSON(http.StatusOK, dto.CancelCollectResponse{
			Code: common.OK,
		})
		return
	} else {
		//参数校验失败的报错
		fmt.Println("bind err: ", err)
		c.JSON(http.StatusOK, dto.CancelCollectResponse{Code: common.ParamInvalid})
	}

}
