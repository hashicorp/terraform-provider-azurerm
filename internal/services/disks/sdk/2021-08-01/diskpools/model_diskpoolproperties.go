package diskpools

type DiskPoolProperties struct {
	AdditionalCapabilities *[]string          `json:"additionalCapabilities,omitempty"`
	AvailabilityZones      []string           `json:"availabilityZones"`
	Disks                  *[]Disk            `json:"disks,omitempty"`
	ProvisioningState      ProvisioningStates `json:"provisioningState"`
	Status                 OperationalStatus  `json:"status"`
	SubnetId               string             `json:"subnetId"`
}
