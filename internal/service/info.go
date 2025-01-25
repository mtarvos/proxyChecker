package service

import "proxyChecker/internal/entity"

type GeoData struct {
	IPAddress          string     `json:"ip_address"`
	City               string     `json:"city"`
	CityGeonameID      int        `json:"city_geoname_id"`
	Region             string     `json:"region"`
	RegionISOCode      string     `json:"region_iso_code"`
	RegionGeonameID    int        `json:"region_geoname_id"`
	PostalCode         string     `json:"postal_code"`
	Country            string     `json:"country"`
	CountryCode        string     `json:"country_code"`
	CountryGeonameID   int        `json:"country_geoname_id"`
	CountryIsEU        bool       `json:"country_is_eu"`
	Continent          string     `json:"continent"`
	ContinentCode      string     `json:"continent_code"`
	ContinentGeonameID int        `json:"continent_geoname_id"`
	Longitude          float64    `json:"longitude"`
	Latitude           float64    `json:"latitude"`
	Security           Security   `json:"security"`
	Timezone           Timezone   `json:"timezone"`
	Flag               Flag       `json:"flag"`
	Currency           Currency   `json:"currency"`
	Connection         Connection `json:"connection"`
}

type Security struct {
	IsVPN bool `json:"is_vpn"`
}

type Timezone struct {
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	GMTOffset    int    `json:"gmt_offset"`
	CurrentTime  string `json:"current_time"`
	IsDST        bool   `json:"is_dst"`
}

type Flag struct {
	Emoji   string `json:"emoji"`
	Unicode string `json:"unicode"`
	PNG     string `json:"png"`
	SVG     string `json:"svg"`
}

type Currency struct {
	CurrencyName string `json:"currency_name"`
	CurrencyCode string `json:"currency_code"`
}

type Connection struct {
	AutonomousSystemNumber       int     `json:"autonomous_system_number"`
	AutonomousSystemOrganization string  `json:"autonomous_system_organization"`
	ConnectionType               *string `json:"connection_type"`
	ISPName                      *string `json:"isp_name"`
	OrganizationName             *string `json:"organization_name"`
}

type infoService struct {
	infoProvider InfoProvider
	infoSaver    InfoSaver
}

type InfoProvider interface {
	GetInfo() (entity.IPInfo, error)
}

type InfoSaver interface {
	ProxyUpdate([]entity.ProxyItem) error
}

func newInfoService(infoProvider InfoProvider, infoSaver InfoSaver) *infoService {
	return &infoService{infoProvider: infoProvider, infoSaver: infoSaver}
}
