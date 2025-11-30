package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	TotalRow int64 `json:"totalRow"`
	Offset   int   `json:"offset"`
	Limit    int   `json:"limit"`
}

type APIResponse struct {
	Status       string      `json:"status"`
	ResponseCode int         `json:"responseCode"`
	ResponseDesc string      `json:"responseDesc"`
	Timestamp    time.Time   `json:"timestamp"`
	Data         interface{} `json:"data,omitempty"`
	DataList     interface{} `json:"dataList,omitempty"`
	Pagination   *Pagination `json:"pagination,omitempty"`
}

func ResponseSuccess(c *gin.Context, data interface{}, desc string) {
	c.JSON(http.StatusOK, APIResponse{
		Status:       "success",
		ResponseCode: 200,
		ResponseDesc: desc,
		Timestamp:    time.Now(),
		Data:         data,
	})
}

func ResponseSuccessList(c *gin.Context, data interface{}, total int64, offset, limit int, desc string) {
	c.JSON(http.StatusOK, APIResponse{
		Status:       "success",
		ResponseCode: 200,
		ResponseDesc: desc,
		Timestamp:    time.Now(),
		DataList:     data,
		Pagination: &Pagination{
			TotalRow: total,
			Offset:   offset,
			Limit:    limit,
		},
	})
}

func ErrorResponse(c *gin.Context, code int, desc string) {
	c.JSON(code, APIResponse{
		Status:       "error",
		ResponseCode: code,
		ResponseDesc: desc,
		Timestamp:    time.Now(),
	})
}
