package autoscalevcores

type AutoScaleVCore struct {
	Id         *string                   `json:"id,omitempty"`
	Location   string                    `json:"location"`
	Name       *string                   `json:"name,omitempty"`
	Properties *AutoScaleVCoreProperties `json:"properties,omitempty"`
	Sku        AutoScaleVCoreSku         `json:"sku"`
	SystemData *SystemData               `json:"systemData,omitempty"`
	Tags       *map[string]string        `json:"tags,omitempty"`
	Type       *string                   `json:"type,omitempty"`
}
