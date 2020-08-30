package tpl

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	templatePath = "templates/web"
	partialPath  = "templates/partial"
)

var (
	pathPrefix string
	tpl        *template.Template
	funcs      = template.FuncMap{"toHTML": toHTML}
)

// Load initializes the `tpl` variable which is the collection of templates.
//
// The path to the templates is based on the current working directory. This will not work as-is for unit tests.
//
// To make it load for unit tests, run the tests with a `-wd` switch.
//
// E.g. go test crgc -run TestTemplate -wd ../../..
func Load(pathPrefix string) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	tpl = template.New("index").Funcs(funcs)
	tpl = template.Must(tpl.ParseGlob(filepath.Join(cwd, pathPrefix, templatePath, "*html")))
	tpl = template.Must(tpl.ParseGlob(filepath.Join(cwd, pathPrefix, partialPath, "*.html")))
}

// Render executes the template and writes to ResponseWriter.
func Render(w http.ResponseWriter, name string, d map[string]interface{}) {
	if e := tpl.ExecuteTemplate(w, name, d); e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
}

// toHTML is a helper function to retain the HTML entities for use with templates.
func toHTML(s string) template.HTML {
	return template.HTML(s)
}
