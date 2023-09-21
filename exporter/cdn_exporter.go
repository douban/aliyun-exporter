package exporter

import (
	"aliyun-exporter/collector"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"sync"
)

const (
	cdnnamespace = "aliyun"
)

type CdnExporter struct {
	client    *cms.Client
	cdnClient *cdn.Client
	rangeTime int64
	delayTime int64
	domains   []string

	fluxHitRate           *prometheus.Desc
	hitRate               *prometheus.Desc
	backSourceBps         *prometheus.Desc
	BPS                   *prometheus.Desc
	l1Acc                 *prometheus.Desc
	backSourceAcc         *prometheus.Desc
	backSourceStatusRatio *prometheus.Desc
	statusRatio           *prometheus.Desc
}

// CdnCloudExporter 实例化
func CdnCloudExporter(cmsClient *cms.Client, cdnClient *cdn.Client, rangeTime int64, delayTime int64) *CdnExporter {
	domains := collector.GetDomains(*cdnClient, "online")
	return &CdnExporter{
		client:    cmsClient,
		cdnClient: cdnClient,
		rangeTime: rangeTime,
		delayTime: delayTime,
		domains:   domains,

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

// Describe 导出
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

// Collect 收集
func (e *CdnExporter) Collect(ch chan<- prometheus.Metric) {
	cdnDashboard := collector.NewCdnExporter(e.client)
	var wg sync.WaitGroup

	//domains := []string{"vt3.doubanio.com"}
	for _, domain := range e.domains {
		domain := domain
		reqHitRate := collector.GetReqHitRate(*e.cdnClient, domain, e.rangeTime, e.delayTime)
		// 去除掉数据量少的域名
		if reqHitRate < 10 {
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			e.hitRate,
			prometheus.GaugeValue,
			reqHitRate,
			domain,
		)
		wg.Add(1)
		go func() {
			defer wg.Done()
			statusProportion := collector.GetStatusCode(*e.cdnClient, domain, e.rangeTime, e.delayTime)
			for status, proportion := range statusProportion {
				proportion, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", proportion), 64)
				ch <- prometheus.MustNewConstMetric(
					e.statusRatio,
					prometheus.GaugeValue,
					proportion,
					domain,
					status,
				)
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			resourceStatus := collector.GetResourceStatusCode(*e.cdnClient, domain, e.rangeTime, e.delayTime)
			for status, proportion := range resourceStatus {
				proportion, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", proportion), 64)
				ch <- prometheus.MustNewConstMetric(
					e.backSourceStatusRatio,
					prometheus.GaugeValue,
					proportion,
					domain,
					status,
				)
			}
		}()

	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, point := range cdnDashboard.RetrieveHitRate(e.rangeTime, e.delayTime) {
			ch <- prometheus.MustNewConstMetric(
				e.fluxHitRate,
				prometheus.GaugeValue,
				point.Average,
				point.InstanceId,
			)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, point := range cdnDashboard.RetrieveOriBps(e.rangeTime, e.delayTime) {
			ch <- prometheus.MustNewConstMetric(
				e.backSourceBps,
				prometheus.GaugeValue,
				point.Average/1000/1000,
				point.InstanceId,
			)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, point := range cdnDashboard.RetrieveL1Acc(e.rangeTime, e.delayTime) {
			ch <- prometheus.MustNewConstMetric(
				e.l1Acc,
				prometheus.GaugeValue,
				point.Average,
				point.InstanceId,
			)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, point := range cdnDashboard.RetrieveOriAcc(e.rangeTime, e.delayTime) {
			ch <- prometheus.MustNewConstMetric(
				e.backSourceAcc,
				prometheus.GaugeValue,
				point.Average,
				point.InstanceId,
			)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, point := range cdnDashboard.RetrieveBPS(e.rangeTime, e.delayTime) {
			ch <- prometheus.MustNewConstMetric(
				e.BPS,
				prometheus.GaugeValue,
				point.Average/1000/1000,
				point.InstanceId,
			)
		}
	}()
	wg.Wait()
}
