package customapis

type WsdlDefinition struct {
	Content      *string           `json:"content,omitempty"`
	ImportMethod *WsdlImportMethod `json:"importMethod,omitempty"`
	Service      *WsdlService      `json:"service,omitempty"`
	Url          *string           `json:"url,omitempty"`
}
