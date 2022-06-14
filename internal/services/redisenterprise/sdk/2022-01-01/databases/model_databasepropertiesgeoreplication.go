package databases

type DatabasePropertiesGeoReplication struct {
	GroupNickname   *string           `json:"groupNickname,omitempty"`
	LinkedDatabases *[]LinkedDatabase `json:"linkedDatabases,omitempty"`
}
