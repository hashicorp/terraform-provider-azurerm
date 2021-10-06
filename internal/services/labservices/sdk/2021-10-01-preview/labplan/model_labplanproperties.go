package labplan

type LabPlanProperties struct {
	AllowedRegions             *[]string              `json:"allowedRegions,omitempty"`
	DefaultAutoShutdownProfile *AutoShutdownProfile   `json:"defaultAutoShutdownProfile,omitempty"`
	DefaultConnectionProfile   *ConnectionProfile     `json:"defaultConnectionProfile,omitempty"`
	DefaultNetworkProfile      *LabPlanNetworkProfile `json:"defaultNetworkProfile,omitempty"`
	LinkedLmsInstance          *string                `json:"linkedLmsInstance,omitempty"`
	ProvisioningState          *ProvisioningState     `json:"provisioningState,omitempty"`
	SharedGalleryId            *string                `json:"sharedGalleryId,omitempty"`
	SupportInfo                *SupportInfo           `json:"supportInfo,omitempty"`
}
