package azurerm

// requiredResourceProviders returns all of the Resource Providers used by the AzureRM Provider
// whilst all may not be used by every user - the intention is that we determine which should be
// registered such that we can avoid obscure errors where Resource Providers aren't registered.
// new Resource Providers should be added to this list as they're used in the Provider
// NOTE: Resource Providers in this list are case sensitive
func requiredResourceProviders() map[string]struct{} {
	return map[string]struct{}{
		"Microsoft.ApiManagement":       {},
		"Microsoft.Authorization":       {},
		"Microsoft.Automation":          {},
		"Microsoft.Cache":               {},
		"Microsoft.Cdn":                 {},
		"Microsoft.CognitiveServices":   {},
		"Microsoft.Compute":             {},
		"Microsoft.ContainerInstance":   {},
		"Microsoft.ContainerRegistry":   {},
		"Microsoft.ContainerService":    {},
		"Microsoft.Databricks":          {},
		"Microsoft.DataLakeStore":       {},
		"Microsoft.DBforMySQL":          {},
		"Microsoft.DBforPostgreSQL":     {},
		"Microsoft.Devices":             {},
		"Microsoft.DevTestLab":          {},
		"Microsoft.DocumentDB":          {},
		"Microsoft.EventGrid":           {},
		"Microsoft.EventHub":            {},
		"Microsoft.KeyVault":            {},
		"microsoft.insights":            {},
		"Microsoft.Logic":               {},
		"Microsoft.ManagedIdentity":     {},
		"Microsoft.Management":          {},
		"Microsoft.Network":             {},
		"Microsoft.NotificationHubs":    {},
		"Microsoft.OperationalInsights": {},
		"Microsoft.Relay":               {},
		"Microsoft.Resources":           {},
		"Microsoft.Search":              {},
		"Microsoft.ServiceBus":          {},
		"Microsoft.ServiceFabric":       {},
		"Microsoft.Sql":                 {},
		"Microsoft.Storage":             {},
	}
}
