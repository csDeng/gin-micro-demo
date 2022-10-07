package config

// 服务配置
type ServiceConfig struct {
	Host string   `mapstructure:"host" json:"host"`
	Port int      `mapstructure:"port" json:"port"`
	Name string   `mapstructure:"name" json:"name"`
	Id   string   `mapstructure:"id" json:"id"`
	Tags []string `mapstructure:"tags" json:"tags"`
}

type ConsulConfig struct {
	Host string `mapstructure:"consul_host" json:"consul_host"`
	Port int    `mapstructure:"consul_port" json:"consul_port"`
}

type RpcCliConfig struct {
	ConsulConfig
	SrvName string
}

type FilterServiceConfig struct {
	ConsulConfig
	SrvName string
}

type NacosConsul struct {
	ConsulDataId string
	ConsulGroup  string
}

type NacosConfig struct {
	NacosConsul
}

type GlobalConfig struct {
}

type ZipkinConfig struct {
	SERVICE_NAME              string `mapstructure:"SERVICE_NAME" json:"SERVICE_NAME"`
	ZIPKIN_RECORDER_HOST_PORT string `mapstructure:"ZIPKIN_RECORDER_HOST_PORT" json:"ZIPKIN_RECORDER_HOST_PORT"`
	ZIPKIN_HTTP_ENDPOINT      string `mapstructure:"ZIPKIN_HTTP_ENDPOINT" json:"ZIPKIN_HTTP_ENDPOINT"`
}

type IpPort struct {
	Ip   string
	Port int
}
