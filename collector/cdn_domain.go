package collector

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"log"
	"strconv"
	"strings"
	"time"
)


func GetDomains(cdnClient cdn.Client, status string) []string {
	var domains []string
	req := cdn.CreateDescribeUserDomainsRequest()
	req.DomainStatus = status
	req.Scheme = "https"

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

func GetStatusCode(cdnClient cdn.Client, domainName string, rangeTime int64, delayTime int64) map[string]float64 {
	req := cdn.CreateDescribeDomainHttpCodeDataRequest()
	req.DomainName = domainName
	req.StartTime = time.Now().UTC().Add(-time.Second * time.Duration(rangeTime)).Format(time.RFC3339)
	req.EndTime = time.Now().UTC().Add(-time.Second * time.Duration(delayTime)).Format(time.RFC3339)
	response, err := cdnClient.DescribeDomainHttpCodeData(req)
	if err != nil {
		log.Fatal("Error response from Aliyun:", err)
	}
	httpStatusCodes := map[string]float64{"200":0, "206":0, "301":0, "302":0, "304":0, "400":0, "403":0, "404":0, "500":0, "502":0, "503":0, "504":0}
	for _, v := range response.HttpCodeData.UsageData{
		for _, status := range v.Value.CodeProportionData{
			proportion, _ := strconv.ParseFloat(status.Proportion, 32)
			httpStatusCodes[status.Code] = httpStatusCodes[status.Code] + proportion
		}
	}
	for status, value := range httpStatusCodes {
		if !strings.HasSuffix(status, "x") {
			value = value / float64(len(response.HttpCodeData.UsageData))
			httpStatusCodes[status] = value

			if strings.HasPrefix(status, "2") {
				httpStatusCodes["2xx"] += value
			} else if strings.HasPrefix(status, "3") {
				httpStatusCodes["3xx"] += value
			} else if strings.HasPrefix(status, "4") {
				httpStatusCodes["4xx"] += value
			} else if strings.HasPrefix(status, "5") {
				httpStatusCodes["5xx"] += value
			}
		}
	}

	return httpStatusCodes
}

func GetResourceStatusCode(cdnClient cdn.Client, domainName string, rangeTime int64, delayTime int64) map[string]float64 {
	req := cdn.CreateDescribeDomainSrcHttpCodeDataRequest()
	req.DomainName = domainName
	req.StartTime = time.Now().UTC().Add(-time.Second * time.Duration(rangeTime)).Format(time.RFC3339)
	req.EndTime = time.Now().UTC().Add(-time.Second * time.Duration(delayTime)).Format(time.RFC3339)
	response, err := cdnClient.DescribeDomainSrcHttpCodeData(req)
	if err != nil {
		log.Fatal("Error response from Aliyun:", err)
	}
	resourceStatusCodes := map[string]float64{"200":0, "206":0, "301":0, "302":0, "304":0, "400":0, "403":0, "404":0, "500":0, "502":0, "503":0, "504":0}
	for _, v := range response.HttpCodeData.UsageData{
		for _, status := range v.Value.CodeProportionData{
			proportion, _ := strconv.ParseFloat(status.Proportion, 32)
			resourceStatusCodes[status.Code] = resourceStatusCodes[status.Code] + proportion
		}
	}
	for status, value := range resourceStatusCodes {
		if !strings.HasSuffix(status, "x") {
			value = value / float64(len(response.HttpCodeData.UsageData))
			resourceStatusCodes[status] = value

			if strings.HasPrefix(status, "2") {
				resourceStatusCodes["2xx"] += value
			} else if strings.HasPrefix(status, "3") {
				resourceStatusCodes["3xx"] += value
			} else if strings.HasPrefix(status, "4") {
				resourceStatusCodes["4xx"] += value
			} else if strings.HasPrefix(status, "5") {
				resourceStatusCodes["5xx"] += value
			}
		}
	}

	return resourceStatusCodes
}
