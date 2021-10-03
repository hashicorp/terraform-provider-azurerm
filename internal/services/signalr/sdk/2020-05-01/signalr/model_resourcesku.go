package signalr

type ResourceSku struct {
	Capacity *int64          `json:"capacity,omitempty"`
	Family   *string         `json:"family,omitempty"`
	Name     string          `json:"name"`
	Size     *string         `json:"size,omitempty"`
	Tier     *SignalRSkuTier `json:"tier,omitempty"`
}
