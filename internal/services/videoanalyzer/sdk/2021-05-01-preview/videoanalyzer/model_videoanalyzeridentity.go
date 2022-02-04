package videoanalyzer

type VideoAnalyzerIdentity struct {
	Type                   string                                  `json:"type"`
	UserAssignedIdentities *map[string]UserAssignedManagedIdentity `json:"userAssignedIdentities,omitempty"`
}
