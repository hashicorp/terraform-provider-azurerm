package managedcluster

type SystemData struct {
	CreatedAt          *string `json:"createdAt,omitempty"`
	CreatedBy          *string `json:"createdBy,omitempty"`
	CreatedByType      *string `json:"createdByType,omitempty"`
	LastModifiedAt     *string `json:"lastModifiedAt,omitempty"`
	LastModifiedBy     *string `json:"lastModifiedBy,omitempty"`
	LastModifiedByType *string `json:"lastModifiedByType,omitempty"`
}
