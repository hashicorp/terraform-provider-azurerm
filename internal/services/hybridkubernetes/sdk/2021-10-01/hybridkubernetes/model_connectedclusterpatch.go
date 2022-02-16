package hybridkubernetes

type ConnectedClusterPatch struct {
	Properties *interface{}       `json:"properties,omitempty"`
	Tags       *map[string]string `json:"tags,omitempty"`
}
