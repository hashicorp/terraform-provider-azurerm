// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_management_private_link_association

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/resourcemanagementprivatelink"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this Private Link Association, which should be a UUID. If `name` is not provided, a UUID will be generated, you should use [the `ignore_changes` attribute to ignore changes to this field](https://developer.hashicorp.com/terraform/language/block/resource#ignore_changes).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},
		"management_group_id": {
			Description:  "The Management Group ID within which this Private Link Association should exist.",
			ForceNew:     true,
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateManagementGroupID,
		},
		"resource_management_private_link_id": {
			Description:  "The Resource ID of Resource Management Private Link.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: resourcemanagementprivatelink.ValidateResourceManagementPrivateLinkID,
		},
		"public_network_access_enabled": {
			Description: "Whether public network access is allowed.",
			Type:        pluginsdk.TypeBool,
			Required:    true,
			ForceNew:    true,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"tenant_id": {
			Description: "The Tenant ID.",
			Computed:    true,
			Type:        pluginsdk.TypeString,
		},
	}
}
