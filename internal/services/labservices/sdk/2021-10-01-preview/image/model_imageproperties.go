package image

type ImageProperties struct {
	Author            *string            `json:"author,omitempty"`
	AvailableRegions  *[]string          `json:"availableRegions,omitempty"`
	Description       *string            `json:"description,omitempty"`
	DisplayName       *string            `json:"displayName,omitempty"`
	EnabledState      EnableState        `json:"enabledState"`
	IconUrl           *string            `json:"iconUrl,omitempty"`
	Offer             *string            `json:"offer,omitempty"`
	OsState           *OsState           `json:"osState,omitempty"`
	OsType            *OsType            `json:"osType,omitempty"`
	Plan              *string            `json:"plan,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	Publisher         *string            `json:"publisher,omitempty"`
	SharedGalleryId   *string            `json:"sharedGalleryId,omitempty"`
	Sku               *string            `json:"sku,omitempty"`
	TermsStatus       *EnableState       `json:"termsStatus,omitempty"`
	Version           *string            `json:"version,omitempty"`
}
