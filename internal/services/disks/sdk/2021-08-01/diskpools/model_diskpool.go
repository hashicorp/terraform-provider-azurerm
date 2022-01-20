package diskpools

type DiskPool struct {
	Id                *string            `json:"id,omitempty"`
	Location          string             `json:"location"`
	ManagedBy         *string            `json:"managedBy,omitempty"`
	ManagedByExtended *[]string          `json:"managedByExtended,omitempty"`
	Name              *string            `json:"name,omitempty"`
	Properties        DiskPoolProperties `json:"properties"`
	Sku               *Sku               `json:"sku,omitempty"`
	SystemData        *SystemMetadata    `json:"systemData,omitempty"`
	Tags              *map[string]string `json:"tags,omitempty"`
	Type              *string            `json:"type,omitempty"`
}
