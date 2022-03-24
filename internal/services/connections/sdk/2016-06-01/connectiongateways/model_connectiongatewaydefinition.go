package connectiongateways

type ConnectionGatewayDefinition struct {
	Etag       *string                                `json:"etag,omitempty"`
	Id         *string                                `json:"id,omitempty"`
	Location   *string                                `json:"location,omitempty"`
	Name       *string                                `json:"name,omitempty"`
	Properties *ConnectionGatewayDefinitionProperties `json:"properties,omitempty"`
	Tags       *map[string]string                     `json:"tags,omitempty"`
	Type       *string                                `json:"type,omitempty"`
}
