package autoscalevcores

type AutoScaleVCoreSku struct {
	Capacity *int64        `json:"capacity,omitempty"`
	Name     string        `json:"name"`
	Tier     *VCoreSkuTier `json:"tier,omitempty"`
}
