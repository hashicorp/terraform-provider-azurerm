// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/failovergroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MsSqlFailoverGroupDataSourceModel struct {
	Databases                            []string                                        `tfschema:"databases"`
	Name                                 string                                          `tfschema:"name"`
	PartnerServers                       []PartnerServerDataSourceModel                  `tfschema:"partner_server"`
	ReadonlyEndpointFailurePolicyEnabled bool                                            `tfschema:"readonly_endpoint_failover_policy_enabled"`
	ServerId                             string                                          `tfschema:"server_id"`
	Tags                                 map[string]string                               `tfschema:"tags"`
	ReadWriteEndpointFailurePolicy       []ReadWriteEndpointFailurePolicyDataSourceModel `tfschema:"read_write_endpoint_failover_policy"`
}

type PartnerServerDataSourceModel struct {
	ID       string `tfschema:"id"`
	Location string `tfschema:"location"`
	Role     string `tfschema:"role"`
}

type ReadWriteEndpointFailurePolicyDataSourceModel struct {
	GraceMinutes int64  `tfschema:"grace_minutes"`
	Mode         string `tfschema:"mode"`
}

var _ sdk.DataSource = MsSqlFailoverGroupDataSource{}

type MsSqlFailoverGroupDataSource struct{}

func (d MsSqlFailoverGroupDataSource) ResourceType() string {
	return "azurerm_mssql_failover_group"
}

func (d MsSqlFailoverGroupDataSource) ModelObject() interface{} {
	return nil
}

func (d MsSqlFailoverGroupDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ValidateMsSqlFailoverGroupName,
		},

		"server_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ServerID,
		},
	}
}

func (d MsSqlFailoverGroupDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"partner_server": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"location": commonschema.LocationComputed(),

					"role": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"databases": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Set: pluginsdk.HashString,
		},

		"readonly_endpoint_failover_policy_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"read_write_endpoint_failover_policy": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"mode": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"grace_minutes": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (d MsSqlFailoverGroupDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.FailoverGroupsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state MsSqlFailoverGroupDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			serverId, err := commonids.ParseSqlServerID(state.ServerId)
			if err != nil {
				return err
			}

			id := failovergroups.NewFailoverGroupID(subscriptionId, serverId.ResourceGroupName, serverId.ServerName, state.Name)
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			state.Name = id.FailoverGroupName
			state.ServerId = serverId.ID()

			if model := existing.Model; model != nil {
				state.Tags = pointer.From(model.Tags)
				if props := model.Properties; props != nil {
					state.Databases = pointer.From(props.Databases)
					state.PartnerServers = flattenPartnerServers(props.PartnerServers)

					if props.ReadOnlyEndpoint != nil && pointer.From(props.ReadOnlyEndpoint.FailoverPolicy) == failovergroups.ReadOnlyEndpointFailoverPolicyEnabled {
						state.ReadonlyEndpointFailurePolicyEnabled = true
					}

					state.ReadWriteEndpointFailurePolicy = []ReadWriteEndpointFailurePolicyDataSourceModel{{
						Mode: string(props.ReadWriteEndpoint.FailoverPolicy),
					}}

					state.ReadWriteEndpointFailurePolicy[0].GraceMinutes = pointer.From(props.ReadWriteEndpoint.FailoverWithDataLossGracePeriodMinutes)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenPartnerServers(input []failovergroups.PartnerInfo) []PartnerServerDataSourceModel {
	output := make([]PartnerServerDataSourceModel, 0)
	if input == nil {
		return output
	}

	for _, partner := range input {
		model := PartnerServerDataSourceModel{
			Location: location.NormalizeNilable(partner.Location),
			Role:     string(pointer.From(partner.ReplicationRole)),
			ID:       partner.Id,
		}

		output = append(output, model)
	}

	return output
}
