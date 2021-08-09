package managedidentity

type SystemAssignedIdentity struct {
	Id         *string                           `json:"id,omitempty"`
	Location   string                            `json:"location"`
	Name       *string                           `json:"name,omitempty"`
	Properties *SystemAssignedIdentityProperties `json:"properties,omitempty"`
	Tags       *map[string]string                `json:"tags,omitempty"`
	Type       *string                           `json:"type,omitempty"`
}
