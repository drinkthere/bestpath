package main

import (
	"bestpath/config"
	"bestpath/context"
	"bestpath/protocol/pb"
	"bestpath/utils/logger"
	zmq "github.com/pebbe/zmq4"
	"google.golang.org/protobuf/proto"
	"os"
)

func StartZmq(cfg *config.Config, globalContext *context.GlobalContext) {
	go func() {
		ctx, err := zmq.NewContext()
		if err != nil {
			logger.Fatal("[ZMQ] Failed to create context, error: %s", err.Error())
			os.Exit(1)
		}
		pub, err := ctx.NewSocket(zmq.PUB)
		if err != nil {
			ctx.Term()
			logger.Fatal("[ZMQ] Failed to create PUB socket, error: %s", err.Error())
			os.Exit(2)
		}
		err = pub.Bind(cfg.BestPathChangedZMQIPC)
		if err != nil {
			ctx.Term()
			logger.Fatal("[ZMQ] Failed to bind IPC %s, error: %s", cfg.BestPathChangedZMQIPC, err.Error())
			os.Exit(3)
		}

		defer pub.Close()
		defer ctx.Term()

		for {
			select {
			case <-globalContext.BestPathChangedCh:
				md := &pb.BestPath{
					SourceIP: globalContext.BestPath.SourceIP,
					TargetIP: globalContext.BestPath.TargetIP,
					AvgRtt:   int64(globalContext.BestPath.AvgRtt),
				}

				data, err := proto.Marshal(md)
				if err != nil {
					logger.Error("[ZMQ] Error marshaling MarketData: %v", err)
					continue
				}
				_, err = pub.Send(string(data), 0)
				if err != nil {
					logger.Error("[ZMQ] Error sending MarketData: %v", err)
					continue
				}
			}
		}
	}()
}
