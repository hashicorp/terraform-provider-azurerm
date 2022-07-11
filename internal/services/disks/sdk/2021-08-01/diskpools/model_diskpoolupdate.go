package diskpools

type DiskPoolUpdate struct {
	ManagedBy         *string                  `json:"managedBy,omitempty"`
	ManagedByExtended *[]string                `json:"managedByExtended,omitempty"`
	Properties        DiskPoolUpdateProperties `json:"properties"`
	Sku               *Sku                     `json:"sku,omitempty"`
	Tags              *map[string]string       `json:"tags,omitempty"`
}
