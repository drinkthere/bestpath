package container

type BestPath struct {
	SourceIP string
	TargetIP string
}

func (bp *BestPath) Init(sourceEth string, targetIP string) {
	bp.SourceIP = sourceEth
	bp.TargetIP = targetIP
}

func (bp *BestPath) UpdateBestPath(bestPath BestPath) {
	bp.SourceIP = bestPath.SourceIP
	bp.TargetIP = bestPath.TargetIP
}
