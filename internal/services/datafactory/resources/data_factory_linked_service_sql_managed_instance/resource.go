// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_sql_managed_instance

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ServicePrincipalKey struct {
	LinkedServiceName string `tfschema:"linked_service_name"`
	SecretName        string `tfschema:"secret_name"`
	SecretVersion     string `tfschema:"secret_version"`
}

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_data_factory_linked_service_sql_managed_instance"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return linkedservices.ValidateLinkedServiceID
}
