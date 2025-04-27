// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourceproviders

import (
	"fmt"
	"strings"
)

// This file contains sets of resource providers which the provider should automatically register, depending on
// configuration by the user. Historically, we have ordained the same set of RPs for all users, and users could
// enable or disable automatic registration of the RPs in that set.
//
// Users may now choose a set of RPs that best meets their needs, via the provider block, or they can provide their
// own set. They can also choose to not register any RPs, and manage these out-of-band, if they wish.
//
// Note that resource providers are case-sensitive.
//
// Official Docs: https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/azure-services-resource-providers

type ResourceProviders map[string]struct{}

const (
	ProviderRegistrationsNone     = "none"
	ProviderRegistrationsLegacy   = "legacy"
	ProviderRegistrationsCore     = "core"
	ProviderRegistrationsExtended = "extended"
	ProviderRegistrationsAll      = "all"
)

func (r ResourceProviders) Add(providers ...string) {
	for _, p := range providers {
		r[p] = struct{}{}
	}
}

func (r ResourceProviders) Merge(a ResourceProviders) ResourceProviders {
	for k, v := range a {
		r[k] = v
	}
	return r
}

// Core is a minimal set of RPs, and will be the default set in a future major version of the provider.
func Core() ResourceProviders {
	return ResourceProviders{
		"Microsoft.Authorization":       {},
		"Microsoft.Compute":             {},
		"Microsoft.CostManagement":      {},
		"Microsoft.ManagedIdentity":     {},
		"Microsoft.MarketplaceOrdering": {},
		"Microsoft.Network":             {},
		"Microsoft.Resources":           {},
		"Microsoft.Storage":             {},
	}
}

// Extended is an expanded set of resource providers, as suggested by the community.
func Extended() ResourceProviders {
	return Core().Merge(ResourceProviders{
		"Microsoft.ApiManagement":        {},
		"Microsoft.AppConfiguration":     {},
		"Microsoft.AppPlatform":          {},
		"Microsoft.Automation":           {},
		"Microsoft.Cache":                {},
		"Microsoft.Cdn":                  {},
		"Microsoft.ContainerInstance":    {},
		"Microsoft.ContainerRegistry":    {},
		"Microsoft.ContainerService":     {},
		"Microsoft.DBforMySQL":           {},
		"Microsoft.DBforPostgreSQL":      {},
		"Microsoft.DataFactory":          {},
		"Microsoft.DataLakeAnalytics":    {},
		"Microsoft.DataLakeStore":        {},
		"Microsoft.DataMigration":        {},
		"Microsoft.DataProtection":       {},
		"Microsoft.Databricks":           {},
		"Microsoft.DevTestLab":           {},
		"Microsoft.Devices":              {},
		"Microsoft.DocumentDB":           {},
		"Microsoft.EventGrid":            {},
		"Microsoft.EventHub":             {},
		"Microsoft.HDInsight":            {},
		"microsoft.insights":             {},
		"Microsoft.KeyVault":             {},
		"Microsoft.Kusto":                {},
		"Microsoft.Logic":                {},
		"Microsoft.Maintenance":          {},
		"Microsoft.Management":           {},
		"Microsoft.NotificationHubs":     {},
		"Microsoft.OperationalInsights":  {},
		"Microsoft.OperationsManagement": {},
		"Microsoft.PowerBIDedicated":     {},
		"Microsoft.Relay":                {},
		"Microsoft.Security":             {},
		"Microsoft.SecurityInsights":     {},
		"Microsoft.ServiceBus":           {},
		"Microsoft.ServiceFabric":        {},
		"Microsoft.SignalRService":       {},
		"Microsoft.Sql":                  {},
		"Microsoft.StreamAnalytics":      {},
		"Microsoft.Web":                  {},
	})
}

