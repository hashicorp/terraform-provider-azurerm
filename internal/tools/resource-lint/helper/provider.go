package helper

import (
	"sync"

	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var (
	typedServices []sdk.TypedServiceRegistration
	servicesOnce  sync.Once
)

// GetTypedServices returns the typed services, initializing once
func GetTypedServices() []sdk.TypedServiceRegistration {
	servicesOnce.Do(func() {
		provider.AzureProvider() // initialize provider
		typedServices = provider.SupportedTypedServices()
	})
	return typedServices
}
