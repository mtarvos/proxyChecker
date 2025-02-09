package entity

type StringField *string

type ProxyItem struct {
	ID       int64            `db:"id" json:"-"`
	IP       string           `db:"proxy" json:"ip"`
	Port     int              `db:"port" json:"port"`
	OutIP    CustomNullString `db:"out_ip" json:"outIP"`
	Country  CustomNullString `db:"country" json:"country"`
	City     CustomNullString `db:"city" json:"city"`
	ISP      CustomNullString `db:"ISP" json:"ISP"`
	Timezone CustomNullInt32  `db:"timezone" json:"timezone"`
	Alive    CustomNullInt32  `db:"alive" json:"alive"`
}

const (
	Eq = "equal"
	Ne = "not_equal"
)

type Operand string

type Filters struct {
	Alive   *int
	Country *StringFilter
	City    *StringFilter
	ISP     *StringFilter
	OutIP   *StringFilter
	Page    int
	Limit   int
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
