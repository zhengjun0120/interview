package conf

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server    Server    `yaml:"server"`
	Mysql     Mysql     `yaml:"mysql"`
	Redis     Redis     `yaml:"redis"`
	Snowflake Snowflake `yaml:"snowflake"`
	Jwt       Jwt       `yaml:"jwt"`
	Smtp      Smtp      `yaml:"smtp"`
	Cos       Cos       `yaml:"cos"`
	AiChat    AiChat    `yaml:"ai_chat"`
}

type Server struct {
	Port               string `yaml:"port"`
	Host               string `yaml:"host"`
	YourFrontendDomain string `yaml:"your_frontend_domain"`
}

type Mysql struct {
	Dsn string `yaml:"dsn"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DB       int    `yaml:"db"`
	Password string `yaml:"password"`
}

type Snowflake struct {
	NodeId int64 `yaml:"node_id"`
}

type Jwt struct {
	SecretKey   string `yaml:"secret_key"`
	ExpireHours int    `yaml:"expire_hours"`
}

type Smtp struct {
	SmtpHost    string `yaml:"smtp_host"`
	SmtpPort    string `yaml:"smtp_port"`
	SmtpUser    string `yaml:"smtp_user"`
	SmtpPass    string `yaml:"smtp_pass"`
	EncodedName string `yaml:"encoded_name"`
}

type Cos struct {
	SecretId  string `yaml:"secret_id"`
	SecretKey string `yaml:"secret_key"`
	BucketUrl string `yaml:"bucket_url"`
}

type AiChat struct {
	ApiKey  string `yaml:"api_key"`
	ModelId string `yaml:"model_id"`
	BaseUrl string `yaml:"base_url"`
}

var config Config

func GetConfig() Config {
	return config
}

func InitConfig() {
	file, err := os.Open("../conf/config.yaml")
	if err != nil {
		panic("配置文件加载失败" + err.Error())
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		panic("配置文件解析失败" + err.Error())
	}

	fmt.Println("配置文件加载成功\n", config)
}
