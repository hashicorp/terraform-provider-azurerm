package diskpools

type DiskPoolUpdateProperties struct {
	Disks *[]Disk `json:"disks,omitempty"`
}
