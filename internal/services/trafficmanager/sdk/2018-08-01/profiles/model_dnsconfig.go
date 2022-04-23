package profiles

type DnsConfig struct {
	Fqdn         *string `json:"fqdn,omitempty"`
	RelativeName *string `json:"relativeName,omitempty"`
	Ttl          *int64  `json:"ttl,omitempty"`
}
