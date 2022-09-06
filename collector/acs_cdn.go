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

func (db *CdnExporter) RetrieveHitRate() []datapoint {
	return retrieve("hitRate", db.project)
}

func (db *CdnExporter) RetrieveOriBps() []datapoint {
	return retrieve("ori_bps", db.project)
}

func (db *CdnExporter) RetrieveL1Acc() []datapoint {
	return retrieve("l1_acc", db.project)
}

func (db *CdnExporter) RetrieveOriAcc() []datapoint {
	return retrieve("ori_acc", db.project)
}

func (db *CdnExporter) RetrieveBPS() []datapoint {
	return retrieve("BPS", db.project)
}

func (db *CdnExporter) RetrieveOriStatusRatio() []statusPoint {
	statusMetrics := map[string]string{
		"code1xx": "1xx",
		"code2xx": "2xx",
		"code3xx": "3xx",
		"code4xx": "4xx",
		"code_ratio_499": "499",
		"code5xx": "5xx",
	}
	var response []statusPoint
	for metric, status := range statusMetrics {
		dataPoints := retrieve(metric, db.project)
		for _, point := range dataPoints{
			// code1xx、code2xx、code3xx 只有 Maximum 值
			if metric == "code1xx" || metric == "code2xx" || metric == "code3xx" {
				point.Average = point.Maximum
			}
			response = append(response, statusPoint{point, status})
		}
	}
	return response
}

func (db *CdnExporter) RetrieveStatusRatio() []statusPoint {
	backSourceStatusMetrics := map[string]string{
		"ori_code_ratio_1xx": "1xx",
		"ori_code_ratio_2xx": "2xx",
		"ori_code_ratio_3xx": "3xx",
		"ori_code_ratio_499": "499",
		"ori_code_ratio_4xx": "4xx",
		"ori_code_ratio_5xx": "5xx",
	}
	var response []statusPoint
	for metric, status := range backSourceStatusMetrics {
		dataPoints := retrieve(metric, db.project)
		for _, point := range dataPoints{
			response = append(response, statusPoint{point, status})
		}
	}
	return response
}
