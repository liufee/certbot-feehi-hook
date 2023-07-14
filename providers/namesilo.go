package providers

import (
	"certbot-feehi-hook/utils"
	"encoding/xml"
	"fmt"
	"net/url"
)

func NewNamesilo(domain, apiKey string) *Namesilo {
	return &Namesilo{
		apiKey: apiKey,
		domain: domain,
	}
}

var namesiloURL = "https://www.namesilo.com/api/"

type Namesilo struct {
	apiKey string
	domain string
}

func (n *Namesilo) ResolveDomainName(dnsType string, pr string, value string) (bool, error) {
	record, err := n.GetRecordByTypeAndPr(dnsType, pr)
	if err != nil {
		return false, err
	}
	if record != nil { // delete
		_, err = n.DeleteRecord(record.RecordID)
		if err != nil {
			return false, err
		}
	}
	// add
	_, err = n.AddRecord(dnsType, pr, value)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (n *Namesilo) AddRecord(dnsType string, pr string, value string) (bool, error) {
	query := url.Values{}
	query.Add("rrtype", dnsType)
	query.Add("domain", n.domain)
	query.Add("rrhost", pr)
	query.Add("rrvalue", value)
	query.Add("rrttl", "3600")
	res, err := n.do("dnsAddRecord", query)
	if err != nil {
		panic(err)
	}
	var response struct {
		Reply struct {
			Code     int    `xml:"code"`
			Detail   string `xml:"detail"`
			RecordID string `xml:"record_id"`
		} `xml:"reply"`
	}
	err = xml.Unmarshal(res, &response)
	if err != nil {
		return false, err
	}
	if response.Reply.Code != 300 {
		return false, fmt.Errorf(response.Reply.Detail)
	}
	return true, nil
}

func (n *Namesilo) GetRecords() ([]ResourceRecord, error) {
	query := url.Values{}
	query.Add("domain", n.domain)
	res, err := n.do("dnsListRecords", query)
	if err != nil {
		panic(err)
	}
	var response struct {
		Reply struct {
			Code            int              `xml:"code"`
			Detail          string           `xml:"detail"`
			ResourceRecords []ResourceRecord `xml:"resource_record"`
		} `xml:"reply"`
	}
	err = xml.Unmarshal(res, &response)
	if err != nil {
		return nil, err
	}
	if response.Reply.Code != 300 {
		return nil, fmt.Errorf(response.Reply.Detail)
	}
	return response.Reply.ResourceRecords, nil
}
func (n *Namesilo) DeleteResolveDomainName(dnsType string, pr string) (bool, error) {
	record, err := n.GetRecordByTypeAndPr(dnsType, pr)
	if err != nil {
		return false, err
	}
	if record != nil {
		b, err := n.DeleteRecord(record.RecordID)
		if err != nil {
			return false, err
		}
		return b, nil
	}
	return true, nil
}

func (n *Namesilo) GetRecordByTypeAndPr(tp, RR string) (*ResourceRecord, error) {
	var record ResourceRecord
	recordsResult, err := n.GetRecords()
	if err != nil {
		return &record, err
	}
	for _, v := range recordsResult {
		if v.Type == tp && v.Host == RR+"."+n.domain {
			return &v, nil
		}
	}
	return nil, nil
}

func (n *Namesilo) DeleteRecord(rrid string) (bool, error) {
	query := url.Values{}
	query.Add("rrid", rrid)
	query.Add("domain", n.domain)
	res, err := n.do("dnsDeleteRecord", query)
	if err != nil {
		return false, err
	}
	var response struct {
		Reply struct {
			Code   int    `xml:"code"`
			Detail string `xml:"detail"`
		} `xml:"reply"`
	}
	err = xml.Unmarshal(res, &response)
	if err != nil {
		return false, err
	}
	if response.Reply.Code != 300 {
		return false, fmt.Errorf(response.Reply.Detail)
	}
	return true, nil
}

func (n *Namesilo) do(operation string, query url.Values) ([]byte, error) {
	query.Add("version", "1")
	query.Add("type", "xml")
	query.Add("key", n.apiKey)
	return utils.HttpGet(namesiloURL+operation, query)
}
