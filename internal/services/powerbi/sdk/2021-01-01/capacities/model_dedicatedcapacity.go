package capacities

type DedicatedCapacity struct {
	Id         *string                      `json:"id,omitempty"`
	Location   string                       `json:"location"`
	Name       *string                      `json:"name,omitempty"`
	Properties *DedicatedCapacityProperties `json:"properties,omitempty"`
	Sku        CapacitySku                  `json:"sku"`
	SystemData *SystemData                  `json:"systemData,omitempty"`
	Tags       *map[string]string           `json:"tags,omitempty"`
	Type       *string                      `json:"type,omitempty"`
}
