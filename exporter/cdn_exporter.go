package exporter

import (
	"aliyun-exporter/collector"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	cdnnamespace = "aliyun"
)

type CdnExporter struct {
	client *cms.Client
	cdnClient *cdn.Client
	rangeTime            int64
	delayTime            int64

	fluxHitRate            *prometheus.Desc
	hitRate                *prometheus.Desc
	backSourceBps          *prometheus.Desc
	BPS                    *prometheus.Desc
	l1Acc                  *prometheus.Desc
	backSourceAcc          *prometheus.Desc
	backSourceStatusRatio  *prometheus.Desc
	statusRatio            *prometheus.Desc
}

//实例化
func CdnCloudExporter(cmsClient *cms.Client, cdnClient *cdn.Client, rangeTime int64, delayTime int64) *CdnExporter {
	return &CdnExporter{
		client: cmsClient,
		cdnClient: cdnClient,
		rangeTime: rangeTime,
		delayTime: delayTime,

		fluxHitRate: prometheus.NewDesc(
			prometheus.BuildFQName(cdnnamespace, "cdn", "flux_hit_rate"),
			"边缘字节命中率(%)",
			[]string{
				"instanceId",
			},
			nil,
		),
		hitRate: prometheus.NewDesc(
			prometheus.BuildFQName(cdnnamespace, "cdn", "hit_rate"),
			"请求命中率(%)",
			[]string{
				"instanceId",
			},
			nil,
		),
		BPS: prometheus.NewDesc(
			prometheus.BuildFQName(cdnnamespace, "cdn", "bandwidth"),
			"边缘网络带宽(Mbps)",
			[]string{
				"instanceId",
			},
			nil,
		),
		backSourceBps: prometheus.NewDesc(
			prometheus.BuildFQName(cdnnamespace, "cdn", "ori_bandwidth"),
			"回源网络带宽(Mbps)",
			[]string{
				"instanceId",
			},
			nil,
		),

		l1Acc: prometheus.NewDesc(
			prometheus.BuildFQName(cdnnamespace, "cdn", "l1_acc"),
			"边缘累加请求数(Count)",
			[]string{
				"instanceId",
			},
			nil,
		),

		backSourceAcc: prometheus.NewDesc(
			prometheus.BuildFQName(cdnnamespace, "cdn", "ori_acc"),
			"回源累加请求数(Count)",
			[]string{
				"instanceId",
			},
			nil,
		),

		backSourceStatusRatio: prometheus.NewDesc(
			prometheus.BuildFQName(cdnnamespace, "cdn", "ori_status_ratio"),
			"回源状态码占比(%)",
			[]string{
				"instanceId",
				"status",
			},
			nil,
		),

		statusRatio: prometheus.NewDesc(
			prometheus.BuildFQName(cdnnamespace, "cdn", "status_ratio"),
			"状态码占比(%)",
			[]string{
				"instanceId",
				"status",
			},
			nil,
		),
	}
}

//导出
func (e *CdnExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.fluxHitRate
	ch <- e.hitRate
	ch <- e.backSourceBps
	ch <- e.BPS
	ch <- e.backSourceAcc
	ch <- e.l1Acc
	ch <- e.backSourceStatusRatio
	ch <- e.statusRatio
}

//收集
func (e *CdnExporter) Collect(ch chan<- prometheus.Metric) {
	cdnDashboard := collector.NewCdnExporter(e.client)

	domains := collector.GetDomains(*e.cdnClient, "online")
	//domains := []string{"vt3.doubanio.com"}
	for _, domain := range domains {
		reqHitRate := collector.GetReqHitRate(*e.cdnClient, domain, e.rangeTime, e.delayTime)
		// 去除掉数据量少的域名
		if reqHitRate < 10 {
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			e.hitRate,
			prometheus.GaugeValue,
			float64(reqHitRate),
			domain,
		)
	}

	for _, point := range cdnDashboard.RetrieveHitRate() {
		ch <- prometheus.MustNewConstMetric(
			e.fluxHitRate,
			prometheus.GaugeValue,
			float64(point.Average),
			point.InstanceId,
		)
	}

	for _, point := range cdnDashboard.RetrieveOriBps() {
		ch <- prometheus.MustNewConstMetric(
			e.backSourceBps,
			prometheus.GaugeValue,
			float64(point.Average / 1000 / 1000),
			point.InstanceId,
		)
	}
	for _, point := range cdnDashboard.RetrieveL1Acc() {
		ch <- prometheus.MustNewConstMetric(
			e.l1Acc,
			prometheus.GaugeValue,
			float64(point.Average),
			point.InstanceId,
		)
	}
	for _, point := range cdnDashboard.RetrieveOriAcc() {
		ch <- prometheus.MustNewConstMetric(
			e.backSourceAcc,
			prometheus.GaugeValue,
			float64(point.Average),
			point.InstanceId,
		)
	}
	for _, point := range cdnDashboard.RetrieveBPS() {
		ch <- prometheus.MustNewConstMetric(
			e.BPS,
			prometheus.GaugeValue,
			float64(point.Average / 1000 / 1000),
			point.InstanceId,
		)
	}
	for _, point := range cdnDashboard.RetrieveOriStatusRatio() {
		ch <- prometheus.MustNewConstMetric(
			e.backSourceStatusRatio,
			prometheus.GaugeValue,
			float64(point.Average),
			point.InstanceId,
			point.Status,
		)
	}
	for _, point := range cdnDashboard.RetrieveStatusRatio() {
		ch <- prometheus.MustNewConstMetric(
			e.statusRatio,
			prometheus.GaugeValue,
			float64(point.Average),
			point.InstanceId,
			point.Status,
		)
	}
}
