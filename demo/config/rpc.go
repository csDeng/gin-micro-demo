package config

import (
	AP "gin-micro-demo/app/A/cmd/rpc/proto"
	BP "gin-micro-demo/app/B/cmd/rpc/proto"
)

type RpcCli struct {
	A_srv AP.AClient
	B_srv BP.BClient
}
