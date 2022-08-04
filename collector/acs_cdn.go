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

func (db *CdnExporter) RetrievehitRate() []datapoint {
	return retrieve("hitRate", db.project)
}

func (db *CdnExporter) RetrieveoriBps() []datapoint {
	return retrieve("ori_bps", db.project)
}

func (db *CdnExporter) Retrievel1Acc() []datapoint {
	return retrieve("l1_acc", db.project)
}

func (db *CdnExporter) RetrieveoriAcc() []datapoint {
	return retrieve("ori_acc", db.project)
}

func (db *CdnExporter) RetrieveBPS() []datapoint {
	return retrieve("BPS", db.project)
}

func (db *CdnExporter) Retrieveoricode4xxRatio() []datapoint {
	return retrieve("ori_code_ratio_4xx", db.project)
}

func (db *CdnExporter) Retrieveoricode5xxRatio() []datapoint {
	return retrieve("ori_code_ratio_5xx", db.project)
}
