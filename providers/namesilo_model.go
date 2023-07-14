package providers

type ResourceRecord struct {
	RecordID string `xml:"record_id"`
	Type     string `xml:"type"`
	Host     string `xml:"host"`
	Value    string `xml:"value"`
	TTL      int    `xml:"ttl"`
	Distance int    `xml:"distance"`
}
