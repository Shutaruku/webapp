package handlers

import (
	"net/http"

	"github.com/YuanData/webapp/pkg/config"
	"github.com/YuanData/webapp/pkg/models"
	"github.com/YuanData/webapp/pkg/render"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "home-page.tpml", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	sidekickMap := make(map[string]string)
	sidekickMap["morty"] = "Ooh, wee!"

	render.RenderTemplate(w, "about-page.tpml", &models.TemplateData{
		StringMap: sidekickMap,
	})
}
