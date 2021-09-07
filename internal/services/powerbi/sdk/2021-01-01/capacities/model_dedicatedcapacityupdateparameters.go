package capacities

type DedicatedCapacityUpdateParameters struct {
	Properties *DedicatedCapacityMutableProperties `json:"properties,omitempty"`
	Sku        *CapacitySku                        `json:"sku,omitempty"`
	Tags       *map[string]string                  `json:"tags,omitempty"`
}
