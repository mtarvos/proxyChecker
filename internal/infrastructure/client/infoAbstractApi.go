package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"proxyChecker/internal/entity"
	"proxyChecker/internal/lib/helpers"
)

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

type AbstractAPI struct {
	log        *slog.Logger
	getInfoURL string
	getInfoKey string
}

func NewAbstractAPI(log *slog.Logger, getInfoURL string, getInfoKey string) *AbstractAPI {
	return &AbstractAPI{log: log, getInfoURL: getInfoURL, getInfoKey: getInfoKey}
}

func (a *AbstractAPI) getFullURL(url string, key string, ip string) (string, error) {
	const fn = "AbstractAPI.getFullURL"
	if url == "" || key == "" || ip == "" {
		return "", fmt.Errorf("%s: url or key or ip are not defined", fn)
	}

	return fmt.Sprintf("%s?api_key=%s&ip_address=%s", url, key, ip), nil
}

func (a *AbstractAPI) GetInfo(ctx context.Context, ip string) (entity.IPInfo, error) {
	const fn = "AbstractAPI.GetInfo"

	fullURL, err := a.getFullURL(a.getInfoURL, a.getInfoKey, ip)
	if err != nil {
		return entity.IPInfo{}, fmt.Errorf("%s: Error get full url for AbstractAPI: %w", fn, err)
	}

	_, result, err := helpers.SendGetRequest(ctx, fullURL)
	if err != nil {
		return entity.IPInfo{}, fmt.Errorf("%s: Error getting ip info by AbstractAPI: %w", fn, err)
	}

	var geoData GeoData
	err = json.Unmarshal([]byte(result), &geoData)
	if err != nil {
		return entity.IPInfo{}, fmt.Errorf("%s: Error unmarshalling result to struct: %w, json: %s", fn, err, result)
	}

	return entity.IPInfo{
		IP:       ip,
		Country:  geoData.Country,
		City:     geoData.City,
		ISP:      a.getPointerValue(geoData.Connection.ISPName),
		Timezone: geoData.Timezone.GMTOffset,
	}, nil
}

func (a *AbstractAPI) getPointerValue(fieldVal *string) string {
	if fieldVal != nil {
		return *fieldVal
	}
	return ""
}
