package entity

type StatsData struct {
	Total        int
	Alive        int
	Dead         int
	UniqCountry  int
	UniqISP      int
	CountryStats []CountryStatsItem
	ISPStats     []ISPStatsItem
}

type CountryStatsItem struct {
	Country string
	Count   int
}

type ISPStatsItem struct {
	ISP   string
	Count int
}
