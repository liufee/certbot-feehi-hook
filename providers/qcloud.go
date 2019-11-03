package providers

import (
	"certbot-feehi-hook/utils"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const qcloudBaseUrl = "https://cns.api.qcloud.com/v2/index.php"

func NewQcloud(domain, qcloudAppkey, qcloudAppsecret string) *qcloud {
	if qcloudAppkey == "" {
		panic("qcloud_SecretId must be passed")
	}
	if qcloudAppsecret == "" {
		panic("qcloud_SecretKey must be passed")
	}
	domain, levelsDomainName := utils.ParseDomain(domain)
	return &qcloud{
		DomainName:       domain,
		LevelsDomainName: levelsDomainName,
		AppKey:           qcloudAppkey,
		AppSecret:        qcloudAppsecret,
	}
}

type qcloud struct {
	AppKey           string
	AppSecret        string
	DomainName       string
	LevelsDomainName string
}

func (q qcloud) ResolveDomainName(dnsType string, RR string, value string) (bool, error) {
	var result bool
	var err error

	if q.LevelsDomainName != "" {
		RR += "." + q.LevelsDomainName
	}

	record, err := q.GetRecordByTypeAndPr(dnsType, RR)

	if err != nil {
		return false, err
	}
	if record == nil { //need new add
		log.Println("add new record")
		result, err = q.AddRecord(dnsType, RR, value)
	} else {
		log.Println(" record already exist, do update")
		result, err = q.UpdateRecord(dnsType, RR, value, strconv.Itoa(record.Id))
	}
	if err != nil {
		return false, err
	}
	return result, nil
}

func (q *qcloud) DeleteResolveDomainName(dnsType string, RR string) (bool, error) {
	if q.LevelsDomainName != "" {
		RR += "." + q.LevelsDomainName
	}
	record, err := q.GetRecordByTypeAndPr(dnsType, RR)
	if err != nil {
		return false, err
	}
	if record != nil {
		result, err := q.DeleteRecord(strconv.Itoa(record.Id))
		if err != nil {
			return false, err
		}
		return result, nil
	}
	log.Println("Dns record not exists")
	return false, nil
}

func (q *qcloud) GetRecords() (qcloudGetRecordsResult, error) {
	var result qcloudGetRecordsResult
	businessParams := map[string]string{
		//"subDomain": RR,
		"length": "100",
	}
	body, err := q.Do("RecordList", "GET", businessParams)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}
	if result.Code != 0 {
		return result, errors.New(result.Message)
	}
	return result, nil
}

func (q *qcloud) AddRecord(dnsType string, RR string, value string) (bool, error) {
	businessParams := map[string]string{
		"recordType": dnsType, //TXT
		"recordLine": "默认",    //www
		"subDomain":  RR,
		"value":      value, //8.8.8.8
	}
	body, err := q.Do("RecordCreate", "GET", businessParams)
	if err != nil {
		return false, err
	}
	var b qcloudResult
	err = json.Unmarshal(body, &b)
	if err != nil {
		return false, err
	}
	if b.Code != 0 {
		return false, errors.New(b.Message)
	}
	return true, nil
}

func (q *qcloud) UpdateRecord(dnsType string, RR string, value string, recordId string) (bool, error) {
	businessParams := map[string]string{
		"recordId":   recordId,
		"recordType": dnsType, //TXT
		"recordLine": "默认",    //www
		"subDomain":  RR,
		"value":      value, //8.8.8.8
	}
	body, err := q.Do("RecordModify", "GET", businessParams)
	if err != nil {
		return false, err
	}
	var b qcloudResult
	err = json.Unmarshal(body, &b)
	if err != nil {
		return false, err
	}
	if b.Code != 0 {
		return false, errors.New(b.Message)
	}
	return true, nil
}

func (q *qcloud) DeleteRecord(recordId string) (bool, error) {
	businessParams := map[string]string{
		"recordId": recordId,
	}
	body, err := q.Do("RecordDelete", "GET", businessParams)
	if err != nil {
		return false, err
	}
	var b qcloudResult
	err = json.Unmarshal(body, &b)
	if err != nil {
		return false, err
	}
	if b.Code != 0 {
		return false, errors.New(b.Message)
	}
	return true, nil
}

func (q *qcloud) Do(action string, httpRequestMethod string, businessParams map[string]string) ([]byte, error) {
	params := map[string]string{
		"Action":    action,
		"Nonce":     utils.GetRandomString(6),
		"Timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		"SecretId":  q.AppKey,
		"domain":    q.DomainName,
	}
	for k, v := range businessParams {
		params[k] = v
	}
	signValue := q.signature(action, httpRequestMethod, params)
	u := url.Values{}
	for k, v := range params {
		u.Add(k, v)
	}
	u.Add("Signature", signValue)
	return utils.HttpGet(qcloudBaseUrl, u)
}

func (q *qcloud) signature(action string, httpRequestMethod string, params map[string]string) string {
	var keys = func() []string {
		var keys []string
		for key := range params {
			keys = append(keys, key)
		}
		return keys
	}()
	sort.Strings(keys)
	var requestParamString string
	for _, key := range keys {
		requestParamString += "&" + strings.Replace(key, "_", ".", -1) + "=" + params[key]
	}
	requestParamString = strings.Trim(requestParamString, "&")
	u := strings.Replace(qcloudBaseUrl, "https://", "", 1)
	stringToSign := httpRequestMethod + u + "?" + requestParamString
	key := []byte(q.AppSecret)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(stringToSign))
	b := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(b)
}

func (q *qcloud) GetRecordByTypeAndPr(tp, RR string) (*qcloudRecord, error) {
	var record qcloudRecord
	recordsResult, err := q.GetRecords()
	if err != nil {
		return &record, err
	}
	for _, v := range recordsResult.Data.Records {
		if v.Type == tp && v.Name == RR {
			return &v, nil
		}
	}
	return nil, nil
}
