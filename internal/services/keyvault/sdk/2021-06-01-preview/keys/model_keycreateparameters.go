package keys

type KeyCreateParameters struct {
	Properties KeyProperties      `json:"properties"`
	Tags       *map[string]string `json:"tags,omitempty"`
}
