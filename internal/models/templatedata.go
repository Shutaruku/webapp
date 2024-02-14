package models

import "github.com/YuanData/webapp/internal/forms"

type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float64
	Data            map[string]interface{}
	CSRFToken       string
	Success         string
	Warning         string
	Error           string
	Form            *forms.Form
	IsAuthenticated int
}
