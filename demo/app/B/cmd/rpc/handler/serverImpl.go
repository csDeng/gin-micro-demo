package handler

import (
	"context"
	AP "gin-micro-demo/app/A/cmd/rpc/proto"
	"gin-micro-demo/app/B/cmd/rpc/proto"
	"gin-micro-demo/config"
	"gin-micro-demo/utils"
)

type BServerImpl struct {
	proto.UnimplementedBServer
}

func (s BServerImpl) HelloB(ctx context.Context, req *proto.BReq) (*proto.BResp, error) {
	resp := new(proto.BResp)
	c, err := config.GetConsulConfig()
	if err != nil {
		panic(err)
	}

	// 获取服务配置
	rc := &config.RpcCliConfig{
		ConsulConfig: config.ConsulConfig{
			Host: c.Host,
			Port: c.Port,
		},
		SrvName: "srv_a_server",
	}
	ArpcConfig := config.RpcCliConfig{
		ConsulConfig: config.ConsulConfig{
			Host: rc.ConsulConfig.Host,
			Port: rc.ConsulConfig.Port,
		},
		SrvName: rc.SrvName,
	}
	ARpcCli, err := utils.RpcGetRpcCli(&ArpcConfig)
	if err != nil {
		return nil, err
	}
	client := AP.NewAClient(ARpcCli)
	r, err := client.HelloA(
		ctx,
		&AP.AReq{
			Name: "B call A",
		})
	if err != nil {
		return nil, err
	}

	resp.Res = r.Res
	return resp, nil
}
