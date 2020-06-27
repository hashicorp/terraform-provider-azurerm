package provider

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/resources/mgmt/resources"
	"github.com/hashicorp/go-azure-helpers/resourceproviders"
)

// RequiredResourceProviders returns all of the Resource Providers used by the AzureRM Provider
// whilst all may not be used by every user - the intention is that we determine which should be
// registered such that we can avoid obscure errors where Resource Providers aren't registered.
// new Resource Providers should be added to this list as they're used in the Provider
// (this is the approach used by Microsoft in their tooling)
func RequiredResourceProviders() map[string]struct{} {
	// NOTE: Resource Providers in this list are case sensitive
	return map[string]struct{}{
		"Microsoft.ApiManagement":           {},
		"Microsoft.AppPlatform":             {},
		"Microsoft.Authorization":           {},
		"Microsoft.Automation":              {},
		"Microsoft.Blueprints":              {},
		"Microsoft.BotService":              {},
		"Microsoft.Cache":                   {},
		"Microsoft.Cdn":                     {},
		"Microsoft.CognitiveServices":       {},
		"Microsoft.Compute":                 {},
		"Microsoft.ContainerInstance":       {},
		"Microsoft.ContainerRegistry":       {},
		"Microsoft.ContainerService":        {},
		"Microsoft.CostManagement":          {},
		"Microsoft.CustomProviders":         {},
		"Microsoft.Databricks":              {},
		"Microsoft.DataLakeAnalytics":       {},
		"Microsoft.DataLakeStore":           {},
		"Microsoft.DataMigration":           {},
		"Microsoft.DBforMariaDB":            {},
		"Microsoft.DBforMySQL":              {},
		"Microsoft.DBforPostgreSQL":         {},
		"Microsoft.Devices":                 {},
		"Microsoft.DevSpaces":               {},
		"Microsoft.DevTestLab":              {},
		"Microsoft.DocumentDB":              {},
		"Microsoft.EventGrid":               {},
		"Microsoft.EventHub":                {},
		"Microsoft.HDInsight":               {},
		"Microsoft.Healthcare":              {},
		"Microsoft.KeyVault":                {},
		"Microsoft.Kusto":                   {},
		"microsoft.insights":                {},
		"Microsoft.Logic":                   {},
		"Microsoft.MachineLearningServices": {},
		"Microsoft.Maintenance":             {},
		"Microsoft.ManagedIdentity":         {},
		"Microsoft.Management":              {},
		"Microsoft.Maps":                    {},
		"Microsoft.MarketplaceOrdering":     {},
		"Microsoft.Media":                   {},
		"Microsoft.MixedReality":            {},
		"Microsoft.Network":                 {},
		"Microsoft.NotificationHubs":        {},
		"Microsoft.OperationalInsights":     {},
		"Microsoft.OperationsManagement":    {},
		"Microsoft.PowerBIDedicated":        {},
		"Microsoft.Relay":                   {},
		"Microsoft.RecoveryServices":        {},
		"Microsoft.Resources":               {},
		"Microsoft.Search":                  {},
		"Microsoft.Security":                {},
		"Microsoft.SecurityInsights":        {},
		"Microsoft.ServiceBus":              {},
		"Microsoft.ServiceFabric":           {},
		"Microsoft.Sql":                     {},
		"Microsoft.Storage":                 {},
		"Microsoft.StorageCache":            {},
		"Microsoft.StreamAnalytics":         {},
		"Microsoft.TimeSeriesInsights":      {},
		"Microsoft.Web":                     {},
	}
}

func EnsureResourceProvidersAreRegistered(ctx context.Context, client resources.ProvidersClient, availableRPs []resources.Provider, requiredRPs map[string]struct{}) error {
	log.Printf("[DEBUG] Determining which Resource Providers require Registration")
	providersToRegister := resourceproviders.DetermineResourceProvidersRequiringRegistration(availableRPs, requiredRPs)

	if len(providersToRegister) > 0 {
		log.Printf("[DEBUG] Registering %d Resource Providers", len(providersToRegister))
		if err := resourceproviders.RegisterForSubscription(ctx, client, providersToRegister); err != nil {
			return err
		}
	} else {
		log.Printf("[DEBUG] All required Resource Providers are registered")
	}

	return nil
}
