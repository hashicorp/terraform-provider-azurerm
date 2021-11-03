package workspaces

type WorkspaceCustomStringParameter struct {
	Type  *CustomParameterType `json:"type,omitempty"`
	Value string               `json:"value"`
}
