package main

import (
	"certbot-feehi-hook/providers"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const ver = "2.0.0"

var (
	h       = flag.Bool("h", false, "This help")
	v       = flag.Bool("v", false, "Show certbot-feehi-hook version")
	logFile = flag.String("log", "/tmp/cert-book-feehi-log.txt", "log file")

	providerType = flag.String("type", "", "Your dns provider. Current support aliyun, qcloud")
	action       = flag.String("action", "", "add or delete. manual-auth-hook use add, manual-cleanup-hook use delete")

	aliyunKey    = flag.String("ali_AccessKey_ID", "", "aliyun access id")
	aliyunSecret = flag.String("ali_Access_Key_Secret", "", "aliyun access key secret")

	qcloudKey    = flag.String("qcloud_SecretId", "", "qcloud SecretId")
	qcloudSecret = flag.String("qcloud_SecretKey", "", "qcloud SecretKey")

	namesiloApiKey = flag.String("namesilo_apikey", "", "namesilo ApiKey")
)

var validationKey string
var domain string

func main() {
	flag.Parse()

	if *h {
		flag.Usage()
		os.Exit(0)
	}
	if *v {
		fmt.Println("Ver", ver)
		os.Exit(0)
	}

	l, err := os.OpenFile(*logFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	log.SetOutput(l)

	log.Println("Auto dns record manage start")
	parseCertbotPassedValue()
	var p providers.Provider
	log.Printf("Your dns provider is %s \n", *providerType)
	switch *providerType {
	case "aliyun":
		p = providers.NewAliyun(domain, *aliyunKey, *aliyunSecret)
	case "qcloud":
		p = providers.NewQcloud(domain, *qcloudKey, *qcloudSecret)
	case "namesilo":
		p = providers.NewNamesilo(domain, *namesiloApiKey)
	default:
		if *providerType == "" {
			panic("--type must be your dns provider type, such as aliyun qcloud namesilo")
		}
		panic("Not support " + *providerType + " yet")
	}
	run(p)
}

func run(p providers.Provider) {
	log.Println("Current action is ", *action)
	switch *action {
	case "add":
		log.Println("Will add record TXT _acme-challenge ", validationKey)
		result, err := p.ResolveDomainName("TXT", "_acme-challenge", validationKey)
		if err != nil {
			log.Fatalf("Auto resolve record TXT _acme-challenge %s comes to error %s \n", validationKey, err)
		}
		if !result {
			log.Fatalln("Auto resolve txt record failed")
		}
		log.Println("Auto resolve txt record success")
		log.Println("Check if record effected")
		for {
			records, _ := net.LookupTXT("_acme-challenge." + domain)
			if len(records) > 0 {
				break
			}
			log.Println("Record not effected", records)
			time.Sleep(time.Second * 20)
		}
		log.Println("feehi Hook finish")
	case "delete":
		log.Println("Will delete record TXT _acme-challenge ")
		result, err := p.DeleteResolveDomainName("TXT", "_acme-challenge")
		if err != nil {
			log.Printf("Auto delete record TXT _acme-challenge %s comes to error %s \n", validationKey, err)
		}
		if result {
			log.Println("Auto delete txt record success")
		} else {
			log.Println("Auto delete txt record failed")
		}
	default:
		panic("action only support add or delete")
	}
	os.Exit(0)
}

func parseCertbotPassedValue() {
	validationKey = os.Getenv("CERTBOT_VALIDATION")
	if validationKey == "" {
		panic("Not get CERTBOT_VALIDATION")
	}
	domain = os.Getenv("CERTBOT_DOMAIN")
	if domain == "" {
		panic("Not get CERTBOT_DOMAIN")
	}
}
