package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitee.com/git-lz/go-tinyid/common/merrors"
)

type Response struct {
	Errno  int         `json:"code"`
	Msg    string      `json:"message"`
	Data   interface{} `json:"data"`
	IsSucc bool        `json:"is_succ"`
}

func NewResponse() *Response {
	return &Response{
		IsSucc: true,
	}
}

func (r *Response) WithDefaultMsg(errno int) *Response {
	tx := r
	msg, ok := merrors.ErnoMsgMap[errno]
	if !ok {
		msg = "unknown error"
	}

	if errno != merrors.ErnoSuccess {
		tx.IsSucc = false
	}

	tx.Msg = msg
	tx.Errno = errno
	return tx
}

func (r *Response) WithMsg(errno int, msg string) *Response {
	tx := r
	if errno != merrors.ErnoSuccess {
		tx.IsSucc = false
	}

	tx.Msg = msg
	tx.Errno = errno
	return tx
}

func (r *Response) WithData(data interface{}) *Response {
	tx := r

	if tx.Errno != merrors.ErnoSuccess {
		tx.IsSucc = false
	}

	tx.Data = data
	return tx
}

func EchoResponse(ctx *gin.Context, resp *Response) {
	if resp.Errno != merrors.ErnoSuccess {
		resp.IsSucc = false
	} else {
		resp.IsSucc = true
	}

	ctx.JSON(http.StatusOK, resp)
}
