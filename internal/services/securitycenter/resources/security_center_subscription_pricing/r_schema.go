// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_subscription_pricing

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" //nolint:staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func securityCenterSubscriptionPricingSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"tier": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInEnumSlice(security.PossiblePricingTierValues(), false),
		},

		"resource_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "VirtualMachines",
			ValidateFunc: validation.StringInSlice([]string{
				"AI",
				"Api",
				"AppServices",
				"ContainerRegistry",
				"KeyVaults",
				"KubernetesService",
				"SqlServers",
				"SqlServerVirtualMachines",
				"StorageAccounts",
				"VirtualMachines",
				"Arm",
				"Dns",
				"OpenSourceRelationalDatabases",
				"Containers",
				"CosmosDbs",
				"CloudPosture",
			}, false),
		},

		"subplan": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"extension": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotWhiteSpace,
					},
					"additional_extension_properties": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
					},
				},
			},
		},
	}
}
