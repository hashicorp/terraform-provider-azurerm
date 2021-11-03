package videoanalyzer

type StorageAccount struct {
	Id       *string           `json:"id,omitempty"`
	Identity *ResourceIdentity `json:"identity,omitempty"`
	Status   *string           `json:"status,omitempty"`
}
