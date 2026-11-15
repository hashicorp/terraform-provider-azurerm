// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_credential_service_principal

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/credentials"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var (
	_ sdk.Resource           = Resource{}
	_ sdk.ResourceWithUpdate = Resource{}
)

func (Resource) ResourceType() string {
	return "azurerm_data_factory_credential_service_principal"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return credentials.ValidateCredentialID
}
