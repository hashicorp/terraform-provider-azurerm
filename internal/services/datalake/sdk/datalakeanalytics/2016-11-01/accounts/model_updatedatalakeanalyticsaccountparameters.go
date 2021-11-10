package accounts

type UpdateDataLakeAnalyticsAccountParameters struct {
	Properties *UpdateDataLakeAnalyticsAccountProperties `json:"properties,omitempty"`
	Tags       *map[string]string                        `json:"tags,omitempty"`
}
