// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_catalog

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/catalogs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Required:     true,
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"dev_center_id": {
			Required:     true,
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			ValidateFunc: catalogs.ValidateDevCenterID,
		},

		"catalog_github": CatalogPropertiesSchema(),

		"catalog_adogit": CatalogPropertiesSchema(),
	}
}

func (r Resource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}
