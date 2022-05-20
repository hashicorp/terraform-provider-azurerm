package vmcollectionupdate

type VMCollectionUpdate struct {
	OperationName *OperationName `json:"operationName,omitempty"`
	VmResourceId  *string        `json:"vmResourceId,omitempty"`
}
