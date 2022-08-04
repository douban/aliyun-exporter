package main

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
)

func newCdnClient() *cms.Client {
	cdnClient, err := cms.NewClientWithAccessKey(
		config.regionId,
		config.accessKeyId,
		config.accessKeySecret,
	)
	//log.Println("testcms")
	if err != nil {
		panic(err)
	}

	return cdnClient
}
