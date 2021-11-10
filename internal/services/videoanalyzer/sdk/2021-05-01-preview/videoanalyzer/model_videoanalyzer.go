package videoanalyzer

type VideoAnalyzer struct {
	Id         *string                        `json:"id,omitempty"`
	Identity   *VideoAnalyzerIdentity         `json:"identity,omitempty"`
	Location   string                         `json:"location"`
	Name       *string                        `json:"name,omitempty"`
	Properties *VideoAnalyzerPropertiesUpdate `json:"properties,omitempty"`
	SystemData *SystemData                    `json:"systemData,omitempty"`
	Tags       *map[string]string             `json:"tags,omitempty"`
	Type       *string                        `json:"type,omitempty"`
}
