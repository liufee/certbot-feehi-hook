package main

import (
	"certbot-feehi-hook/providers"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

const ver = "1.0.0"

var (
	h = flag.Bool("h", false, "This help")
	v = flag.Bool("v", false, "Show certbot-feehi-hook version")

	providerType = flag.String("type", "", "Your dns provider. Current support aliyun")
	action       = flag.String("action", "", "add or delete. manual-auth-hook use add, manual-cleanup-hook use delete")
	aliyunKey    = flag.String("ali_AccessKey_ID", "", "aliyun access id")
	aliyunSecret = flag.String("ali_Access_Key_Secret", "", "aliyun access key secret")
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
	log.Println("Auto dns record manage start")
	parseCertbotPassedValue()
	var p providers.Provider
	log.Printf("Your dns provider is %s \n", *providerType)
	switch *providerType {
	case "aliyun":
		p = providers.NewAliyun(domain, *aliyunKey, *aliyunSecret)
	default:
		if *providerType == "" {
			panic("--type must be your dns provider type, such as aliyun qcloud")
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
		log.Println("Time sleep 20s for record effective")
		time.Sleep(time.Second * 20)
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
