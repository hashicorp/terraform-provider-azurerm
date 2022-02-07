package workspaces

type WorkspaceCustomBooleanParameter struct {
	Type  *CustomParameterType `json:"type,omitempty"`
	Value bool                 `json:"value"`
}
