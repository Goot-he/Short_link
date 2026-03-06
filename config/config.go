package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

var GlobalCfg *GlobalConfig

type GlobalConfig struct {
	Server    ServerConfig    `yaml:"Server"`
	DataBase  DataBaseConfig  `yaml:"Database"`
	Redis     RedisConfig     `yaml:"Redis"`
	Bloom     BloomConfig     `yaml:"Bloom"`
	SnowFlake SnowFlakeConfig `yaml:"Snowflake"`
	Logger    LoggerConfig    `yaml:"Logger"`
}

type ServerConfig struct {
	HttpPort string `yaml:"HttpPort"`
	Mode     string `yaml:"Mode"`
}

type DataBaseConfig struct {
	DbType             string `yaml:"DbType"`
	DbHost             string `yaml:"DbHost"`
	DbPort             string `yaml:"DbPort"`
	DbUser             string `yaml:"DbUser"`
	DbPassword         string `yaml:"DbPassword"`
	DbName             string `yaml:"DbName"`
	MaxIdleConnections int    `yaml:"MaxIdleConnections"`
	MaxOpenConnections int    `yaml:"MaxOpenConnections"`
}

type RedisConfig struct {
	RdbPort  string `yaml:"RdbPort"`
	Password string `yaml:"Password"`
	Rdb      int    `yaml:"Rdb"`
}

type BloomConfig struct {
	P float64 `yaml:"P"`
	N int64   `yaml:"N"`
}
type SnowFlakeConfig struct {
	MachineID int64 `yaml:"MachineID"`
	Epoch     int64 `yaml:"Epoch"`
}

type LoggerConfig struct {
	Level   string `yaml:"Level"`
	LogType string `yaml:"LogType"`
}

func InitGlobalConfig() error {
	filePath := "./config/config.yaml"

	// 1.打开配置文件
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 2.解析YAML到结构体
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&GlobalCfg)
	return err
}
