package cognitiveservicesaccounts

type CheckSkuAvailabilityParameter struct {
	Kind string   `json:"kind"`
	Skus []string `json:"skus"`
	Type string   `json:"type"`
}
