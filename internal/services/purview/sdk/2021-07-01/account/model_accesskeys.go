package account

type AccessKeys struct {
	AtlasKafkaPrimaryEndpoint   *string `json:"atlasKafkaPrimaryEndpoint,omitempty"`
	AtlasKafkaSecondaryEndpoint *string `json:"atlasKafkaSecondaryEndpoint,omitempty"`
}
