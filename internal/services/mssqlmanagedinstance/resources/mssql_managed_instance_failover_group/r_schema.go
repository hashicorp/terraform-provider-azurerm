// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_failover_group

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/instancefailovergroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlManagedInstanceFailoverGroupName,
		},

		"location": commonschema.Location(),

		"managed_instance_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseSqlManagedInstanceIDInsensitively),
		},

		"partner_managed_instance_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"readonly_endpoint_failover_policy_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"read_write_endpoint_failover_policy": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"mode": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(instancefailovergroups.PossibleValuesForReadWriteEndpointFailoverPolicy(), false),
					},

					"grace_minutes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(60),
					},
				},
			},
		},

		"secondary_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(instancefailovergroups.SecondaryInstanceTypeGeo),
			ValidateFunc: validation.StringInSlice(instancefailovergroups.PossibleValuesForSecondaryInstanceType(), false),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"partner_region": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"location": commonschema.LocationComputed(),

					"role": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"role": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
