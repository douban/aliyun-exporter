package collector

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"log"
)

type datapoint struct {
	Average    float64 `json:"Average"`
	Maximum    float64 `json:"Maximum"`
	Minimum    float64 `json:"Minimum"`
	Value      float64 `json:"Value"`
	InstanceId string  `json:"instanceId"`
	Timestamp  int64   `json:"timestamp"`
}

type GetResponseFunc func(client *cms.Client, request *cms.DescribeMetricLastRequest) (string, error)

type Project struct {
	client      *cms.Client
	getResponse GetResponseFunc
	Namespace   string
}

func defaultGetResponseFunc(client *cms.Client, request *cms.DescribeMetricLastRequest) (string, error) {
	response, err := client.DescribeMetricLast(request)
	if err != nil {
		log.Fatal(err)
		return "", err
	} else {
		return response.Datapoints, nil
	}

}

func retrieve(metric string, p Project) []datapoint {
	request := cms.CreateDescribeMetricLastRequest()
	request.Namespace = p.Namespace
	request.MetricName = metric
	requestsStats.Inc()

	datapoints := make([]datapoint, 0)

	getResponseFunc := p.getResponse
	if getResponseFunc == nil {
		getResponseFunc = defaultGetResponseFunc
	}
	//log.Println("\ntest retrieve\n")
	source, err := getResponseFunc(p.client, request)
	if err != nil {
		responseError.Inc()
		log.Fatal(err)
	} else if err := json.Unmarshal([]byte(source), &datapoints); err != nil {
		responseFormatError.Inc()
		log.Fatal(err)
	}
	return datapoints
}
