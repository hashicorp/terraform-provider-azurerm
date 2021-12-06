package nodetype

type NodeTypeActionParameters struct {
	Force *bool    `json:"force,omitempty"`
	Nodes []string `json:"nodes"`
}
