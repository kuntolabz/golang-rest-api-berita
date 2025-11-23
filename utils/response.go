package utils

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	ResponseCode    int         `json:"responseCode"`
	ResponseMessage string      `json:"ResponseMessage"`
	Data            interface{} `json:"data,omitempty"`
	Error           interface{} `json:"error,omitempty"`
}

func SuccessResponse(c *gin.Context, data interface{}, total int64, offset int, limit int) {
	value := reflect.ValueOf(data)
	isSlice := value.Kind() == reflect.Slice

	response := gin.H{
		"responseCode": 200,
	}

	if isSlice {
		// Slice dengan lebih dari 1 row → dataList + pagination
		if value.Len() > 1 {
			response["dataList"] = data
			response["pagination"] = gin.H{
				"totalRow": total,
				"offset":   offset,
				"limit":    limit,
			}
		} else if value.Len() == 1 {
			// Slice dengan 1 row → data tanpa pagination
			response["data"] = value.Index(0).Interface()
		} else {
			// Slice kosong → dataList kosong, tanpa pagination
			response["dataList"] = []interface{}{}
		}
	} else {
		// Bukan slice → single object → tanpa pagination
		response["data"] = data
	}

	c.JSON(http.StatusOK, response)
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, ResponseData{
		ResponseCode:    200,
		ResponseMessage: "Success",
		Data:            data,
	})
}

func ErrorResponse(c *gin.Context, code int, message string, errors ...string) {
	resp := ResponseData{
		ResponseCode:    code,
		ResponseMessage: message,
	}

	// Jika ada errors, isi
	if len(errors) > 0 && errors[0] != "" {
		resp.Error = errors[0]
	}

	c.JSON(code, resp)
}
