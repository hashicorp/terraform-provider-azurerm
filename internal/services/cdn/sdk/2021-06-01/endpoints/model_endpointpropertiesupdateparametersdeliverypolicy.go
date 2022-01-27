package endpoints

type EndpointPropertiesUpdateParametersDeliveryPolicy struct {
	Description *string        `json:"description,omitempty"`
	Rules       []DeliveryRule `json:"rules"`
}
