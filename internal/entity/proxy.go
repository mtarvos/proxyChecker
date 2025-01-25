package entity

type ProxyItem struct {
	Ip       string
	Port     int
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
