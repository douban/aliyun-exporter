package collector

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"log"
	"strconv"
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

func GetReqHitRate(cdnClient cdn.Client, domainName string) float64 {
	var reqHitRateSum float64
	req := cdn.CreateDescribeDomainReqHitRateDataRequest()
	req.DomainName = domainName

	response, err := cdnClient.DescribeDomainReqHitRateData(req)
	if err != nil {
		log.Fatal("Error response from Aliyun:", err)
	}
	// 最后一条记录数据不准确，影响平均值计算，去除后再计算
	validResponse := response.ReqHitRateInterval.DataModule[:len(response.ReqHitRateInterval.DataModule)-1]
	for _, reqHitRate := range validResponse{
		value, _ :=strconv.ParseFloat(reqHitRate.Value, 64)
		reqHitRateSum += value
	}
	reqHitRateAverage, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", reqHitRateSum / float64(len(validResponse))), 64)
	return reqHitRateAverage
}
