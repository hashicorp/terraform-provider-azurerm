// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SubnetServiceEndpointName() pluginsdk.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"Microsoft.AzureActiveDirectory",
		"Microsoft.AzureCosmosDB",
		"Microsoft.CognitiveServices",
		"Microsoft.ContainerRegistry",
		"Microsoft.EventHub",
		"Microsoft.KeyVault",
		"Microsoft.ServiceBus",
		"Microsoft.Sql",
		"Microsoft.Storage",
		"Microsoft.Storage.Global",
		"Microsoft.Web",
	}, false)
}
