package handler

import (
	"context"
	"fmt"
	"gin-micro-demo/app/B/cmd/rpc/proto"
	"log"
)

type BServerImpl struct {
	proto.UnimplementedBServer
}

func (s BServerImpl) HelloB(ctx context.Context, req *proto.BReq) (*proto.BResp, error) {
	log.Printf("recv %s 's msg \r\n", req.Name)
	resp := new(proto.BResp)
	resp.Res = fmt.Sprintf("hello %s", req.Name)
	return resp, nil
}
