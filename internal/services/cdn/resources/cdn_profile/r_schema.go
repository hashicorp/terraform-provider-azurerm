// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_profile

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cdnProfileSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(cdn.SkuNameStandardAkamai),
				string(cdn.SkuNameStandardChinaCdn),
				string(cdn.SkuNameStandardVerizon),
				string(cdn.SkuNameStandardMicrosoft),
				string(cdn.SkuNamePremiumVerizon),
			}, false),
		},

		"tags": commonschema.Tags(),
	}
}
