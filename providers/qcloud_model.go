package providers

type qcloudResult struct {
	Code    int
	Message string
}

type qcloudRecord struct {
	Id         int
	Ttl        int
	Value      string
	Enabled    int
	Status     string
	UpdatedOn  string `json:"updated_on"`
	QProjectId int    `json:"q_project_id"`
	Name       string
	Line       string
	LineId     string `json:"line_id"`
	Type       string
	Remark     string
	Mx         int
	Hold       string
}

type qcloudGetRecordsResult struct {
	Code     int
	Message  string
	CodeDesc string
	Data     struct {
		Domain struct {
			Id         string
			Name       string
			Punycode   string
			Grade      string
			Owner      string
			ExtStatus  string `json:"ext_status"`
			Ttl        int
			MinTtl     int      `json:"min_ttl"`
			DnspodNs   []string `json:"dnspod_ns"`
			Status     string
			QProjectId int `json:"q_project_id"`
		}
		Info struct {
			SubDomains  string `json:"sub_domains"`
			RecordTotal string `json:"record_total"`
		}
		Records []qcloudRecord
	}
}
