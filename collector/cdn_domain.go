package collector

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"log"
	"strconv"
	"time"
)


func GetDomains(cdnClient cdn.Client, status string) []string {
	var domains []string
	req := cdn.CreateDescribeUserDomainsRequest()
	req.DomainStatus = status

	response, err := cdnClient.DescribeUserDomains(req)
	if err != nil {
		log.Fatal("Error response from Aliyun:", err)
	}
	for _, res := range response.Domains.PageData{
		domains = append(domains, res.DomainName)
	}
	return domains
}

func GetReqHitRate(cdnClient cdn.Client, domainName string, rangeTime int64, delayTime int64) float64 {
	var reqHitRateSum float64
	req := cdn.CreateDescribeDomainReqHitRateDataRequest()
	req.DomainName = domainName
	// 靠近当前时间数据会不太准确，调整前移抓取时间
	req.StartTime = time.Now().UTC().Add(-time.Second * time.Duration(rangeTime)).Format(time.RFC3339)
	req.EndTime = time.Now().UTC().Add(-time.Second * time.Duration(delayTime)).Format(time.RFC3339)
	response, err := cdnClient.DescribeDomainReqHitRateData(req)
	if err != nil {
		log.Fatal("Error response from Aliyun:", err)
	}
	for _, reqHitRate := range response.ReqHitRateInterval.DataModule{
		value, _ :=strconv.ParseFloat(reqHitRate.Value, 64)
		reqHitRateSum += value
	}
	reqHitRateAverage, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", reqHitRateSum / float64(len(response.ReqHitRateInterval.DataModule))), 64)
	return reqHitRateAverage
}
