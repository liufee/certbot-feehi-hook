package providers

import (
	"certbot-feehi-hook/utils"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const aliBaseUrl = "http://alidns.aliyuncs.com"

func NewAliyun(domain, aliKey, aliSecret string) *Aliyun {
	if aliKey == "" {
		panic("ali_AccessKey_ID must be passed")
	}
	if aliSecret == "" {
		panic("ali_Access_Key_Secret must be passed")
	}
	domain, levelsDomainName := utils.ParseDomain(domain)
	return &Aliyun{
		DomainName:       domain,
		LevelsDomainName: levelsDomainName,
		AppKey:           aliKey,
		AppSecret:        aliSecret,
	}
}

type Aliyun struct {
	AppKey           string
	AppSecret        string
	DomainName       string
	LevelsDomainName string
}

func (a *Aliyun) ResolveDomainName(dnsType string, RR string, value string) (bool, error) {
	var result bool
	var err error
	if a.LevelsDomainName != "" {
		RR += "." + a.LevelsDomainName
	}
	record, err := a.GetRecordByTypeAndPr(dnsType, RR)
	if err != nil {
		return false, err
	}
	if record == nil { //need new add
		log.Println("add new record")
		result, err = a.AddRecord(dnsType, RR, value)
	} else {
		log.Println(" record already exist, do update")
		result, err = a.UpdateRecord(dnsType, RR, value, record.RecordId)
	}
	if err != nil {
		return false, err
	}
	return result, nil
}

func (a *Aliyun) DeleteResolveDomainName(dnsType string, RR string) (bool, error) {
	if a.LevelsDomainName != "" {
		RR += "." + a.LevelsDomainName
	}
	record, err := a.GetRecordByTypeAndPr(dnsType, RR)
	if err != nil {
		return false, err
	}
	if record != nil {
		result, err := a.DeleteRecord(record.RecordId)
		if err != nil {
			return false, err
		}
		return result, nil
	}
	log.Println("Dns record not exists")
	return false, nil
}

func (a *Aliyun) GetRecords() (aliyunGetRecordsResult, error) {
	var result aliyunGetRecordsResult
	businessParams := map[string]string{
		"Action":   "DescribeDomainRecords",
		"PageSize": "100",
	}
	body, err := a.Do(businessParams)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (a *Aliyun) AddRecord(dnsType string, RR string, value string) (bool, error) {
	businessParams := map[string]string{
		"Action": "AddDomainRecord",
		"Type":   dnsType, //TXT
		"RR":     RR,      //www
		"Value":  value,   //8.8.8.8
	}
	body, err := a.Do(businessParams)
	if err != nil {
		return false, err
	}
	var b aliyunResult
	err = json.Unmarshal(body, &b)
	if err != nil {
		return false, err
	}
	if b.Code != "" {
		return false, errors.New(b.Code + ":" + b.Message)
	}
	return true, nil
}

func (a *Aliyun) UpdateRecord(dnsType string, RR string, value string, recordId string) (bool, error) {
	businessParams := map[string]string{
		"Action":   "UpdateDomainRecord",
		"Type":     dnsType, //TXT
		"RR":       RR,      //www
		"Value":    value,   //8.8.8.8
		"RecordId": recordId,
	}
	body, err := a.Do(businessParams)
	if err != nil {
		return false, err
	}
	var b aliyunResult
	err = json.Unmarshal(body, &b)
	if err != nil {
		return false, err
	}
	if b.Code != "" {
		return false, errors.New(b.Code + ":" + b.Message)
	}
	return true, nil
}

func (a *Aliyun) DeleteRecord(recordId string) (bool, error) {
	businessParams := map[string]string{
		"Action":   "DeleteDomainRecord",
		"RecordId": recordId,
	}
	body, err := a.Do(businessParams)
	if err != nil {
		return false, err
	}
	var b aliyunResult
	err = json.Unmarshal(body, &b)
	if err != nil {
		return false, err
	}
	if b.Code != "" {
		return false, errors.New(b.Code + ":" + b.Message)
	}
	return true, nil
}

func (a *Aliyun) Do(businessParams map[string]string) ([]byte, error) {
	md5.Sum([]byte(strconv.Itoa(rand.Int())))
	params := map[string]string{
		"DomainName":       a.DomainName,
		"Format":           "JSON",
		"Version":          "2015-01-09",
		"AccessKeyId":      a.AppKey,
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05-0700"),
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   utils.GetRandomString(16),
	}

	for k, v := range businessParams {
		params[k] = v
	}
	signValue := a.signature(params, "GET")
	u := url.Values{}
	for k, v := range params {
		u.Add(k, v)
	}
	u.Add("Signature", signValue)
	return utils.HttpGet(aliBaseUrl, u)
}

func (a *Aliyun) signature(params map[string]string, method string) string {
	var keys = func() []string {
		var keys []string
		for key := range params {
			keys = append(keys, key)
		}
		return keys
	}()
	sort.Strings(keys)
	stringToSign := strings.ToUpper(method) + "&" + url.QueryEscape("/") + "&"

	var requestParamString = ""

	for _, key := range keys {
		requestParamString += "&" + url.QueryEscape(key) + "=" + url.QueryEscape(params[key])
	}
	requestParamString = strings.Trim(requestParamString, "&")
	stringToSign = stringToSign + url.QueryEscape(requestParamString)

	key := []byte(a.AppSecret + "&")
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(stringToSign))
	b := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(b)
}

func (a *Aliyun) GetRecordByTypeAndPr(tp, RR string) (*aliyunRecord, error) {
	var record aliyunRecord
	recordsResult, err := a.GetRecords()
	if err != nil {
		return &record, err
	}
	for _, v := range recordsResult.DomainRecords.Record {
		if v.Type == tp && v.RR == RR {
			return &v, nil
		}
	}
	return nil, nil
}
