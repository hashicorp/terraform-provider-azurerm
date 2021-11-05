package accounts

type CreateDataLakeAnalyticsAccountParameters struct {
	Location   string                                   `json:"location"`
	Properties CreateDataLakeAnalyticsAccountProperties `json:"properties"`
	Tags       *map[string]string                       `json:"tags,omitempty"`
}
