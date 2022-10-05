package utils

import (
	"fmt"
	"gin-micro-demo/config"

	"github.com/gin-gonic/gin"

	// consul 协议解析器
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
)

func GetRpcCli(c *config.RpcCliConfig) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=15s", c.ConsulConfig.Host, c.ConsulConfig.Port, c.SrvName),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "round_robin"}`)),
	)
	return conn, err
}

// 实现健康检查接口
func RegisterHealth(g *gin.Engine) {
	g.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"msg": "ok",
		})
	})
}
