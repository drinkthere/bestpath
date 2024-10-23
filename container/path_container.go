package container

import "time"

type BestPath struct {
	SourceIP string
	TargetIP string
	AvgRtt   time.Duration
}

func (bp *BestPath) Init(sourceEth string, targetIP string) {
	bp.SourceIP = sourceEth
	bp.TargetIP = targetIP
	bp.AvgRtt = time.Duration(1<<63 - 1)
}

func (bp *BestPath) UpdateBestPath(bestPath BestPath) {
	bp.SourceIP = bestPath.SourceIP
	bp.TargetIP = bestPath.TargetIP
	bp.AvgRtt = bestPath.AvgRtt
}
