package configurationstores

type ConfigurationStoreUpdateParameters struct {
	Identity   *ResourceIdentity                             `json:"identity,omitempty"`
	Properties *ConfigurationStorePropertiesUpdateParameters `json:"properties,omitempty"`
	Sku        *Sku                                          `json:"sku,omitempty"`
	Tags       *map[string]string                            `json:"tags,omitempty"`
}
