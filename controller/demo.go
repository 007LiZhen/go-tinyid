package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/007LiZhen/go-tinyid/common/dto/demo"
	"github.com/007LiZhen/go-tinyid/common/dto/response"
	"github.com/007LiZhen/go-tinyid/common/merrors"
	"github.com/007LiZhen/go-tinyid/logic/idsequence"
)

type ID struct{}

func NewID() *ID {
	return &ID{}
}

func (id *ID) Get(ctx *gin.Context) {
	resp := &response.Response{}
	defer func() {
		response.EchoResponse(ctx, resp)
	}()

	var req = new(demo.GetIdReq)
	if err := ctx.BindQuery(&req); err != nil {
		resp.WithMsg(merrors.ErnoRequestBindFailed, err.Error())
		return
	}

	idSequence, err := idsequence.GetIdSequence(req.Biz)
	if err != nil {
		resp.WithMsg(merrors.ErnoDataNotSupport, err.Error())
		return
	}

	nextId, err := idSequence.GetOne()
	if err != nil {
		resp.WithMsg(merrors.ErnoGetNextIdFailed, err.Error())
		return
	}

	resp.WithData(nextId)
	return
}
