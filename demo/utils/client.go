package utils

import (
	"fmt"
	"gin-micro-demo/config"
	"gin-micro-demo/otgrpc"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	// consul 协议解析器
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
)

func RpcGetRpcCli(c *config.RpcCliConfig) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=15s", c.ConsulConfig.Host, c.ConsulConfig.Port, c.SrvName),
		grpc.WithInsecure(),

		// rpc 负载均衡
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "round_robin"}`)),

		// 链路追踪 span 注入
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	return conn, err
}

func GinGetRpcCli(c *config.RpcCliConfig) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=15s", c.ConsulConfig.Host, c.ConsulConfig.Port, c.SrvName),
		grpc.WithInsecure(),

		// rpc 负载均衡
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "round_robin"}`)),

		// 链路追踪 span 注入
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
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
