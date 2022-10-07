package handler

import (
	"context"
	"fmt"
	"gin-micro-demo/app/A/cmd/rpc/proto"
)

type AServerImpl struct {
	proto.UnimplementedAServer
}

func (s AServerImpl) HelloA(ctx context.Context, req *proto.AReq) (*proto.AResp, error) {
	resp := new(proto.AResp)
	resp.Res = fmt.Sprintf("A_RPC_send: hello %s", req.Name)
	return resp, nil
}
