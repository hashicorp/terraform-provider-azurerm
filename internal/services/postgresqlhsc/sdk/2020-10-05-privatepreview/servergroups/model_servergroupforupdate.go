package servergroups

type ServerGroupForUpdate struct {
	Location   *string                         `json:"location,omitempty"`
	Properties *ServerGroupPropertiesForUpdate `json:"properties,omitempty"`
	Tags       *map[string]string              `json:"tags,omitempty"`
}
