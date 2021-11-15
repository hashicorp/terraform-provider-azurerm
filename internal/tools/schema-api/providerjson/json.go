package providerjson

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
)

const (
	DataSourcesPath = "/schema-data/v1/data-sources/"
	ResourcesPath   = "/schema-data/v1/resources/"
	DataSourcesList = "/schema-data/v1/data-sources"
	ResourcesList   = "/schema-data/v1/resources"
)

type ProviderData struct {
	*schema.Provider `json:"provider"`
}

type Provider schema.Provider
type Schema schema.Schema
type Resource schema.Resource

func (p *Provider) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Schema        map[string]*schema.Schema   `json:"schema"`
		ResourceMap   map[string]*schema.Resource `json:"resource_map,omitempty"`
		DataSourceMap map[string]*schema.Resource `json:"data_source_map,omitempty"`
	}{
		Schema:        p.Schema,
		ResourceMap:   p.ResourcesMap,
		DataSourceMap: p.DataSourcesMap,
	})
}

func (s *Schema) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type          schema.ValueType        `json:"type"`
		ConfigMode    schema.SchemaConfigMode `json:"config_mode,omitempty"`
		Optional      bool                    `json:"optional,omitempty"`
		Required      bool                    `json:"required,omitempty"`
		Default       interface{}             `json:"default,omitempty"`
		Description   string                  `json:"description,omitempty"`
		InputDefault  string                  `json:"-"`
		Computed      bool                    `json:"computed,omitempty"`
		ForceNew      bool                    `json:"force_new,omitempty"`
		Elem          interface{}             `json:"elem,omitempty"`
		MaxItems      int                     `json:"max_items,omitempty"`
		MinItems      int                     `json:"min_items,omitempty"`
		Set           interface{}             `json:"-"`
		ComputedWhen  []string                `json:"computed_when,omitempty"`
		ConflictsWith []string                `json:"conflicts_with,omitempty"`
		ExactlyOneOf  []string                `json:"exactly_one_of,omitempty"`
		AtLeastOneOf  []string                `json:"at_least_one_of,omitempty"`
		RequiredWith  []string                `json:"required_with,omitempty"`
		Deprecated    string                  `json:"deprecated,omitempty"`
		Sensitive     bool                    `json:"sensitive,omitempty"`
	}{
		Type:          s.Type,
		ConfigMode:    s.ConfigMode,
		Optional:      s.Optional,
		Required:      s.Required,
		Default:       s.Default,
		Description:   s.Description,
		Computed:      s.Computed,
		ForceNew:      s.ForceNew,
		Elem:          s.Elem,
		MaxItems:      s.MaxItems,
		MinItems:      s.MinItems,
		ComputedWhen:  s.ComputedWhen,
		ConflictsWith: s.ConflictsWith,
		ExactlyOneOf:  s.ExactlyOneOf,
		AtLeastOneOf:  s.AtLeastOneOf,
		RequiredWith:  s.RequiredWith,
		Deprecated:    s.Deprecated,
		Sensitive:     s.Sensitive,
	})
}

func (r *Resource) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Schema        map[string]*schema.Schema `json:"schema,omitempty"`
		SchemaVersion int                       `json:"schema_version,omitempty"`
	}{
		Schema:        r.Schema,
		SchemaVersion: r.SchemaVersion,
	})
}

type JsonSchema struct {
	Type          interface{} `json:"type"`
	ConfigMode    interface{} `json:"config_mode,omitempty"`
	Optional      bool        `json:"optional,omitempty"`
	Required      bool        `json:"required,omitempty"`
	Default       interface{} `json:"default,omitempty"`
	Description   string      `json:"description,omitempty"`
	InputDefault  string      `json:"-"`
	Computed      bool        `json:"computed,omitempty"`
	ForceNew      bool        `json:"force_new,omitempty"`
	Elem          interface{} `json:"elem,omitempty"`
	MaxItems      int         `json:"max_items,omitempty"`
	MinItems      int         `json:"min_items,omitempty"`
	Set           interface{} `json:"-"`
	ComputedWhen  []string    `json:"computed_when,omitempty"`
	ConflictsWith []string    `json:"conflicts_with,omitempty"`
	ExactlyOneOf  []string    `json:"exactly_one_of,omitempty"`
	AtLeastOneOf  []string    `json:"at_least_one_of,omitempty"`
	RequiredWith  []string    `json:"required_with,omitempty"`
	Deprecated    string      `json:"deprecated,omitempty"`
	Sensitive     bool        `json:"sensitive,omitempty"`
}

