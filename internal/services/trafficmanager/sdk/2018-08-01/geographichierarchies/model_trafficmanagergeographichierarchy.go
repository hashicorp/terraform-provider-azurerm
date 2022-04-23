package geographichierarchies

type TrafficManagerGeographicHierarchy struct {
	Id         *string                        `json:"id,omitempty"`
	Name       *string                        `json:"name,omitempty"`
	Properties *GeographicHierarchyProperties `json:"properties,omitempty"`
	Type       *string                        `json:"type,omitempty"`
}
