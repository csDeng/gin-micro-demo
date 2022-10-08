package main

import (
	"fmt"
	"gin-micro-demo/app/B/cmd/api/handler"
	"gin-micro-demo/config"
	"gin-micro-demo/utils"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	ip_port := utils.GetIpPort()
	ip, port := ip_port.Ip, ip_port.Port

	g := gin.Default()
	name := "b_api"
	// 注入限流中间件
	utils.InitSentinel(name)
	utils.SetSentinelMiddleware(g, name)

	zipConfig, err := config.GetZipkinConfig()
	if err != nil {
		panic(err)
	}

	zipConfig.SERVICE_NAME = name
	zipConfig.ZIPKIN_HTTP_ENDPOINT = fmt.Sprintf("%s:%d", ip, port)

	// 设置全局链路追踪
	reporter := utils.SetGinTracer(zipConfig, g)
	handler.Init_api(g)

	defer func() {
		if reporter != nil {
			reporter.Close()
		}
	}()

	// 启动api
	go func() {
		err := g.Run(fmt.Sprintf(":%d", port))
		if err != nil {
			panic(err)
		}
		log.Printf("%s is running at %s:%d", name, ip, port)
	}()

	apiConfig := &config.ServiceConfig{
		Host: ip,
		Port: port,
		Name: name,
		Id:   utils.GetUUId(),
		Tags: []string{name},
	}

	// 服务注册
	if err = utils.RegisterApi(apiConfig); err != nil {
		log.Println("api 注册失败")
		panic(err)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err = utils.DeRegister(apiConfig); err != nil {
		log.Printf("[%s] 服务注销失败; error = %v \r\n", name, err)
	} else {
		log.Printf("[%s] 服务注销成功; \r\n", name)
	}
}
