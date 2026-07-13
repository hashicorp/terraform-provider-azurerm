// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

// Package server exposes the provider schema over HTTP. It is the API-server
// half of the schema-api tool and depends only on the shared providerschema
// model, keeping the model package free of HTTP concerns.
package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/providerschema"
)

const (
	DataSourcesList = "/ProviderSchema-data/v1/data-sources"  // Lists all data sources in the Provider
	ResourcesList   = "/ProviderSchema-data/v1/resources"     // Lists all Resources in the Provider
	DataSourcesPath = "/ProviderSchema-data/v1/data-sources/" // Gets all ProviderSchema data for a data source
	ResourcesPath   = "/ProviderSchema-data/v1/resources/"    // Gets all ProviderSchema data for a Resource
	DumpSchema      = "/ProviderSchema-data/v1/dump/"         // Gets all ProviderSchema
)

// Server serves the provider schema over HTTP.
type Server struct {
	Provider *providerschema.ProviderJSON
}

// New returns a Server backed by the given provider schema.
func New(p *providerschema.ProviderJSON) *Server {
	return &Server{Provider: p}
}

func (s *Server) DataSourcesHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	dsRaw := strings.Split(req.URL.RequestURI(), DataSourcesPath)
	ds := strings.Split(dsRaw[1], "/")[0]
	data, err := providerschema.ResourceFromRaw(s.Provider.DataSourcesMap[ds])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(fmt.Fprintf(w, "[{\"error\": \"Could not process ProviderSchema for %q from provider: %+v\"}]", ds, err))
	} else if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Fprintf(w, "Marshall error: %+v", err))
	}
}

func (s *Server) ResourcesHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	dsRaw := strings.Split(req.URL.RequestURI(), ResourcesPath)
	ds := strings.Split(dsRaw[1], "/")[0]
	data, err := providerschema.ResourceFromRaw(s.Provider.ResourcesMap[ds])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(fmt.Fprintf(w, "[{\"error\": \"Could not process ProviderSchema for %q from provider: %+v\"}]", ds, err))
	} else if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Fprintf(w, "Marshall error: %+v", err))
	}
}

func (s *Server) ListResources(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(s.Provider.Resources()); err != nil {
		log.Println(fmt.Fprintf(w, "Marshall error: %+v", err))
	}
}

func (s *Server) ListDataSources(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(s.Provider.DataSources()); err != nil {
		log.Println(fmt.Fprintf(w, "Marshall error: %+v", err))
	}
}

func (s *Server) DumpAllSchema(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	provider, err := providerschema.ProviderFromRaw(s.Provider)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Fprintf(w, "[{\"error\": \"Could not process provider: %+v\"}]", err))
	}
	if err := json.NewEncoder(w).Encode(provider); err != nil {
		log.Println(fmt.Fprintf(w, "Marshall error: %+v", err))
	}
}
