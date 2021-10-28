package managedhsms

type DeletedManagedHsm struct {
	Id         *string                      `json:"id,omitempty"`
	Name       *string                      `json:"name,omitempty"`
	Properties *DeletedManagedHsmProperties `json:"properties,omitempty"`
	Type       *string                      `json:"type,omitempty"`
}
