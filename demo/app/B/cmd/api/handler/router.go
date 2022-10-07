package handler

import (
	"gin-micro-demo/app/B/cmd/rpc/proto"
	"gin-micro-demo/config"
	"gin-micro-demo/utils"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Init_api(g *gin.Engine) {
	utils.RegisterHealth(g)
	g.GET("/hello", func(ctx *gin.Context) {
		func(g *gin.Context) {
			child, reporter := utils.SetContextTracer(g, "hello_call_ctx")
			defer reporter.Close()
			defer child.Finish()
			log.Println("ctx 传递跟踪")
			time.Sleep(time.Second)
		}(ctx)
		ctx.JSON(200, gin.H{
			"mag": "hello",
		})
	})

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
		SrvName: "srv_b_server",
	}

	g.GET("/b", func(ctx *gin.Context) {

		conn, err := utils.GinGetRpcCli(rc)
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		BRpcCli := proto.NewBClient(conn)
		resp, err := BRpcCli.HelloB(ctx, &proto.BReq{
			Name: "api",
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
}
