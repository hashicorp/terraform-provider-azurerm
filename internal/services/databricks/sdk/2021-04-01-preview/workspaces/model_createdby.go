package workspaces

type CreatedBy struct {
	ApplicationId *string `json:"applicationId,omitempty"`
	Oid           *string `json:"oid,omitempty"`
	Puid          *string `json:"puid,omitempty"`
}
