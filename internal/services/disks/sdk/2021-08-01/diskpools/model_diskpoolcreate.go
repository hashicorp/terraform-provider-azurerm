package diskpools

type DiskPoolCreate struct {
	Id                *string                  `json:"id,omitempty"`
	Location          string                   `json:"location"`
	ManagedBy         *string                  `json:"managedBy,omitempty"`
	ManagedByExtended *[]string                `json:"managedByExtended,omitempty"`
	Name              *string                  `json:"name,omitempty"`
	Properties        DiskPoolCreateProperties `json:"properties"`
	Sku               Sku                      `json:"sku"`
	Tags              *map[string]string       `json:"tags,omitempty"`
	Type              *string                  `json:"type,omitempty"`
}
