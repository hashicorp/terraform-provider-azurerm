package routes

type ActivatedResourceReference struct {
	Id       *string `json:"id,omitempty"`
	IsActive *bool   `json:"isActive,omitempty"`
}
