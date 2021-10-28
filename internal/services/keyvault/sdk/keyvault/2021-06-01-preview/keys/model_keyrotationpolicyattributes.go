package keys

type KeyRotationPolicyAttributes struct {
	Created    *int64  `json:"created,omitempty"`
	ExpiryTime *string `json:"expiryTime,omitempty"`
	Updated    *int64  `json:"updated,omitempty"`
}
