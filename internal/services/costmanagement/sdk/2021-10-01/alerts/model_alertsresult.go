package alerts

type AlertsResult struct {
	NextLink *string  `json:"nextLink,omitempty"`
	Value    *[]Alert `json:"value,omitempty"`
}
