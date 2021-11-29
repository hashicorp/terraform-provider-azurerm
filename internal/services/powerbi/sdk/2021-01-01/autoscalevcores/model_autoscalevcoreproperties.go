package autoscalevcores

type AutoScaleVCoreProperties struct {
	CapacityLimit     *int64                  `json:"capacityLimit,omitempty"`
	CapacityObjectId  *string                 `json:"capacityObjectId,omitempty"`
	ProvisioningState *VCoreProvisioningState `json:"provisioningState,omitempty"`
}
