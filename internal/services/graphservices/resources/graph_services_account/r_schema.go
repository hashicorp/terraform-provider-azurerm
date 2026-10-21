// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package graph_services_account

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			ForceNew:     true,
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"application_id": {
			ForceNew:     true,
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.IsUUID,
		},
		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"billing_plan_id": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
	}
}
