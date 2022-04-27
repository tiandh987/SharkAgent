package core

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
	"github.com/tiandh987/SharkAgent/pkg/log"
	"net/http"
)

// ErrResponse 定义发生错误时的返回消息。
// 如果不存在引用将被省略。
// swagger:model
type ErrResponse struct {
	// Code defines the business error code.
	Code int `json:"code"`

	// Message contains the detail of this message.
	// This message is suitable to be exposed to external
	Message string `json:"message"`

	// Reference returns the reference document which maybe useful to solve this error.
	Reference string `json:"reference,omitempty"`
}

// WriteResponse 将 错误 或 响应数据 写入 http 响应正文。
// 它使用 errors.ParseCoder 将任何错误解析为 errors.Coder
// errors.Coder 包含 错误代码、用户安全错误消息 和 http 状态代码。
func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		log.Errorf("%#+v", err)
		coder := errors.ParseCoder(err)
		c.JSON(coder.HTTPStatus(), ErrResponse{
			Code:      coder.Code(),
			Message:   coder.String(),
			Reference: coder.Reference(),
		})

		return
	}

	c.JSON(http.StatusOK, data)
}