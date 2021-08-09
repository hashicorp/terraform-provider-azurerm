package signalr

type SignalRFeature struct {
	Flag       FeatureFlags       `json:"flag"`
	Properties *map[string]string `json:"properties,omitempty"`
	Value      string             `json:"value"`
}
