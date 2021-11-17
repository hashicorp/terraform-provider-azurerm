package topictypes

type TopicTypeProperties struct {
	Description              *string                     `json:"description,omitempty"`
	DisplayName              *string                     `json:"displayName,omitempty"`
	Provider                 *string                     `json:"provider,omitempty"`
	ProvisioningState        *TopicTypeProvisioningState `json:"provisioningState,omitempty"`
	ResourceRegionType       *ResourceRegionType         `json:"resourceRegionType,omitempty"`
	SourceResourceFormat     *string                     `json:"sourceResourceFormat,omitempty"`
	SupportedLocations       *[]string                   `json:"supportedLocations,omitempty"`
	SupportedScopesForSource *[]SupportedScopesForSource `json:"supportedScopesForSource,omitempty"`
}
