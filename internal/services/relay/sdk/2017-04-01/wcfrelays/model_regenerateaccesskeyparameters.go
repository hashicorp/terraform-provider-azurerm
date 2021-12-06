package wcfrelays

type RegenerateAccessKeyParameters struct {
	Key     *string `json:"key,omitempty"`
	KeyType KeyType `json:"keyType"`
}
