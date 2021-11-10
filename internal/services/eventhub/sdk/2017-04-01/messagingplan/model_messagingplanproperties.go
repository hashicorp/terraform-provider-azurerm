package messagingplan

type MessagingPlanProperties struct {
	Revision             *int64  `json:"revision,omitempty"`
	SelectedEventHubUnit *int64  `json:"selectedEventHubUnit,omitempty"`
	Sku                  *int64  `json:"sku,omitempty"`
	UpdatedAt            *string `json:"updatedAt,omitempty"`
}
