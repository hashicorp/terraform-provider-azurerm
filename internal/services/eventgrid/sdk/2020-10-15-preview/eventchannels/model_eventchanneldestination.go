package eventchannels

type EventChannelDestination struct {
	AzureSubscriptionId *string `json:"azureSubscriptionId,omitempty"`
	PartnerTopicName    *string `json:"partnerTopicName,omitempty"`
	ResourceGroup       *string `json:"resourceGroup,omitempty"`
}
