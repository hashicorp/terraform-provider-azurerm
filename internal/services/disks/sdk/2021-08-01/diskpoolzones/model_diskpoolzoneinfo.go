package diskpoolzones

type DiskPoolZoneInfo struct {
	AdditionalCapabilities *[]string `json:"additionalCapabilities,omitempty"`
	AvailabilityZones      *[]string `json:"availabilityZones,omitempty"`
	Sku                    *Sku      `json:"sku,omitempty"`
}
