package namespaces

type RegenerateAccessKeyParameters struct {
	Key     *string `json:"key,omitempty"`
	KeyType KeyType `json:"keyType"`
}
