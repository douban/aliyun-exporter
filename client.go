package main

import (
	"log"
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
		log.Fatal("client init failed",err)
	}

	return cdnClient
}
