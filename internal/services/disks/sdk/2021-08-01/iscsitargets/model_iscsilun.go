package iscsitargets

type IscsiLun struct {
	Lun                        *int64 `json:"lun,omitempty"`
	ManagedDiskAzureResourceId string `json:"managedDiskAzureResourceId"`
	Name                       string `json:"name"`
}
