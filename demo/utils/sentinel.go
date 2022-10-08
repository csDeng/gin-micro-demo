package utils

import (
	"context"
	"errors"
	"log"
	"net/http"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// 配置一个 5qps 的限流器
func InitSentinel(srv_name string) {
	// We should initialize Sentinel first.
	conf := config.NewDefaultConfig()
	// for testing, logging output to console
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		log.Fatalf("初始化异常: %+v", err)
	}

	// 配置限流规则
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource: srv_name,
			// TokenCalculateStrategy当前流量控制器的Token计算策略。Direct表示直接使用字段 Threshold 作为阈值；WarmUp表示使用预热方式计算Token的阈值。
			TokenCalculateStrategy: flow.Direct,

			ControlBehavior: flow.Reject,
			Threshold:       5,
			// StatIntervalInMs: 规则对应的流量控制器的独立统计结构的统计周期。如果StatIntervalInMs是1000，也就是统计QPS。
			StatIntervalInMs: 1000,
		},
	})
	if err != nil {
		log.Fatalf("加载规则失败 error: %+v", err)
		return
	}

}

// gin 使用限流中间件
func SetSentinelMiddleware(g *gin.Engine, srv_name string) {
	g.Use(func(c *gin.Context) {
		e, b := sentinel.Entry(srv_name, sentinel.WithTrafficType(base.Inbound))
		if b != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"msg": "请求过于频繁,请稍后重试!!!",
			})
			return
		} else {
			// 处理请求
			c.Next()
			e.Exit()
		}
	})
}

type RpcInterceptor func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error)

// rpc 服务端的拦截器
// 一元RPC拦截器
func GetUnaryThrottleInterceptor(srv_name string) RpcInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		e, b := sentinel.Entry(srv_name, sentinel.WithTrafficType(base.Inbound))
		if b != nil {
			err := errors.New("请求过于频繁,请稍后重试!!!")
			return nil, err
		}
		// 执行具体rpc方法
		resp, err = handler(ctx, req)
		if err != nil {
			return resp, err
		}
		e.Exit()
		return resp, nil
	}
}
