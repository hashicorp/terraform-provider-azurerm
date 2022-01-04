package databases

type Persistence struct {
	AofEnabled   *bool         `json:"aofEnabled,omitempty"`
	AofFrequency *AofFrequency `json:"aofFrequency,omitempty"`
	RdbEnabled   *bool         `json:"rdbEnabled,omitempty"`
	RdbFrequency *RdbFrequency `json:"rdbFrequency,omitempty"`
}
