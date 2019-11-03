package utils

import (
	"log"
	"strings"
)

var threePointDomains = []string{
	".gov.cn",
	".net.cn",
	".com.cn",
	".org.cn",
	".co.uk",
}

func ParseDomain(domain string) (rootDomain, levelsDomain string) {
	rootDomain = ""
	levelsDomain = ""
	dotNum := strings.Count(domain, ".")
	if dotNum < 1 {
		log.Fatalln("Domain format error")
	}
	if dotNum == 1 {
		rootDomain = domain
		return
	}

	var t = 0
	var dotTime = 2
	for _, item := range threePointDomains {
		if strings.HasSuffix(domain, item) {
			dotTime = 3
			break
		}
	}
	for i := len(domain) - 1; i > -1; i-- {
		char := domain[i]
		if char == byte('.') {
			t++
		}
		if t < dotTime {
			rootDomain += string(char)
		} else {
			levelsDomain += string(char)
		}
	}
	rootDomain = reverseString(rootDomain)
	levelsDomain = reverseString(levelsDomain)
	levelsDomain = strings.TrimRight(levelsDomain, ".")
	return
}

func reverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}
