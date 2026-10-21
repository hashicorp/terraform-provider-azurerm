// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subnet_service_endpoint_storage_policy

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	mgValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func subnetServiceEndpointStoragePolicySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SubnetServiceEndpointStoragePolicyName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"definition": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 2,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.SubnetServiceEndpointStoragePolicyDefinitionName,
					},

					"service_resources": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.Any(
								azure.ValidateResourceID,
								mgValidate.ManagementGroupID,
								validation.StringInSlice([]string{
									"/services/Azure",
									"/services/Azure/Batch",
									"/services/Azure/Databricks",
									"/services/Azure/DataFactory",
									"/services/Azure/MachineLearning",
									"/services/Azure/ManagedInstance",
									"/services/Azure/WebPI",
								}, false),
							),
						},
					},

					"description": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringLenBetween(0, 140),
					},

					"service": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "Microsoft.Storage",
						ValidateFunc: validation.StringInSlice([]string{
							"Microsoft.Storage",
							"Global",
						}, false),
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