type ResourceData struct {
	Schema        map[string]JsonSchema `json:"schema,omitempty"`
	SchemaVersion int                   `json:"schema_version,omitempty"`
	//MigrateState         interface{}        `json:"-"`
	//StateUpgraders       interface{}        `json:"-"`
	//Create               interface{}        `json:"-"`
	//Read                 interface{}        `json:"-"`
	//Update               interface{}        `json:"-"`
	//Delete               interface{}        `json:"-"`
	//Exists               interface{}        `json:"-"`
	//CreateContext        interface{}        `json:"-"`
	//ReadContext          interface{}        `json:"-"`
	//UpdateContext        interface{}        `json:"-"`
	//DeleteContext        interface{}        `json:"-"`
	//CreateWithoutTimeout interface{}        `json:"-"`
	//ReadWithoutTimeout   interface{}        `json:"-"`
	//UpdateWithoutTimeout interface{}        `json:"-"`
	//DeleteWithoutTimeout interface{}        `json:"-"`
	// CustomizeDiff        interface{}        `json:"-"`
}

func ResourceCopy(input *schema.Resource) (r ResourceData) {
	if input == nil {
		return r
	}
	r.Schema = schemaCopy(input.Schema)
	r.SchemaVersion = input.SchemaVersion

	return r
}

func schemaCopy(input map[string]*schema.Schema) map[string]JsonSchema {
	s := make(map[string]JsonSchema, 0)
	for k, p := range input {
		v := *p
		s[k] = JsonSchema{
			Type:          v.Type,
			ConfigMode:    v.ConfigMode,
			Optional:      v.Optional,
			Required:      v.Required,
			Default:       v.Default,
			Description:   v.Description,
			InputDefault:  v.InputDefault,
			Computed:      v.Computed,
			ForceNew:      v.ForceNew,
			Elem:          v.Elem,
			MaxItems:      v.MaxItems,
			MinItems:      v.MinItems,
			Set:           v.Set,
			ComputedWhen:  v.ComputedWhen,
			ConflictsWith: v.ConflictsWith,
			ExactlyOneOf:  v.ExactlyOneOf,
			AtLeastOneOf:  v.AtLeastOneOf,
			RequiredWith:  v.RequiredWith,
			Deprecated:    v.Deprecated,
			Sensitive:     v.Sensitive,
		}
	}

	return s
}

func (p *ProviderData) LoadData() {
	p.Provider = provider.AzureProvider()
}

func (p *ProviderData) DataSourcesHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if path := req.URL; path != nil {
		if strings.EqualFold(path.RequestURI(), DataSourcesList) {
			if err := json.NewEncoder(w).Encode(p.Provider.DataSources()); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Failed to process provider data: %+v", err)))
			}
		} else {
			dsRaw := strings.Split(path.RequestURI(), DataSourcesPath)
			ds := strings.Split(dsRaw[1], "/")[0]
			if len(ds) > 0 {
				dsResource := ResourceCopy(p.Provider.DataSourcesMap[ds])
				if err := json.NewEncoder(w).Encode(dsResource); err != nil {
					w.Write([]byte(fmt.Sprintf("Failed to process provider data for data source %q: %+v", ds, err)))
				}
			}

		}
	}
}

func (p *ProviderData) ListResources(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(p.Provider.Resources()); err != nil {
		panic(err)
	}
}

func (p *Provider) ShowResource(w http.ResponseWriter, _ *http.Request) {

}
