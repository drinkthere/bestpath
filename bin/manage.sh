#!/bin/bash

# 根据输入参数执行相应的操作
case "$1" in
    "start")
        echo "Starting Service..."
        nohup ./bestpath ../config/config.json > ./logs/nohup-bp.log 2>&1 &
        ;;
    "stop")
        echo "Stopping Service..."
        pkill -f "./bestpath ../config/config.json"
        ;;
    *)
        echo "Usage: $0 {start|stop}"
        exit 1
        ;;
esac