package get

type GroupIdInformationProperties struct {
	GroupId           *string   `json:"groupId,omitempty"`
	RequiredMembers   *[]string `json:"requiredMembers,omitempty"`
	RequiredZoneNames *[]string `json:"requiredZoneNames,omitempty"`
}
