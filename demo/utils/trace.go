package utils

import (
	"gin-micro-demo/config"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
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
	// optionally set as Global OpenTracing tracer instance
	opentracing.SetGlobalTracer(tracer)
	// use zipkin-go-opentracing to wrap our tracer
	return tracer, reporter

}

// 获取开启链路追踪的 rpc_server_opt
func GetGrpcOpt(c *config.ZipkinConfig) (grpc.ServerOption, reporter.Reporter) {
	tracer, reporter := getTrace(c)
	opts := grpc.UnaryInterceptor(
		// otgrpc.LogPayloads 是否记录 入参和出参
		// otgrpc.SpanDecorator 装饰器，回调函数
		// otgrpc.IncludingSpans 是否记录
		// grpc 拦截器
		otgrpc.OpenTracingServerInterceptor(
			tracer,
			otgrpc.LogPayloads(),
			// IncludingSpans是请求前回调
			otgrpc.IncludingSpans(func(parentSpanCtx opentracing.SpanContext, method string, req, resp interface{}) bool {
				if method == "/grpc.health.v1.Health/Check" {
					// 健康检查不打印
					return true
				}
				log.Printf("\r\n parent= %+v \r\n", parentSpanCtx)
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
	)
	return opts, reporter
}

// g 链路追踪
func SetGinTracer(c *config.ZipkinConfig, g *gin.Engine) reporter.Reporter {
	trace, reporter := getTrace(c)
	// 将tracer注入到gin的中间件中
	g.Use(func(c *gin.Context) {
		span := trace.StartSpan(c.FullPath())
		defer span.Finish()
		c.Next()
	})
	return reporter

}
