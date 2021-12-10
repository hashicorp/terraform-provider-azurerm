package accounts

type SystemData struct {
	CreatedAt          *string        `json:"createdAt,omitempty"`
	CreatedBy          *string        `json:"createdBy,omitempty"`
	CreatedByType      *CreatedByType `json:"createdByType,omitempty"`
	LastModifiedAt     *string        `json:"lastModifiedAt,omitempty"`
	LastModifiedBy     *string        `json:"lastModifiedBy,omitempty"`
	LastModifiedByType *CreatedByType `json:"lastModifiedByType,omitempty"`
}
