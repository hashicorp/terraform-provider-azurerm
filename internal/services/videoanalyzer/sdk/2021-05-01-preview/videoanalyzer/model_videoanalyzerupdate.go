package videoanalyzer

type VideoAnalyzerUpdate struct {
	Identity   *VideoAnalyzerIdentity         `json:"identity,omitempty"`
	Properties *VideoAnalyzerPropertiesUpdate `json:"properties,omitempty"`
	Tags       *map[string]string             `json:"tags,omitempty"`
}
