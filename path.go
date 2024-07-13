package main

import (
	"bestpath/config"
	"bestpath/container"
	"bestpath/context"
	"bestpath/utils/logger"
	"github.com/drinkthere/ping"
	"sync"
	"time"
)

func LoopPingAws(cfg *config.Config, globalContext *context.GlobalContext) {
	go func() {
		for {
			time.Sleep(1 * time.Hour)
			var bestPath container.BestPath
			bestLatency := time.Duration(1<<63 - 1)
			// Use a WaitGroup to track the completion of all ping operations
			var wg sync.WaitGroup

			for _, sourceIP := range cfg.SourceIPs {
				for _, targetIP := range cfg.TargetIPs {
					wg.Add(1)
					go func(localEth, targetIP string) {
						defer wg.Done()

						pinger, err := ping.NewPinger(targetIP)
						if err != nil {
							panic(err)
						}
						pinger.SetSource(localEth)
						pinger.Count = 3
						err = pinger.Run() // Blocks until finished.
						if err != nil {
							panic(err)
						}
						stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
						logger.Info("Average latency from %s to %s: %.5fms", localEth, targetIP, float64(stats.AvgRtt.Microseconds())/1000)

						if stats.AvgRtt < bestLatency {
							bestLatency = stats.AvgRtt
							bestPath = container.BestPath{
								SourceIP: localEth,
								TargetIP: targetIP,
							}
						}
					}(sourceIP, targetIP)
				}
			}

			// Wait for all ping operations to complete
			wg.Wait()

			logger.Info("bestP path is %+v", bestPath)
			currentBestPath := globalContext.BestPath
			if currentBestPath.SourceIP != bestPath.SourceIP || currentBestPath.TargetIP != bestPath.TargetIP {
				logger.Warn("Best path changed from %s->%s to %s->%s", currentBestPath.SourceIP,
					currentBestPath.TargetIP, bestPath.SourceIP, bestPath.TargetIP)
				globalContext.BestPath.UpdateBestPath(bestPath)
				globalContext.BestPathChangedCh <- struct{}{}
			} else {
				logger.Info("Best path does not change")
			}
		}
	}()
}
