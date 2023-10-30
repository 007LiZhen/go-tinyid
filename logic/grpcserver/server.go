package grpcserver

import (
	"context"

	"gitee.com/git-lz/go-tinyid/common/xgrpc/proto"
	"gitee.com/git-lz/go-tinyid/logic/idsequence"
)

func (is *IdSequenceServer) Get(ctx context.Context, req *proto.SendRequest) (
	resp *proto.SendResponse,
	err error) {
	idSequence, err := idsequence.GetIdSequence(req.Biz)
	if err != nil {
		return nil, err
	}

	resp = new(proto.SendResponse)
	resp.Id, err = idSequence.GetOne()
	return resp, err
}
