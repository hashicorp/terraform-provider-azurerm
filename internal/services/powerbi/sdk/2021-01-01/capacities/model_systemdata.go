package capacities

type SystemData struct {
	CreatedAt          *string       `json:"createdAt,omitempty"`
	CreatedBy          *string       `json:"createdBy,omitempty"`
	CreatedByType      *IdentityType `json:"createdByType,omitempty"`
	LastModifiedAt     *string       `json:"lastModifiedAt,omitempty"`
	LastModifiedBy     *string       `json:"lastModifiedBy,omitempty"`
	LastModifiedByType *IdentityType `json:"lastModifiedByType,omitempty"`
}
