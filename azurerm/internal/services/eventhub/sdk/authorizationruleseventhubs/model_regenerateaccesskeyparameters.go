package authorizationruleseventhubs

type RegenerateAccessKeyParameters struct {
	Key     *string `json:"key,omitempty"`
	KeyType KeyType `json:"keyType"`
}
