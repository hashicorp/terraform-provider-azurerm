package autoscalevcores

type AutoScaleVCoreUpdateParameters struct {
	Properties *AutoScaleVCoreMutableProperties `json:"properties,omitempty"`
	Sku        *AutoScaleVCoreSku               `json:"sku,omitempty"`
	Tags       *map[string]string               `json:"tags,omitempty"`
}
