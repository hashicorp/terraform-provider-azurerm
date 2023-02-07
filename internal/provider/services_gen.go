package provider

// NOTE: this file is generated - manual changes will be overwritten.

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadtestservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedidentity"
)

func autoRegisteredTypedServices() []sdk.TypedServiceRegistration {
	return []sdk.TypedServiceRegistration{
		containers.Registration{},
		loadtestservice.Registration{},
		managedidentity.Registration{},
	}
}
