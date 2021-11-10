package managedclusterversion

type ManagedClusterCodeVersionResult struct {
	Id         *string                       `json:"id,omitempty"`
	Name       *string                       `json:"name,omitempty"`
	Properties *ManagedClusterVersionDetails `json:"properties,omitempty"`
	Type       *string                       `json:"type,omitempty"`
}
