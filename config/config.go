package config

import (
	"encoding/json"
	"go.uber.org/zap/zapcore"
	"os"
)

type Config struct {
	// 日志配置
	LogLevel zapcore.Level
	LogPath  string

	InitSourceIP string // 初始化时的本地最优IP
	InitTargetIP string // 初始化时的aws最优IP

	SourceIPs []string // 本地IP组
	TargetIPs []string // aws内网IP组

	BestPathChangedZMQIPC string // 当最优路径变化时通过这个IPC PUB 消息
}

func LoadConfig(filename string) *Config {
	config := new(Config)
	reader, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	// 加载配置
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	return config
}