// All is intended to be a complete set of RPs that might be needed to utilize any functionality in the provider.
func All() ResourceProviders {
	return Extended().Merge(ResourceProviders{
		"Github.Network":                    {},
		"Microsoft.AVS":                     {},
		"Microsoft.AlertsManagement":        {},
		"Microsoft.Blueprint":               {},
		"Microsoft.BotService":              {},
		"Microsoft.CognitiveServices":       {},
		"Microsoft.CustomProviders":         {},
		"Microsoft.Dashboard":               {},
		"Microsoft.DesktopVirtualization":   {},
		"Microsoft.GuestConfiguration":      {},
		"Microsoft.HealthcareApis":          {},
		"Microsoft.IoTCentral":              {},
		"Microsoft.MachineLearningServices": {},
		"Microsoft.ManagedServices":         {},
		"Microsoft.Maps":                    {},
		"Microsoft.MixedReality":            {},
		"Microsoft.Monitor":                 {},
		"Microsoft.PolicyInsights":          {},
		"Microsoft.RecoveryServices":        {},
		"Microsoft.Search":                  {},
	})
}

// Legacy is the set of automatically registered RPs from earlier versions of the provider, and is provided for
// forwards compatibility. This set should not be changed going forward, and will be removed in a future major release.
func Legacy() ResourceProviders {
	return ResourceProviders{
		"Microsoft.AVS":                     {},
		"Microsoft.ApiManagement":           {},
		"Microsoft.AppConfiguration":        {},
		"Microsoft.AppPlatform":             {},
		"Microsoft.Authorization":           {},
		"Microsoft.Automation":              {},
		"Microsoft.Blueprint":               {},
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
		"Microsoft.DBforMariaDB":            {},
		"Microsoft.DBforMySQL":              {},
		"Microsoft.DBforPostgreSQL":         {},
		"Microsoft.DataFactory":             {},
		"Microsoft.DataLakeAnalytics":       {},
		"Microsoft.DataLakeStore":           {},
		"Microsoft.DataMigration":           {},
		"Microsoft.DataProtection":          {},
		"Microsoft.Databricks":              {},
		"Microsoft.DesktopVirtualization":   {},
		"Microsoft.DevTestLab":              {},
		"Microsoft.Devices":                 {},
		"Microsoft.DocumentDB":              {},
		"Microsoft.EventGrid":               {},
		"Microsoft.EventHub":                {},
		"Microsoft.GuestConfiguration":      {},
		"Microsoft.HDInsight":               {},
		"Microsoft.HealthcareApis":          {},
		"Microsoft.KeyVault":                {},
		"Microsoft.Kusto":                   {},
		"Microsoft.Logic":                   {},
		"Microsoft.MachineLearningServices": {},
		"Microsoft.Maintenance":             {},
		"Microsoft.ManagedIdentity":         {},
		"Microsoft.ManagedServices":         {},
		"Microsoft.Management":              {},
		"Microsoft.Maps":                    {},
		"Microsoft.MarketplaceOrdering":     {},
		"Microsoft.MixedReality":            {},
		"Microsoft.Network":                 {},
		"Microsoft.NotificationHubs":        {},
		"Microsoft.OperationalInsights":     {},
		"Microsoft.OperationsManagement":    {},
		"Microsoft.PolicyInsights":          {},
		"Microsoft.PowerBIDedicated":        {},
		"Microsoft.RecoveryServices":        {},
		"Microsoft.Relay":                   {},
		"Microsoft.Resources":               {},
		"Microsoft.Search":                  {},
		"Microsoft.Security":                {},
		"Microsoft.SecurityInsights":        {},
		"Microsoft.ServiceBus":              {},
		"Microsoft.ServiceFabric":           {},
		"Microsoft.SignalRService":          {},
		"Microsoft.Sql":                     {},
		"Microsoft.Storage":                 {},
		"Microsoft.StreamAnalytics":         {},
		"Microsoft.Web":                     {},
		"microsoft.insights":                {},
	}
}

func GetResourceProvidersSet(input string) (ResourceProviders, error) {
	empty := make(ResourceProviders)
	switch strings.ToLower(input) {
	case ProviderRegistrationsLegacy:
		return Legacy(), nil
	case ProviderRegistrationsCore:
		return Core(), nil
	case ProviderRegistrationsAll:
		return All(), nil
	case ProviderRegistrationsExtended:
		return Extended(), nil
	case ProviderRegistrationsNone:
		return empty, nil
	}

	return empty, fmt.Errorf("unsupported value %q for provider property `resource_provider_registrations`", input)
}
