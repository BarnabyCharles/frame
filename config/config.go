package config

type AppConfig struct {
	Mysql MysqlConfig `yaml:"Mysql" mapstructure:"Mysql"`
	Nacos NacosConfig `yaml:"Nacos" json:"Nacos" mapstructure:"Nacos"`
	Redis RedisConfig `json:"Redis" mapstruture:"Redis"`
}

type MysqlConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
}

type RedisConfig struct {
	Host string `json:"host" yaml:"host" mapstruture:"host"`
	Port int    `json:"port" yaml:"port" mapstruture:"port"`
}

type NacosConfig struct {
	Host        string `json:"Host" yaml:"Host" mapstructure:"Host"`
	Port        int    `json:"Port" yaml:"Port" mapstructure:"Port"`
	ServerName  string `json:"ServerName" yaml:"ServerName" mapstructure:"ServerName"`
	Group       string `json:"Group" yaml:"Group" mapstructure:"Group"`
	NamespaceId string `json:"NamespaceId" mapstructure:"NamespaceId"`
}
