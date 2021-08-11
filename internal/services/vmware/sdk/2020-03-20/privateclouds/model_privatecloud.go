package privateclouds

type PrivateCloud struct {
	Id         *string                `json:"id,omitempty"`
	Location   string                 `json:"location"`
	Name       *string                `json:"name,omitempty"`
	Properties PrivateCloudProperties `json:"properties"`
	Sku        Sku                    `json:"sku"`
	Tags       *map[string]string     `json:"tags,omitempty"`
	Type       *string                `json:"type,omitempty"`
}
