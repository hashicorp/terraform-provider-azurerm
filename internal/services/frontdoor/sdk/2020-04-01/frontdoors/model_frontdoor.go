package frontdoors

type FrontDoor struct {
	Id         *string              `json:"id,omitempty"`
	Location   *string              `json:"location,omitempty"`
	Name       *string              `json:"name,omitempty"`
	Properties *FrontDoorProperties `json:"properties,omitempty"`
	Tags       *map[string]string   `json:"tags,omitempty"`
	Type       *string              `json:"type,omitempty"`
}
