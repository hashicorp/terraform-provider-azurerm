// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_cache_access_policy_assignment

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/rediscacheaccesspolicyassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Redis Cache Access Policy Assignment.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"redis_cache_id": {
			Description:  "The ID of the Redis Cache.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: rediscacheaccesspolicyassignments.ValidateRediID,
		},
		"access_policy_name": {
			Description: "The name of the Access Policy to be assigned.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"object_id": {
			Description: "The principal ID to be assigned the Access Policy.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"object_id_alias": {
			Description:  "The alias of the principal ID. User-friendly name for object ID. Also represents username for token based authentication.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
