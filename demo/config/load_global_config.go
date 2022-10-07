package config

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var (
	GlobalConfigObj *GlobalConfig
	once            sync.Once
	nacosConfig     *NacosConfig

	nacosCli  config_client.IConfigClient
	nacosOnce sync.Once
)

func getNacosClient() config_client.IConfigClient {
	once.Do(func() {
		serverConfigs := []constant.ServerConfig{
			{
				IpAddr:      "http://127.0.0.1",
				ContextPath: "/nacos",
				Port:        8848,
				Scheme:      "http",
			},
		}

		clientConfig := constant.ClientConfig{
			NamespaceId:         "", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogDir:              "tmp/nacos/log",
			CacheDir:            "tmp/nacos/cache",
			LogLevel:            "debug",
		}
		var err error
		nacosCli, err = clients.CreateConfigClient(map[string]interface{}{
			"serverConfigs": serverConfigs,
			"clientConfig":  clientConfig,
		})

		if err != nil {
			panic(err)
		}
	})
	return nacosCli
}

func getNacosContent(dataId, group string) (string, error) {
	nacosCli := getNacosClient()
	content, err := nacosCli.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group})
	if err != nil {
		return "", err
	}

	// 监听配置文件的变化
	err = nacosCli.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("文件发生变化")
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})
	if err != nil {
		return "", err
	}
	return content, nil
}

func GetConsulConfig() (*ConsulConfig, error) {
	consul_config := NacosConsul{
		ConsulDataId: "demo_consul",
		ConsulGroup:  "dev",
	}
	content, err := getNacosContent(
		consul_config.ConsulDataId, consul_config.ConsulGroup)
	if err != nil {
		return nil, err
	}
	config := new(ConsulConfig)
	err = json.Unmarshal([]byte(content), config)
	if err != nil {
		return nil, err
	}
	log.Printf("从配置中心，获取 consul 配置 %+v \r\n", config)
	return config, nil
}

// @TODO 可以使用单例模式 + 监听配置更改优化
func GetZipkinConfig() (*ZipkinConfig, error) {
	consul_config := NacosConsul{
		ConsulDataId: "demo_zipkin",
		ConsulGroup:  "dev",
	}
	content, err := getNacosContent(
		consul_config.ConsulDataId, consul_config.ConsulGroup)
	if err != nil {
		return nil, err
	}
	config := new(ZipkinConfig)
	err = json.Unmarshal([]byte(content), config)
	if err != nil {
		return nil, err
	}
	log.Printf("从配置中心，获取 zipkin 配置 %+v \r\n", config)
	return config, nil
}
