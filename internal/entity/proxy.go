package entity

type ProxyItem struct {
	ID       int64  `db:"id"`
	IP       string `db:"proxy"`
	Port     int    `db:"port"`
	OutIP    string `db:"out_ip"`
	Country  string `db:"country"`
	City     string `db:"city"`
	ISP      string `db:"ISP"`
	Timezone int    `db:"timezone"`
	Alive    int    `db:"alive"`
}

const (
	Eq = "equal"
	Ne = "not_equal"
)

type Operand string

type Filters struct {
	AliveOnly *bool
	Country   *StringFilter
	City      *StringFilter
	ISP       *StringFilter
	OutIP     *StringFilter
}

type StringFilter struct {
	Val interface{}
	Op  Operand
}

type IPInfo struct {
	IP       string
	Country  string
	City     string
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
