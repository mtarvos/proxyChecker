package sqlite

import (
	"math/rand"
	"os"
	"proxyChecker/internal/entity"
	"testing"
)

const testStoragePath = "./storage/test_storage.db"

var testList = []entity.ProxyItem{
	{ip: "1.2.3.5", port: 123, country: "DE", ISP: "Bingo telecom", timezone: 3200, alive: 1, status: 0},
	{ip: "1.2.3.4", port: 123, country: "US", ISP: "Mango telecom", timezone: 3200, alive: 1, status: 0},
	{ip: "1.2.3.6", port: 123, country: "DE", ISP: "Orange telecom", timezone: 3200, alive: 1, status: 0},
	{ip: "1.2.3.7", port: 123, country: "SP", ISP: "Mango telecom", timezone: 3200, alive: 1, status: 0},
	{ip: "1.2.3.8", port: 123, country: "CA", ISP: "Mango telecom", timezone: 3200, alive: 1, status: 0},
	{ip: "1.2.3.9", port: 123, country: "CN", ISP: "Mango telecom", timezone: 3200, alive: 1, status: 1},
	{ip: "1.2.3.10", port: 123, country: "US", ISP: "Orange telecom", timezone: 3200, alive: 1, status: 1},
	{ip: "1.2.3.11", port: 123, country: "US", ISP: "Mango telecom", timezone: 3200, alive: 1, status: 1},
	{ip: "1.2.3.12", port: 123, country: "DE", ISP: "Mango telecom", timezone: 3200, alive: 1, status: 1},
	{ip: "1.2.3.13", port: 123, country: "US", ISP: "Mango telecom", timezone: 3200, alive: 0, status: 1},
}

func TestGetProxyByISP(t *testing.T) {
	var ISPs []string
	for _, item := range testList {
		ISPs = append(ISPs, item.ISP)
	}

	i := rand.Intn(len(ISPs))
	isp := ISPs[i]

	controlList := getListForISP(testList, isp)

	os.Remove(testStoragePath)

	s, err := New(testStoragePath)
	if err != nil {
		t.Error("Failed to init test storage", "error", err.Error())
	}

	if err = s.ClearProxy(); err != nil {
		t.Error("can not clear proxy", "error", err.Error())
	}

	s.SaveProxy(testList)

	list, err := s.GetProxyByISP(isp)
	if err != nil {
		t.Error("can not get proxy by isp", "error", err.Error())
	}

	if !equalProxyList(controlList, list) {
		t.Error("list and controlList are not equal!")
	}
}

func TestGetProxyByCountry(t *testing.T) {
	var countries []string
	for _, item := range testList {
		countries = append(countries, item.country)
	}

	i := rand.Intn(len(countries))
	country := countries[i]

	controlList := getListForCountry(testList, country)

	os.Remove(testStoragePath)

	s, err := New(testStoragePath)
	if err != nil {
		t.Error("Failed to init test storage", "error", err.Error())
	}

	if err = s.ClearProxy(); err != nil {
		t.Error("can not clear proxy", "error", err.Error())
	}

	s.SaveProxy(testList)

	list, err := s.GetProxyByCountry(country)
	if err != nil {
		t.Error("can not get proxy by country", "error", err.Error())
	}

	if !equalProxyList(controlList, list) {
		t.Error("list and controlList are not equal!")
	}
}

func TestGetAliveProxy(t *testing.T) {
	os.Remove(testStoragePath)

	s, err := New(testStoragePath)
	if err != nil {
		t.Error("Failed to init test storage", "error", err.Error())
	}

	if err = s.ClearProxy(); err != nil {
		t.Error("can not clear proxy", "error", err.Error())
	}

	s.SaveProxy(testList)

	porxyList, err := s.GetProxyByAlive(false)
	if err != nil {
		t.Error("Failed to get proxy by alive", "error", err.Error())
	}

	controlList := getAliveList(testList, false)
	if !equalProxyList(porxyList, controlList) {
		t.Error("GetProxyByAlive(false) has wrong result ")
	}

	porxyList, err = s.GetProxyByAlive(true)
	if err != nil {
		t.Error("Failed to get proxy by alive", "error", err.Error())
	}

	controlList = getAliveList(testList, true)
	if !equalProxyList(porxyList, controlList) {
		t.Error("GetProxyByAlive(true) has wrong result ")
	}
}

func TestSaveAndGetProxy(t *testing.T) {
	os.Remove(testStoragePath)

	s, err := New(testStoragePath)
	if err != nil {
		t.Error("Failed to init test storage", "error", err.Error())
	}

	if err = s.ClearProxy(); err != nil {
		t.Error("can not clear proxy", "error", err.Error())
	}

	if err = s.SaveProxy(testList); err != nil {
		t.Error("can not save test proxy", "error", err.Error())
	}

	proxyList, err := s.GetAll()
	if err != nil {
		t.Error("Failed to get proxy list", "error", err.Error())
	}

	if !equalProxyList(testList, proxyList) {
		t.Error("List after save/get is not equal")
	}
}

func getAliveList(proxyList []ProxyItem, alive bool) []ProxyItem {
	if len(proxyList) == 0 {
		return []ProxyItem{}
	}

	aliveInt := 0
	if alive {
		aliveInt = 1
	}

	var resList []ProxyItem

	for _, item := range proxyList {
		if item.alive == aliveInt {
			resList = append(resList, item)
		}
	}

	return resList
}

func getListForISP(list []ProxyItem, isp string) []ProxyItem {
	var resList []ProxyItem

	for _, item := range list {
		if item.ISP == isp {
			resList = append(resList, item)
		}
	}

	return resList
}

func getListForCountry(list []ProxyItem, country string) []ProxyItem {
	var resList []ProxyItem

	for _, item := range list {
		if item.country == country {
			resList = append(resList, item)
		}
	}

	return resList
}

func proxyListContainItem(proxyList []ProxyItem, item ProxyItem) bool {
	if len(proxyList) == 0 {
		return false
	}

	for _, proxyItem := range proxyList {
		if proxyItem.ip == item.ip && proxyItem.port == item.port {
			return true
		}
	}
	return false
}

func equalProxyList(list1 []ProxyItem, list2 []ProxyItem) bool {
	if len(list1) != len(list2) || len(list1) == 0 {
		return false
	}

	for _, item := range list1 {
		if !proxyListContainItem(list2, item) {
			return false
		}
	}

	return true
}
