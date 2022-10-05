package provider

// NOTE: this file is generated - manual changes will be overwritten.

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadtestservice"
)

func autoRegisteredTypedServices() []sdk.TypedServiceRegistration {
	return []sdk.TypedServiceRegistration{
		loadtestservice.Registration{},
	}
}
