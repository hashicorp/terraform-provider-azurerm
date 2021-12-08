package application

type ApplicationUpdateParameters struct {
	Tags *map[string]string `json:"tags,omitempty"`
}
