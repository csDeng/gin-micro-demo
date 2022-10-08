package utils

import (
	"fmt"
	"gin-micro-demo/config"
	"gin-micro-demo/ctx_const"
	"gin-micro-demo/otgrpc"
	"log"

	"github.com/gin-gonic/gin"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"google.golang.org/grpc"
)

func getTrace(c *config.ZipkinConfig) (opentracing.Tracer, reporter.Reporter) {
	// create our local service endpoint
	//记录服务名称和端口
	endpoint, err := zipkin.NewEndpoint(c.SERVICE_NAME, c.ZIPKIN_HTTP_ENDPOINT)
	log.Printf("zipkin endpointer: [ server_name ] = %s, [ hostPort ]  = %s", c.SERVICE_NAME, c.ZIPKIN_HTTP_ENDPOINT)

	// set up a span reporter
	//链路日志输出到哪
	reporter := zipkinhttp.NewReporter(c.ZIPKIN_RECORDER_HOST_PORT)
	log.Printf("ziplin reporter: %s", c.ZIPKIN_RECORDER_HOST_PORT)

	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// 采样器
	// initialize our tracer
	nativeTracer, err := zipkin.NewTracer(
		reporter,
		zipkin.WithLocalEndpoint(endpoint),
	)
	tracer := zipkinot.Wrap(nativeTracer)
	// use zipkin-go-opentracing to wrap our tracer
	return tracer, reporter

}

// 获取开启链路追踪的 rpc_server_opt
func GetGrpcOpt(c *config.ZipkinConfig) (grpc.ServerOption, reporter.Reporter) {
	tracer, reporter := getTrace(c)
	// optionally set as Global OpenTracing tracer instance
	opentracing.SetGlobalTracer(tracer)

	opts := grpc.UnaryInterceptor(
		// otgrpc.LogPayloads 是否记录 入参和出参
		// otgrpc.SpanDecorator 装饰器，回调函数
		// otgrpc.IncludingSpans 是否记录
		// grpc 拦截器
		grpc_middleware.ChainUnaryServer(
			// 链路追踪
			otgrpc.OpenTracingServerInterceptor(
				tracer,
				otgrpc.LogPayloads(),
				// IncludingSpans是请求前回调
				otgrpc.IncludingSpans(func(parentSpanCtx opentracing.SpanContext, method string, req, resp interface{}) bool {
					if method == "/grpc.health.v1.Health/Check" {
						// 健康检查不打印
						return false
					}
					log.Printf("method: %s", method)
					log.Printf("req: %+v", req)
					log.Printf("resp: %+v", resp)
					log.Println("\r\n----------------")
					return true
				}),
				// SpanDecorator是请求后回调
				otgrpc.SpanDecorator(func(span opentracing.Span, method string, req, resp interface{}, grpcError error) {
					if method == "/grpc.health.v1.Health/Check" {
						// 健康检查不打印
						return
					}
					log.Printf("method: %s", method)
					log.Printf("req: %+v", req)
					log.Printf("resp: %+v", resp)
					log.Printf("grpcError: %+v", grpcError)
					log.Println("\r\n----------------")
				}),
			),

			// 限流，注意限流应该在追踪链后面
			grpc.UnaryServerInterceptor(GetUnaryThrottleInterceptor(c.SERVICE_NAME)),
		),
	)
	return opts, reporter
}

// gin 全局中间件 链路追踪
func SetGinTracer(c *config.ZipkinConfig, g *gin.Engine) reporter.Reporter {
	tracer, reporter := getTrace(c)
	opentracing.SetGlobalTracer(tracer)
	// 将tracer注入到gin的中间件中
	g.Use(func(c *gin.Context) {
		if c.FullPath() == "/health" {
			// 健康检查不记录
			return
		}
		tracer = opentracing.GlobalTracer()
		c.Set(ctx_const.TracerCtxName, tracer)
		parentSpan := tracer.StartSpan("gin_middleware_" + c.FullPath() + c.RemoteIP())
		defer parentSpan.Finish()
		c.Set(ctx_const.ParentSpanCtxName, parentSpan)
		c.Set(ctx_const.GinContextName, c)
		c.Next()
	})
	return reporter

}

// gin context 链路追踪
func SetContextTracer(g *gin.Context, srv_name string) (opentracing.Span, reporter.Reporter) {
	ip_port := GetIpPort()
	ip, port := ip_port.Ip, ip_port.Port
	c, err := config.GetZipkinConfig()
	if err != nil {
		panic(err)
	}

	c.SERVICE_NAME = srv_name

	c.ZIPKIN_HTTP_ENDPOINT = fmt.Sprintf("%s:%d", ip, port)
	tracer, reporter := getTrace(c)
	var child opentracing.Span

	if span, ok := g.Get(ctx_const.ParentSpanCtxName); ok {
		parentSpan := span.(opentracing.Span)
		child = tracer.StartSpan("gin_context_withParentSpan_tracer_"+g.FullPath(), opentracing.ChildOf(parentSpan.Context()))
	} else {
		child = tracer.StartSpan("gin_context_tracer_" + g.FullPath() + "_" + g.RemoteIP())
	}

	return child, reporter

}
