package handler

import (
	"context"
	"fmt"
	"gin-micro-demo/app/A/cmd/rpc/proto"
	"log"
)

type AServerImpl struct {
	proto.UnimplementedAServer
}

func (s AServerImpl) HelloA(ctx context.Context, req *proto.AReq) (*proto.AResp, error) {
	log.Printf("recv %s 's msg \r\n", req.Name)
	resp := new(proto.AResp)
	resp.Res = fmt.Sprintf("A: hello %s", req.Name)
	log.Println("HelloA resp: ", resp)
	return resp, nil
}
