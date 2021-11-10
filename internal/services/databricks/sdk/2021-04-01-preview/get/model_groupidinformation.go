package get

type GroupIdInformation struct {
	Id         *string                      `json:"id,omitempty"`
	Name       *string                      `json:"name,omitempty"`
	Properties GroupIdInformationProperties `json:"properties"`
	Type       *string                      `json:"type,omitempty"`
}
