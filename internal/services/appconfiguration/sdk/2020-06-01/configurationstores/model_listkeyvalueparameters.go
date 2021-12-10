package configurationstores

type ListKeyValueParameters struct {
	Key   string  `json:"key"`
	Label *string `json:"label,omitempty"`
}
