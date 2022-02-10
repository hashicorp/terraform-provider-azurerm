package geographichierarchies

type Region struct {
	Code    *string   `json:"code,omitempty"`
	Name    *string   `json:"name,omitempty"`
	Regions *[]Region `json:"regions,omitempty"`
}
