package signalr

type SkuList struct {
	NextLink *string `json:"nextLink,omitempty"`
	Value    *[]Sku  `json:"value,omitempty"`
}
