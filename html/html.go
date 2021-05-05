// Package html manages page rendering
package html

import (
	"html/template"
	"io"
	"log"
)

var pageTpl *template.Template

// RobotInfo shows how many times a user agent has been seen.
type RobotInfo struct {
	Seen      int
	UserAgent string
}

// Page shows how many times robots have been seen and how many tried to look at our secrets!
type Page struct {
	Title    string
	Robots   []RobotInfo
	Cheaters []RobotInfo
}

func init() {
	tpl := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>{{.Title}}</title>
	<link rel="stylesheet" href="https://gist.githubusercontent.com/fabjan/322b61203b3e0fa36c862d331f360793/raw/c3853ef186ac0ef48070a686a0cad7007ad11341/better.css">
</head>
<body>

<h2>Robots</h2>
{{if not .Robots}}<p><em>(no data)</em></p>
{{else}}
<table>
    <tr>
        <th>Seen</th>
        <th>User-Agent</th>
    </tr>
    {{ range .Robots}}
        <tr>
            <td>{{ .Seen }}</td>
            <td>{{ .UserAgent }}</td>
        </tr>
	{{ end}}
</table>
{{end}}

<h2>Cheaters</h2>
{{if not .Cheaters}}<p><em>(no data)</em></p>
{{else}}
<table>
    <tr>
        <th>Seen</th>
        <th>User-Agent</th>
    </tr>
    {{ range .Cheaters}}
        <tr>
            <td>{{ .Seen }}</td>
            <td>{{ .UserAgent }}</td>
        </tr>
	{{ end}}
</table>
{{end}}

</html>
`
	var err error
	pageTpl, err = template.New("page").Parse(tpl)
	if err != nil {
		log.Fatal(err)
	}
}

// Render renders a HTML page to the given writer.
func Render(w io.Writer, p Page) error {
	return pageTpl.Execute(w, p)
}
