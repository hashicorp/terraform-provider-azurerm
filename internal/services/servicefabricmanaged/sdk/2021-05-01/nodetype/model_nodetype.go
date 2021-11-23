package nodetype

type NodeType struct {
	Id         *string             `json:"id,omitempty"`
	Name       *string             `json:"name,omitempty"`
	Properties *NodeTypeProperties `json:"properties,omitempty"`
	SystemData *SystemData         `json:"systemData,omitempty"`
	Tags       *map[string]string  `json:"tags,omitempty"`
	Type       *string             `json:"type,omitempty"`
}
