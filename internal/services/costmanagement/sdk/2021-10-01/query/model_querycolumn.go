package query

type QueryColumn struct {
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
}
