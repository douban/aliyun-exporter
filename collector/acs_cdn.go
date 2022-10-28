package collector

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
)

type CdnExporter struct {
	project Project
}

func NewCdnExporter(c *cms.Client) *CdnExporter {
	return &CdnExporter{
		project: Project{
			client:    c,
			Namespace: "acs_cdn",
		},
	}
}

type statusPoint struct {
	datapoint
	Status  string  `json:"status"`
}

func (db *CdnExporter) RetrieveHitRate(rangeTime int64, delayTime int64) []datapoint {
	return retrieve("hitRate", db.project, rangeTime, delayTime)
}

func (db *CdnExporter) RetrieveOriBps(rangeTime int64, delayTime int64) []datapoint {
	return retrieve("ori_bps", db.project, rangeTime, delayTime)
}

func (db *CdnExporter) RetrieveL1Acc(rangeTime int64, delayTime int64) []datapoint {
	return retrieve("l1_acc", db.project, rangeTime, delayTime)
}

func (db *CdnExporter) RetrieveOriAcc(rangeTime int64, delayTime int64) []datapoint {
	return retrieve("ori_acc", db.project, rangeTime, delayTime)
}

func (db *CdnExporter) RetrieveBPS(rangeTime int64, delayTime int64) []datapoint {
	return retrieve("BPS", db.project, rangeTime, delayTime)
}
