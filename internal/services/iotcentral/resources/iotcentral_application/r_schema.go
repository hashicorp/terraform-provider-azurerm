// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iotcentral_application

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iotcentral/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func iotcentralApplicationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the IotHub resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ApplicationName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sub_domain": {
			Description:  "A `sub_domain` name. Subdomain for the IoT Central URL. Each application must have a unique subdomain.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ApplicationSubdomain,
		},

		"display_name": {
			Description:  "A `display_name` name. Custom display name for the IoT Central application.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validate.ApplicationDisplayName,
		},

		"identity": commonschema.SystemAssignedIdentityOptional(),

		"public_network_access_enabled": {
			Description: "Whether public network access is allowed for the IoT Central Application.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"sku": {
			Description:  "A `sku` name.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(apps.AppSkuSTOne),
			ValidateFunc: validation.StringInSlice(apps.PossibleValuesForAppSku(), false),
		},
		"template": {
			Description:  "A `template` name. IoT Central application template name.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "iotc-pnp-preview@1.0.0",
			ValidateFunc: validate.ApplicationTemplateName,
		},

		"tags": commonschema.Tags(),
	}
}
