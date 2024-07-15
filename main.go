package main

import (
	"bestpath/config"
	"bestpath/context"
	"bestpath/utils/logger"
	"fmt"
	"os"
	"runtime"
	"time"
)

var globalConfig config.Config
var globalContext context.GlobalContext

func main() {
	runtime.GOMAXPROCS(1)
	// 参数判断
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s config_file\n", os.Args[0])
		os.Exit(1)
	}

	// 加载配置文件
	globalConfig = *config.LoadConfig(os.Args[1])

	// 设置日志级别, 并初始化日志
	logger.InitLogger(globalConfig.LogPath, globalConfig.LogLevel)

	// 解析config，加载杠杆和合约交易对，初始化context，账户初始化设置，拉取仓位、余额等
	globalContext.Init(&globalConfig)

	logger.Info("Staring Now")
	StartZmq(&globalConfig, &globalContext)

	// 开始定时ping aws服务器获取最优出口ip
	LoopPingAws(&globalConfig, &globalContext)

	// 阻塞主进程
	for {
		time.Sleep(24 * time.Hour)
	}
}
