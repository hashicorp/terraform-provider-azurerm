package profiles

type ProfileUpdateParameters struct {
	Tags *map[string]string `json:"tags,omitempty"`
}
