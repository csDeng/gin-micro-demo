package utils

import (
	"fmt"
	"gin-micro-demo/config"
	"log"
	"sync"

	"github.com/hashicorp/consul/api"
)

var (
	consul_config *config.ConsulConfig
	once          sync.Once
)

func getConsul() {
	once.Do(func() {
		var err error
		consul_config, err = config.GetConsulConfig()
		if err != nil {
			panic(err)
		}
	})
}

// 服务注册
func RegisterRpc(c *config.ServiceConfig) error {
	id, name, tags, host, port := c.Id, c.Name, c.Tags, c.Host, c.Port
	getConsul()
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", consul_config.Host, consul_config.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	// 生成注册对象
	registration := &api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Tags:    tags,
		Port:    port,
		Address: host,

		// grpc 健康检查
		Check: &api.AgentServiceCheck{
			Interval:                       "5s",
			Timeout:                        "5s",
			GRPC:                           fmt.Sprintf("%s:%d", host, port),
			GRPCUseTLS:                     false,
			DeregisterCriticalServiceAfter: "15s",
		},
	}
	err = client.Agent().ServiceRegister(registration)
	log.Printf("服务注册:  %+v \r\n", registration)
	log.Printf("健康检查: %+v \r\n", registration.Check.GRPC)
	return err
}

func RegisterApi(c *config.ServiceConfig) error {
	id, name, tags, host, port := c.Id, c.Name, c.Tags, c.Host, c.Port
	getConsul()
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", consul_config.Host, consul_config.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	// 生成注册对象
	registration := &api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Tags:    tags,
		Port:    port,
		Address: host,

		// http 健康检查
		Check: &api.AgentServiceCheck{
			Interval:                       "5s",
			Timeout:                        "5s",
			HTTP:                           fmt.Sprintf("http://%s:%d/health", host, port),
			DeregisterCriticalServiceAfter: "15s",
		},
	}
	err = client.Agent().ServiceRegister(registration)
	log.Printf("Api注册:  %+v \r\n", registration)
	log.Printf("健康检查: %+v \r\n", registration.Check.GRPC)
	return err
}

// 服务注销
func DeRegister(c *config.ServiceConfig) error {
	getConsul()
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", consul_config.Host, consul_config.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	return client.Agent().ServiceDeregister(c.Id)
}

// 服务发现

// 服务过滤
// 服务过滤
func FilterService(c *config.FilterServiceConfig) (map[string]*api.AgentService, error) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", c.ConsulConfig.Host, c.ConsulConfig.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	return client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", c.SrvName))
}
