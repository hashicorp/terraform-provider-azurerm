// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_failover_group

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/failovergroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Failover Group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlFailoverGroupName,
		},

		"server_id": {
			Description:  "The ID of the primary SQL Server on which to create the failover group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseSqlServerIDInsensitively),
		},

		"partner_server": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Description:  "The ID of a partner SQL server to include in the failover group.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},

					"location": commonschema.LocationComputed(),

					"role": {
						Description: "The replication role of the partner server.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"databases": {
			Description: "A set of database names to include in the failover group.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Set:         pluginsdk.HashString,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"readonly_endpoint_failover_policy_enabled": {
			Description: "Whether failover is enabled for the readonly endpoint.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
		},

		"read_write_endpoint_failover_policy": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"mode": {
						Description:  "The failover policy of the read-write endpoint for the failover group.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(failovergroups.PossibleValuesForReadWriteEndpointFailoverPolicy(), false),
					},
					"grace_minutes": {
						Description:  "The grace period in minutes, before failover with data loss is attempted for the read-write endpoint.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(60),
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return nil
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model MsSqlFailoverGroupModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("DecodeDiff: %+v", err)
			}

			if rwPolicy := model.ReadWriteEndpointFailurePolicy; len(rwPolicy) > 0 {
				if rwPolicy[0].Mode == string(failovergroups.ReadWriteEndpointFailoverPolicyAutomatic) && rwPolicy[0].GraceMinutes < 60 {
					return fmt.Errorf("`grace_minutes` should be %d or greater when `mode` is %q", 60, failovergroups.ReadWriteEndpointFailoverPolicyAutomatic)
				}
				if rwPolicy[0].Mode == string(failovergroups.ReadWriteEndpointFailoverPolicyManual) && rwPolicy[0].GraceMinutes > 0 {
					return fmt.Errorf("`grace_minutes` should not be specified when `mode` is %q", failovergroups.ReadWriteEndpointFailoverPolicyManual)
				}
			}

			return nil
		},
	}
}
