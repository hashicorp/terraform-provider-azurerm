package edgenodes

type CidrIpAddress struct {
	BaseIpAddress *string `json:"baseIpAddress,omitempty"`
	PrefixLength  *int64  `json:"prefixLength,omitempty"`
}
