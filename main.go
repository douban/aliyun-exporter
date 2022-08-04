package main

import (
	"aliyun-exporter/exporter"
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var config struct {
	accessKeyId     string
	accessKeySecret string
	regionId        string
	host            string
	port            int
	service         string
}

func main() {
	flag.StringVar(&(config.accessKeyId), "id", os.Getenv("ACCESS_KEY_ID"), "阿里云AccessKey ID")
	flag.StringVar(&(config.accessKeySecret), "secret", os.Getenv("ACCESS_KEY_SECRET"), "阿里云AccessKey Secret")
	flag.StringVar(&(config.regionId), "region", os.Getenv("REGIONID"), "阿里云Region ID")
	flag.StringVar(&(config.host), "host", "0.0.0.0", "服务监听地址")
	flag.IntVar(&(config.port), "port", 9180, "服务监听端口")
	flag.StringVar(&(config.service), "service", "acs_cdn", "输出Metrics的服务，默认为全部")
	flag.Parse()
	serviceArr := strings.Split(config.service, ",")
	for _, ae := range serviceArr {
		switch ae {
		case "acs_cdn":
			cdn := exporter.CdnCloudExporter(newCdnClient())
			prometheus.MustRegister(cdn)
		default:
			log.Println("暂不支持该服务，请根据提示选择服务。")
		}
	}

	listenAddress := net.JoinHostPort(config.host, strconv.Itoa(config.port))
	log.Println(listenAddress)
	log.Println("Running on", listenAddress)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
