package exporter

import (
	"aliyun-exporter/collector"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	cdnnamespace = "aliyun"
)

type CdnExporter struct {
	client *cms.Client

	cdnhitRate         *prometheus.Desc
	cdnoriBps          *prometheus.Desc
	cdnl1Acc           *prometheus.Desc
	cdnoriAcc          *prometheus.Desc
	cdnBPS             *prometheus.Desc
	cdnoricode4xxratio *prometheus.Desc
	cdnoricode5xxratio *prometheus.Desc
}

//实例化
func CdnCloudExporter(c *cms.Client) *CdnExporter {
	return &CdnExporter{
		client: c,

		cdnhitRate: prometheus.NewDesc(
			prometheus.BuildFQName(cdnnamespace, "cdn", "hit_rate"),
			"边缘字节命中率(%)",
			[]string{
				"instanceId",
			},
			nil,
		),
		cdnBPS: prometheus.NewDesc(
			prometheus.BuildFQName(cdnnamespace, "cdn", "BPS"),
			"边缘网络带宽(bit/s)",
			[]string{
				"instanceId",
			},
			nil,
		),
		cdnoriBps: prometheus.NewDesc(
			prometheus.BuildFQName(cdnnamespace, "cdn", "ori_bps"),
			"回源网络带宽(bit/s)",
			[]string{
				"instanceId",
			},
			nil,
		),

		cdnl1Acc: prometheus.NewDesc(
			prometheus.BuildFQName(cdnnamespace, "cdn", "l1_acc"),
			"边缘累加请求数(Count)",
			[]string{
				"instanceId",
			},
			nil,
		),

		cdnoriAcc: prometheus.NewDesc(
			prometheus.BuildFQName(cdnnamespace, "cdn", "ori_acc"),
			"回源累加请求数(Count)",
			[]string{
				"instanceId",
			},
			nil,
		),

		cdnoricode4xxratio: prometheus.NewDesc(
			prometheus.BuildFQName(cdnnamespace, "cdn", "ori_code_ratio_4xx"),
			"回源状态码4XX占比(%)",
			[]string{
				"instanceId",
			},
			nil,
		),

		cdnoricode5xxratio: prometheus.NewDesc(
			prometheus.BuildFQName(cdnnamespace, "cdn", "ori_code_ratio_5xx"),
			"回源状态码5XX占比(%)",
			[]string{
				"instanceId",
			},
			nil,
		),
	}
}

//导出
func (e *CdnExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.cdnhitRate
	ch <- e.cdnoriBps
	ch <- e.cdnl1Acc
	ch <- e.cdnoriAcc
	ch <- e.cdnBPS
	ch <- e.cdnoricode4xxratio
	ch <- e.cdnoricode5xxratio
}

//收集
func (e *CdnExporter) Collect(ch chan<- prometheus.Metric) {
	cdnDashboard := collector.NewCdnExporter(e.client)

	for _, point := range cdnDashboard.RetrievehitRate() {
		ch <- prometheus.MustNewConstMetric(
			e.cdnhitRate,
			prometheus.GaugeValue,
			float64(point.Average),
			point.InstanceId,
		)
	}
	for _, point := range cdnDashboard.RetrieveoriBps() {
		ch <- prometheus.MustNewConstMetric(
			e.cdnoriBps,
			prometheus.GaugeValue,
			float64(point.Average),
			point.InstanceId,
		)
	}
	for _, point := range cdnDashboard.Retrievel1Acc() {
		ch <- prometheus.MustNewConstMetric(
			e.cdnl1Acc,
			prometheus.GaugeValue,
			float64(point.Average),
			point.InstanceId,
		)
	}
	for _, point := range cdnDashboard.RetrieveoriAcc() {
		ch <- prometheus.MustNewConstMetric(
			e.cdnoriAcc,
			prometheus.GaugeValue,
			float64(point.Average),
			point.InstanceId,
		)
	}
	for _, point := range cdnDashboard.RetrieveBPS() {
		ch <- prometheus.MustNewConstMetric(
			e.cdnBPS,
			prometheus.GaugeValue,
			float64(point.Average),
			point.InstanceId,
		)
	}
	for _, point := range cdnDashboard.Retrieveoricode4xxRatio() {
		ch <- prometheus.MustNewConstMetric(
			e.cdnoricode4xxratio,
			prometheus.GaugeValue,
			float64(point.Average),
			point.InstanceId,
		)
	}
	for _, point := range cdnDashboard.Retrieveoricode5xxRatio() {
		ch <- prometheus.MustNewConstMetric(
			e.cdnoricode5xxratio,
			prometheus.GaugeValue,
			float64(point.Average),
			point.InstanceId,
		)
	}
}
