package main

import (
	"fmt"
	"testing"

	"github.com/YuanData/webapp/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
	default:
		t.Error(fmt.Sprintf("Type mismatch: Expected *chi.Mux, got %T", v))
	}
}
