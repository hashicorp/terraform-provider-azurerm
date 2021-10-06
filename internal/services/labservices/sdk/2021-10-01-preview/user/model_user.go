package user

type User struct {
	Id         *string        `json:"id,omitempty"`
	Name       *string        `json:"name,omitempty"`
	Properties UserProperties `json:"properties"`
	SystemData *SystemData    `json:"systemData,omitempty"`
	Type       *string        `json:"type,omitempty"`
}
