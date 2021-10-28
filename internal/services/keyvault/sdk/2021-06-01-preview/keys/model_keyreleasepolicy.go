package keys

type KeyReleasePolicy struct {
	ContentType *string `json:"contentType,omitempty"`
	Data        *string `json:"data,omitempty"`
}
