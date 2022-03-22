package customapis

type ApiOAuthSettingsParameter struct {
	Options      *interface{} `json:"options,omitempty"`
	UiDefinition *interface{} `json:"uiDefinition,omitempty"`
	Value        *string      `json:"value,omitempty"`
}
