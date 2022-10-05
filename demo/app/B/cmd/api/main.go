package main

import (
	"fmt"
	"gin-micro-demo/app/B/cmd/rpc/proto"
	"gin-micro-demo/config"
	"gin-micro-demo/utils"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	port, err := utils.GetFreePort()
	if err != nil {
		panic(err)
	}
	ip, err := utils.GetIp()
	if err != nil {
		panic(err)
	}
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
		SrvName: "srv_b",
	}

	conn, err := utils.GetRpcCli(rc)
	if err != nil {
		panic(err)
	}
	BRpcCli := proto.NewBClient(conn)

	g := gin.Default()
	g.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"mag": "hello",
		})
	})
	g.GET("/b", func(ctx *gin.Context) {
		resp, err := BRpcCli.HelloB(ctx, &proto.BReq{
			Name: "from api",
		})
		if err != nil {
			log.Println(err)
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"msg": resp.Res,
		})

	})
	utils.RegisterHealth(g)

	name := "b_api"

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
	if conn != nil {
		conn.Close()
	}
	if err = utils.DeRegister(apiConfig); err != nil {
		log.Printf("[%s] 服务注销失败; error = %v \r\n", name, err)
	} else {
		log.Printf("[%s] 服务注销成功; \r\n", name)
	}
}
