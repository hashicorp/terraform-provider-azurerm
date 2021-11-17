package eventsubscriptions

type EventSubscriptionIdentity struct {
	Type                 *EventSubscriptionIdentityType `json:"type,omitempty"`
	UserAssignedIdentity *string                        `json:"userAssignedIdentity,omitempty"`
}
