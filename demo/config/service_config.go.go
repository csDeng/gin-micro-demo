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
