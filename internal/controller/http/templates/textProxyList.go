package templates

const TEXTProxyList = "{{range .}}{{.IP}};{{.Port}};{{.OutIP.String}};{{.Country.String}};{{.City.String}};{{.ISP.String}};{{.Timezone.Int32}}\n{{end}}"
