package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/YuanData/webapp/internal/config"
	"github.com/YuanData/webapp/internal/models"
	"github.com/justinas/nosurf"
)

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Success = app.Session.PopString(r.Context(), "success")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	return td
}

var app *config.AppConfig
var pathToTemplates = "./templates"

func NewRenderer(a *config.AppConfig) {
	app = a
}

func Template(w http.ResponseWriter, r *http.Request, tpml string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tpml]
	if !ok {
		return errors.New("template not in cache for some reason")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	theCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*-page.tpml", pathToTemplates))
	if err != nil {
		return theCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return theCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*-layout.tpml", pathToTemplates))
		if err != nil {
			return theCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*-layout.tpml", pathToTemplates))
			if err != nil {
				return theCache, err
			}
		}

		theCache[name] = ts
	}
	return theCache, nil
}
