package entity

type ProxyItem struct {
	Ip       string
	Port     int
	OutIP    string
	Country  string
	ISP      string
	Timezone int
	Alive    int
	Status   int
}

type Filters struct {
	AliveOnly *bool
	Country   string
	ISP       string
}

type IPInfo struct {
	Ip       string
	Country  string
	ISP      string
	Timezone int
}

type ProxyCheckerResponse struct {
	IP     string `json:"ip"`
	Status string `json:"status"`
	Info   string `json:"info"`
}

type Status string

const (
	SOCKS       Status = "SOCKS"
	HTTP_PROXY  Status = "HTTP_PROXY"
	HTTPS_PROXY Status = "HTTPS_PROXY"
)
