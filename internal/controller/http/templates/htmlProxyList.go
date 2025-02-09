package templates

const HTMLProxyList = ` 
<!DOCTYPE html>
<html>
<head>
<style>
body {
font-family: Arial, sans-serif;
padding: 20px;
background-color: #f5f5f5;
}

.container {
max-width: 1200px;
margin: 0 auto;
background-color: white;
border-radius: 8px;
box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
padding: 20px;
}

table {
width: 100%;
border-collapse: collapse;
margin-top: 20px;
}

th, td {
padding: 12px 15px;
text-align: left;
border-bottom: 1px solid #ddd;
}

th {
background-color: #4CAF50;
color: white;
font-weight: bold;
}

tr:nth-child(even) {
background-color: #f9f9f9;
}

tr:hover {
background-color: #f5f5f5;
}

.alive {
color: #4CAF50;
font-weight: bold;
}

.dead {
color: #f44336;
font-weight: bold;
}

.unknown {
  color: #9E9E9E; 
  font-weight: normal; 
}

.pagination {
    display: flex;
    justify-content: center;
    align-items: center;
    margin-top: 20px;
    padding: 10px;
}
.pagination a {
    color: black;
    padding: 8px 16px;
    text-decoration: none;
    border: 1px solid #ddd;
    margin: 0 4px;
    border-radius: 4px;
}
.pagination a:hover {
    background-color: #f5f5f5;
}
.pagination a.active {
    background-color: #4CAF50;
    color: white;
    border: 1px solid #4CAF50;
}
.pagination a.disabled {
    color: #9E9E9E;
    pointer-events: none;
    border: 1px solid #ddd;
}

</style>
</head>
<body>
<div class="container">
<table>
<thead>
<tr>
<th>IP</th>
<th>Port</th>
<th>Out IP</th>
<th>Country</th>
<th>City</th>
<th>ISP</th>
<th>Timezone</th>
<th>Alive</th>
</tr>
</thead>
<tbody>
{{range .ProxyList}}
<tr>
<td>{{.IP}}</td>
<td>{{.Port}}</td>
<td>{{.OutIP.String}}</td>
<td>{{.Country.String}}</td>
<td>{{.City.String}}</td>
<td>{{.ISP.String}}</td>
<td>{{.Timezone.Int32}}</td>
<td class="{{if eq .Alive.Int32 0}}unknown{{else if eq .Alive.Int32 1}}dead{{else if eq .Alive.Int32 2}}alive{{end}}">
  {{if eq .Alive.Int32 0}}Unknown{{else if eq .Alive.Int32 1}}No{{else if eq .Alive.Int32 2}}Yes{{end}}
</td>
</tr>
{{end}}
</tbody>
</table>
<div class="pagination">
        {{if gt .CurrentPage 1}}
            <a href="?page=1&limit={{.Limit}}">&laquo; First</a>
            <a href="?page={{subtract .CurrentPage 1}}&limit={{.Limit}}">&lsaquo; Previous</a>
        {{else}}
            <a class="disabled">&laquo; First</a>
            <a class="disabled">&lsaquo; Previous</a>
        {{end}}
        
        {{range .Pages}}
            {{if eq . $.CurrentPage}}
                <a class="active">{{.}}</a>
            {{else}}
                <a href="?page={{.}}&limit={{$.Limit}}">{{.}}</a>
            {{end}}
        {{end}}
        
        {{if lt .CurrentPage .TotalPages}}
            <a href="?page={{add .CurrentPage 1}}&limit={{.Limit}}">Next &rsaquo;</a>
            <a href="?page={{.TotalPages}}&limit={{.Limit}}">Last &raquo;</a>
        {{else}}
            <a class="disabled">Next &rsaquo;</a>
            <a class="disabled">Last &raquo;</a>
        {{end}}
    </div>
</div>
</body>
</html>
`
