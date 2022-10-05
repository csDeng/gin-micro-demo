package main

import (
	"fmt"
	"gin-micro-demo/app/B/cmd/rpc/handler"
	"gin-micro-demo/app/B/cmd/rpc/proto"
	"gin-micro-demo/config"
	"gin-micro-demo/utils"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 获取本机ip
	ip, err := utils.GetIp()
	if err != nil {
		panic(err)
	}

	// 获取随机端口
	port, err := utils.GetFreePort()
	if err != nil {
		panic(err)
	}

	name := "srv_b"

	server := grpc.NewServer()
	proto.RegisterBServer(server, handler.BServerImpl{})
	log.Printf("%s is running at %s:%d", name, ip, port)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))

	// 注册服务健康检查
	healthServer := health.NewServer()
	system := name
	healthServer.SetServingStatus(system, grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(server, healthServer)

	// 生成注册对象
	cfg := &config.ServiceConfig{
		Host: ip,
		Port: port,
		Name: name,
		Id:   utils.GetUUId(),
		Tags: []string{"B", "rpc"},
	}
	// 注册rpc
	err = utils.RegisterRpc(cfg)
	if err != nil {
		panic(err)
	}

	// 开启 grpc 反射，方便调试
	reflection.Register(server)
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()
	// 优雅重启，退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = utils.DeRegister(cfg); err != nil {
		log.Printf("[%s] 服务注销失败; error = %v \r\n", name, err)
	} else {
		log.Printf("[%s] 服务注销成功; \r\n", name)
	}

}
