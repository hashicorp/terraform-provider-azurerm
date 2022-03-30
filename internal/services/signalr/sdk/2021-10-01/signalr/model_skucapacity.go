package signalr

type SkuCapacity struct {
	AllowedValues *[]int64   `json:"allowedValues,omitempty"`
	Default       *int64     `json:"default,omitempty"`
	Maximum       *int64     `json:"maximum,omitempty"`
	Minimum       *int64     `json:"minimum,omitempty"`
	ScaleType     *ScaleType `json:"scaleType,omitempty"`
}
