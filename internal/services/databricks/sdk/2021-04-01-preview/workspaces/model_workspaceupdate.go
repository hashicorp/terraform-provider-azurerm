package workspaces

type WorkspaceUpdate struct {
	Tags *map[string]string `json:"tags,omitempty"`
}
