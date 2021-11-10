package videoanalyzer

type VideoStreamingToken struct {
	ExpirationDate *string `json:"expirationDate,omitempty"`
	Token          *string `json:"token,omitempty"`
}
