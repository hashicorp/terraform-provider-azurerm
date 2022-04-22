package redisenterprise

type Module struct {
	Args    *string `json:"args,omitempty"`
	Name    string  `json:"name"`
	Version *string `json:"version,omitempty"`
}
