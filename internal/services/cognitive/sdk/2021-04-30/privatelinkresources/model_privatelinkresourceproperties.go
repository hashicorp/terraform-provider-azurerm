package privatelinkresources

type PrivateLinkResourceProperties struct {
	DisplayName       *string   `json:"displayName,omitempty"`
	GroupId           *string   `json:"groupId,omitempty"`
	RequiredMembers   *[]string `json:"requiredMembers,omitempty"`
	RequiredZoneNames *[]string `json:"requiredZoneNames,omitempty"`
}
