package locations

type Quota struct {
	HostsRemaining *map[string]int64 `json:"hostsRemaining,omitempty"`
	QuotaEnabled   *QuotaEnabled     `json:"quotaEnabled,omitempty"`
}
