package main

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"log"
)

func CmsClient() *cms.Client {
	cmsClient, err := cms.NewClientWithAccessKey(
		config.regionId,
		config.accessKeyId,
		config.accessKeySecret,
	)
	//log.Println("testcms")
	if err != nil {
		log.Fatal("client init failed",err)
	}

	return cmsClient
}

func CdnClient () *cdn.Client {
	cdnClient, err := cdn.NewClientWithAccessKey(
		config.regionId,
		config.accessKeyId,
		config.accessKeySecret,
	)
	if err != nil {
		log.Fatal("client init failed",err)
	}

	return cdnClient
}
