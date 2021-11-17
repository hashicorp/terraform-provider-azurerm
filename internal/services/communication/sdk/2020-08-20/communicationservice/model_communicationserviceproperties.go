package communicationservice

type CommunicationServiceProperties struct {
	DataLocation        string             `json:"dataLocation"`
	HostName            *string            `json:"hostName,omitempty"`
	ImmutableResourceId *string            `json:"immutableResourceId,omitempty"`
	NotificationHubId   *string            `json:"notificationHubId,omitempty"`
	ProvisioningState   *ProvisioningState `json:"provisioningState,omitempty"`
	Version             *string            `json:"version,omitempty"`
}
