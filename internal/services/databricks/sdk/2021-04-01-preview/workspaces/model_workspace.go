package workspaces

type Workspace struct {
	Id         *string             `json:"id,omitempty"`
	Location   string              `json:"location"`
	Name       *string             `json:"name,omitempty"`
	Properties WorkspaceProperties `json:"properties"`
	Sku        *Sku                `json:"sku,omitempty"`
	SystemData *SystemData         `json:"systemData,omitempty"`
	Tags       *map[string]string  `json:"tags,omitempty"`
	Type       *string             `json:"type,omitempty"`
}
