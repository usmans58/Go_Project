package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/usmans58/bookings/packages/config"
	"github.com/usmans58/bookings/packages/models"
)

var functions = template.FuncMap{}
var app *config.AppConfig

// this function sets the the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}
func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}
func RenderTemplate(w http.ResponseWriter, tmpl string, template_data *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template Cache")

	}

	buf := new(bytes.Buffer) //bytes buffer
	template_data = AddDefaultData(template_data)
	_ = t.Execute(buf, template_data) //execute the template t and store it in buffer
	_, err := buf.WriteTo(w)          //write buffer to the response writer
	if err != nil {
		fmt.Println("Error eriting template the browser", err)
	}
	// parsedTemplate, _ := template.ParseFiles("./HTML_Templates/" + tmpl)
	// err = parsedTemplate.Execute(w, nil)
	// if err != nil {
	// 	fmt.Println("Error Parsing template", err)
	// 	return
	// }
}

//creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	pages, err := filepath.Glob("./HTML_Templates/*.page.html")
	if err != nil {
		return myCache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		// fmt.Println("The current page is :", page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./HTML_Templates/*.layout.html")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./HTML_Templates/*.layout.html")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}

	return myCache, nil
}
