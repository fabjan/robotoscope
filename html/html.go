// Package html manages page rendering
package html

import (
	_ "embed" // for index.html.tpl
	"html/template"
	"io"
	"log"

	"github.com/fabjan/robotoscope/core"
)

var pageTpl *template.Template

// Page shows how many times robots have been seen and how many tried to look at our secrets!
type Page struct {
	Title    string
	Robots   []core.RobotInfo
	Cheaters []core.RobotInfo
}

//go:embed index.html.tpl
var tpl string

func init() {
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
