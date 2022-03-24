package managedapis

type ApiResourceMetadata struct {
	ApiType          *ApiType           `json:"apiType,omitempty"`
	BrandColor       *string            `json:"brandColor,omitempty"`
	ConnectionType   *string            `json:"connectionType,omitempty"`
	HideKey          *string            `json:"hideKey,omitempty"`
	Source           *string            `json:"source,omitempty"`
	Tags             *map[string]string `json:"tags,omitempty"`
	WsdlImportMethod *WsdlImportMethod  `json:"wsdlImportMethod,omitempty"`
	WsdlService      *WsdlService       `json:"wsdlService,omitempty"`
}
