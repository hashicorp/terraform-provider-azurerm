package services

type ServiceUpdateParameters struct {
	Tags *map[string]string `json:"tags,omitempty"`
}
