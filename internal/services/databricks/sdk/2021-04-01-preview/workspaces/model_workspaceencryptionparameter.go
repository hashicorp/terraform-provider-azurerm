package workspaces

type WorkspaceEncryptionParameter struct {
	Type  *CustomParameterType `json:"type,omitempty"`
	Value *Encryption          `json:"value,omitempty"`
}
