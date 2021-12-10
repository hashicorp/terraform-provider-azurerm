package workspaces

type WorkspaceCustomObjectParameter struct {
	Type  *CustomParameterType `json:"type,omitempty"`
	Value interface{}          `json:"value"`
}
