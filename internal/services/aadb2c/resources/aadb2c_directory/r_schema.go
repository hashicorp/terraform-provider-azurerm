// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package aadb2c_directory

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/aadb2c/2021-04-01-preview/tenants"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"domain_name": {
			Description:  "Domain name of the B2C tenant, including onmicrosoft.com suffix.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"country_code": {
			Description:  "Country code of the B2C tenant. See https://aka.ms/B2CDataResidency for valid country codes.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"data_residency_location": {
			Description:  "Location in which the B2C tenant is hosted and data resides. See https://aka.ms/B2CDataResidency for more information.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(tenants.PossibleValuesForLocation(), false),
		},

		"display_name": {
			Description:  "The initial display name of the B2C tenant.",
			Type:         pluginsdk.TypeString,
			Optional:     true, // azignore:AZS007 - pre-existing violation
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"sku_name": {
			Description: "Billing SKU for the B2C tenant. See https://aka.ms/b2cBilling for more information.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(tenants.SkuNamePremiumP1),
				string(tenants.SkuNamePremiumP2),
				// string(tenants.SkuNameStandard), // API doesn't seem to support "Standard", it's ignored and "PremiumP1" is used instead, even when patching
			}, false),
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"billing_type": {
			Description: "The type of billing for the B2C tenant. Possible values include: `MAU` or `Auths`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"effective_start_date": {
			Description: "The date from which the billing type took effect. May not be populated until after the first billing cycle.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tenant_id": {
			Description: "The Tenant ID for the B2C tenant.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
