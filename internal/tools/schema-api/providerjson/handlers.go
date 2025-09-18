// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package providerjson

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	DataSourcesList = "/ProviderSchema-data/v1/data-sources"  // Lists all data sources in the Provider
	ResourcesList   = "/ProviderSchema-data/v1/resources"     // Lists all Resources in the Provider
	DataSourcesPath = "/ProviderSchema-data/v1/data-sources/" // Gets all ProviderSchema data for a data source
	ResourcesPath   = "/ProviderSchema-data/v1/resources/"    // Gets all ProviderSchema data for a Resource
	DumpSchema      = "/ProviderSchema-data/v1/dump/"         // Gets all ProviderSchema
)

func (p *ProviderJSON) DataSourcesHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	dsRaw := strings.Split(req.URL.RequestURI(), DataSourcesPath)
	ds := strings.Split(dsRaw[1], "/")[0]
	data, err := resourceFromRaw(p.DataSourcesMap[ds])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(fmt.Fprintf(w, "[{\"error\": \"Could not process ProviderSchema for %q from provider: %+v\"}]", ds, err))
	} else if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Fprintf(w, "Marshall error: %+v", err))
	}
}

func (p *ProviderJSON) ResourcesHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	dsRaw := strings.Split(req.URL.RequestURI(), ResourcesPath)
	ds := strings.Split(dsRaw[1], "/")[0]
	data, err := resourceFromRaw(p.ResourcesMap[ds])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(fmt.Fprintf(w, "[{\"error\": \"Could not process ProviderSchema for %q from provider: %+v\"}]", ds, err))
	} else if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Fprintf(w, "Marshall error: %+v", err))
	}
}

func (p *ProviderJSON) ListResources(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(p.Resources()); err != nil {
		log.Println(fmt.Fprintf(w, "Marshall error: %+v", err))
	}
}

func (p *ProviderJSON) ListDataSources(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(p.DataSources()); err != nil {
		log.Println(fmt.Fprintf(w, "Marshall error: %+v", err))
	}
}

func (p *ProviderJSON) DumpAllSchema(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	provider, err := ProviderFromRaw(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Fprintf(w, "[{\"error\": \"Could not process provider: %+v\"}]", err))
	}
	if err := json.NewEncoder(w).Encode(provider); err != nil {
		log.Println(fmt.Fprintf(w, "Marshall error: %+v", err))
	}
}
