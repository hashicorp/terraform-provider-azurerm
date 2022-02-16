package resource

type ServiceCosmosDbConfigurationInfo struct {
	KeyVaultKeyUri  *string `json:"keyVaultKeyUri,omitempty"`
	OfferThroughput *int64  `json:"offerThroughput,omitempty"`
}
