package diskpools

type DiskPoolCreateProperties struct {
	AdditionalCapabilities *[]string `json:"additionalCapabilities,omitempty"`
	AvailabilityZones      *[]string `json:"availabilityZones,omitempty"`
	Disks                  *[]Disk   `json:"disks,omitempty"`
	SubnetId               string    `json:"subnetId"`
}
