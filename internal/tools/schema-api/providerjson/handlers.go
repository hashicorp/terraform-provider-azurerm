package providerjson

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	DataSourcesList = "/schema-data/v1/data-sources"  // Lists all data sources in the Provider
	ResourcesList   = "/schema-data/v1/resources"     // Lists all Resources in the Provider
	DataSourcesPath = "/schema-data/v1/data-sources/" // Gets all schema data for a data source
	ResourcesPath   = "/schema-data/v1/resources/"    // Gets all schema data for a Resource
)

func (p ProviderJSON) DataSourcesHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	dsRaw := strings.Split(req.URL.RequestURI(), DataSourcesPath)
	ds := strings.Split(dsRaw[1], "/")[0]
	data, err := resourceFromRaw(p.DataSourcesMap[ds])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(w.Write([]byte(fmt.Sprintf("[{\"error\": \"Could not process schema for %q from provider: %+v\"}]", ds, err))))
	} else if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(w.Write([]byte(fmt.Sprintf("Marshall error: %+v", err))))
	}
}

func (p ProviderJSON) ResourcesHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	dsRaw := strings.Split(req.URL.RequestURI(), ResourcesPath)
	ds := strings.Split(dsRaw[1], "/")[0]
	data, err := resourceFromRaw(p.ResourcesMap[ds])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(w.Write([]byte(fmt.Sprintf("[{\"error\": \"Could not process schema for %q from provider: %+v\"}]", ds, err))))
	} else if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(w.Write([]byte(fmt.Sprintf("Marshall error: %+v", err))))
	}
}

func (p *ProviderJSON) ListResources(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(p.Resources()); err != nil {
		log.Println(w.Write([]byte(fmt.Sprintf("Marshall error: %+v", err))))
	}
}

func (p *ProviderJSON) ListDataSources(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(p.DataSources()); err != nil {
		log.Println(w.Write([]byte(fmt.Sprintf("Marshall error: %+v", err))))
	}
}
