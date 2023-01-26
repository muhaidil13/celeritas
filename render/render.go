package render

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/CloudyKit/jet/v6"
)

type Render struct {
	Renderer   string
	RootPath   string
	Secure     bool
	ServerName string
	Port       string
	JetViews   *jet.Set
}

type TemplateData struct {
	IsAuthenticated bool
	intMap          map[string]int
	FloatMap        map[string]int
	Data            map[string]interface{}
	StringMap       map[string]string
	CSRFToken       string
	Port            string
	ServerName      string
	Secure          bool
}

func (c *Render) Page(w http.ResponseWriter, r *http.Request, view string, variable, data interface{}) error {
	switch strings.ToLower(c.Renderer) {
	case "go":
		return c.GoPage(w, r, view, data)
	case "jet":
		return c.JetPage(w, r, view, variable, data)
	}

	return nil
}

func (c *Render) GoPage(w http.ResponseWriter, r *http.Request, view string, data interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.page.tmpl", c.RootPath, view))
	if err != nil {
		return err
	}
	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}
	err = tmpl.Execute(w, &td)
	if err != nil {
		return err
	}
	return nil
}

// render jet Template
func (c *Render) JetPage(w http.ResponseWriter, r *http.Request, templateName string, variable, data interface{}) error {
	// type for passing data into template
	var vars jet.VarMap
	if variable == nil {
		vars = make(jet.VarMap)
	} else {
		vars = variable.(jet.VarMap)
	}

	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}
	t, err := c.JetViews.GetTemplate(fmt.Sprintf("%s.jet", templateName))
	if err != nil {
		log.Println(err)
		return err
	}
	if err = t.Execute(w, vars, td); err != nil {
		log.Println(err)

		return err
	}
	return nil
}
