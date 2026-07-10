package helpers

type DocType string

const DocTypeDataSource DocType = "data-source"
const DocTypeResource DocType = "resource"

type Schema struct {
	Type        string      `json:"type,omitempty"`
	ConfigMode  string      `json:"config_mode,omitempty"`
	Optional    bool        `json:"optional,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Description string      `json:"description,omitempty"`
	Computed    bool        `json:"computed,omitempty"`
	ForceNew    bool        `json:"force_new,omitempty"`
	Elem        interface{} `json:"elem,omitempty"`
	MaxItems    int         `json:"max_items,omitempty"`
	MinItems    int         `json:"min_items,omitempty"`
}

type Resource struct {
	Schema   map[string]Schema `json:"schema"`
	Timeouts *Timeouts         `json:"timeouts,omitempty"`
}

type Timeouts struct {
	Create int `json:"create,omitempty"`
	Read   int `json:"read,omitempty"`
	Delete int `json:"delete,omitempty"`
	Update int `json:"update,omitempty"`
}
