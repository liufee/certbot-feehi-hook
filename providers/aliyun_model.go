package providers

type aliyunGetRecordsResult struct {
	PageNumber    int
	TotalCount    int
	PageSize      int
	RequestId     string
	DomainRecords struct {
		Record []aliyunRecord
	}
}

type aliyunRecord struct {
	RR         string `json:"RR"`
	Status     string
	Value      string
	Weight     int
	RecordId   string
	Type       string
	DomainName string
	Locked     bool
	Line       string
	TTL        int
}
