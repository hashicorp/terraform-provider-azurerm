package redisenterprise

type Database struct {
	Id         *string             `json:"id,omitempty"`
	Name       *string             `json:"name,omitempty"`
	Properties *DatabaseProperties `json:"properties,omitempty"`
	Type       *string             `json:"type,omitempty"`
}
