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
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/servers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MsSqlFailoverGroupModel struct {
	Databases                            []string             `tfschema:"databases"`
	Name                                 string               `tfschema:"name"`
	PartnerServers                       []PartnerServerModel `tfschema:"partner_server"`
	ReadonlyEndpointFailurePolicyEnabled bool                 `tfschema:"readonly_endpoint_failover_policy_enabled"`
	ServerId                             string               `tfschema:"server_id"`
	Tags                                 map[string]string    `tfschema:"tags"`

	ReadWriteEndpointFailurePolicy []ReadWriteEndpointFailurePolicyModel `tfschema:"read_write_endpoint_failover_policy"`
}

type PartnerServerModel struct {
	ID       string `tfschema:"id"`
	Location string `tfschema:"location"`
	Role     string `tfschema:"role"`
}

type ReadWriteEndpointFailurePolicyModel struct {
	GraceMinutes int64  `tfschema:"grace_minutes"`
	Mode         string `tfschema:"mode"`
}

var (
	_ sdk.Resource           = MsSqlFailoverGroupResource{}
	_ sdk.ResourceWithUpdate = MsSqlFailoverGroupResource{}
)

type MsSqlFailoverGroupResource struct{}

func (r MsSqlFailoverGroupResource) ResourceType() string {
	return "azurerm_mssql_failover_group"
}

func (r MsSqlFailoverGroupResource) ModelObject() interface{} {
	return &MsSqlFailoverGroupModel{}
}

func (r MsSqlFailoverGroupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.FailoverGroupID
}

func (r MsSqlFailoverGroupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlFailoverGroupName,
		},

		"server_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ServerID,
		},

		"partner_server": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
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
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Set: pluginsdk.HashString,
		},

		"readonly_endpoint_failover_policy_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"read_write_endpoint_failover_policy": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"mode": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(failovergroups.ReadWriteEndpointFailoverPolicyAutomatic),
							string(failovergroups.ReadWriteEndpointFailoverPolicyManual),
						}, false),
					},
					"grace_minutes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(60),
					},
				},
			},
		},

		"tags": tags.Schema(),
	}
}

func (r MsSqlFailoverGroupResource) Attributes() map[string]*pluginsdk.Schema {
	return nil
}

func (r MsSqlFailoverGroupResource) CustomizeDiff() sdk.ResourceFunc {
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

func (r MsSqlFailoverGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.FailoverGroupsClient
			serversClient := metadata.Client.MSSQL.ServersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model MsSqlFailoverGroupModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			serverId, err := commonids.ParseSqlServerID(model.ServerId)
			if err != nil {
				return err
			}

			if _, err = serversClient.Get(ctx, *serverId, servers.DefaultGetOperationOptions()); err != nil {
				return fmt.Errorf("retrieving %s: %+v", serverId, err)
			}

			id := failovergroups.NewFailoverGroupID(subscriptionId, serverId.ResourceGroupName, serverId.ServerName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			readOnlyFailoverPolicy := failovergroups.ReadOnlyEndpointFailoverPolicyDisabled
			if model.ReadonlyEndpointFailurePolicyEnabled {
				readOnlyFailoverPolicy = failovergroups.ReadOnlyEndpointFailoverPolicyEnabled
			}

			properties := failovergroups.FailoverGroup{
				Properties: &failovergroups.FailoverGroupProperties{
					Databases: &model.Databases,
					ReadOnlyEndpoint: &failovergroups.FailoverGroupReadOnlyEndpoint{
						FailoverPolicy: &readOnlyFailoverPolicy,
					},
					ReadWriteEndpoint: failovergroups.FailoverGroupReadWriteEndpoint{},
					PartnerServers:    r.expandPartnerServers(model.PartnerServers),
				},
				Tags: pointer.To(model.Tags),
			}

			if rwPolicy := model.ReadWriteEndpointFailurePolicy; len(rwPolicy) > 0 {
				properties.Properties.ReadWriteEndpoint.FailoverPolicy = failovergroups.ReadWriteEndpointFailoverPolicy(rwPolicy[0].Mode)
				if rwPolicy[0].Mode == string(failovergroups.ReadWriteEndpointFailoverPolicyAutomatic) {
					properties.Properties.ReadWriteEndpoint.FailoverWithDataLossGracePeriodMinutes = utils.Int64(rwPolicy[0].GraceMinutes)
				}
			}

			err = client.CreateOrUpdateThenPoll(ctx, id, properties)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r MsSqlFailoverGroupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.FailoverGroupsClient

			id, err := failovergroups.ParseFailoverGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Info("Decoding state...")
			var state MsSqlFailoverGroupModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			metadata.Logger.Infof("updating %s", id)

			readOnlyFailoverPolicy := failovergroups.ReadOnlyEndpointFailoverPolicyDisabled
			if state.ReadonlyEndpointFailurePolicyEnabled {
				readOnlyFailoverPolicy = failovergroups.ReadOnlyEndpointFailoverPolicyEnabled
			}

			properties := failovergroups.FailoverGroup{
				Properties: &failovergroups.FailoverGroupProperties{
					Databases: &state.Databases,
					ReadOnlyEndpoint: &failovergroups.FailoverGroupReadOnlyEndpoint{
						FailoverPolicy: &readOnlyFailoverPolicy,
					},
					ReadWriteEndpoint: failovergroups.FailoverGroupReadWriteEndpoint{
						FailoverPolicy: failovergroups.ReadWriteEndpointFailoverPolicy(state.ReadWriteEndpointFailurePolicy[0].Mode),
					},
					PartnerServers: r.expandPartnerServers(state.PartnerServers),
				},
				Tags: pointer.To(state.Tags),
			}

			if state.ReadWriteEndpointFailurePolicy[0].Mode == string(failovergroups.ReadWriteEndpointFailoverPolicyAutomatic) {
				properties.Properties.ReadWriteEndpoint.FailoverWithDataLossGracePeriodMinutes = pointer.To(state.ReadWriteEndpointFailurePolicy[0].GraceMinutes)
			}

			// client.Update doesn't support changing the PartnerServers
			err = client.CreateOrUpdateThenPoll(ctx, *id, properties)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r MsSqlFailoverGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.MSSQL.FailoverGroupsClient

			id, err := failovergroups.ParseFailoverGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			serverId := parse.NewServerID(subscriptionId, id.ResourceGroupName, id.ServerName)

			model := MsSqlFailoverGroupModel{
				Name:     id.FailoverGroupName,
				ServerId: serverId.ID(),
			}

			if existing.Model != nil {

				model.Tags = pointer.From(existing.Model.Tags)

				if props := existing.Model.Properties; props != nil {
					if props.Databases != nil {
						model.Databases = *props.Databases
					}

					model.PartnerServers = r.flattenPartnerServers(props.PartnerServers)

					if props.ReadOnlyEndpoint != nil && pointer.From(props.ReadOnlyEndpoint.FailoverPolicy) == failovergroups.ReadOnlyEndpointFailoverPolicyEnabled {
						model.ReadonlyEndpointFailurePolicyEnabled = true
					}

					model.ReadWriteEndpointFailurePolicy = []ReadWriteEndpointFailurePolicyModel{{
						Mode: string(props.ReadWriteEndpoint.FailoverPolicy),
					}}

					model.ReadWriteEndpointFailurePolicy[0].GraceMinutes = pointer.From(props.ReadWriteEndpoint.FailoverWithDataLossGracePeriodMinutes)

				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r MsSqlFailoverGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.FailoverGroupsClient

			id, err := failovergroups.ParseFailoverGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if existing, err := client.Get(ctx, *id); err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			err = client.DeleteThenPoll(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r MsSqlFailoverGroupResource) flattenPartnerServers(input []failovergroups.PartnerInfo) []PartnerServerModel {
	output := make([]PartnerServerModel, 0)
	if input == nil {
		return output
	}

	for _, partner := range input {
		model := PartnerServerModel{
			Location: location.NormalizeNilable(partner.Location),
			Role:     string(pointer.From(partner.ReplicationRole)),
			ID:       partner.Id,
		}

		output = append(output, model)
	}

	return output
}

func (r MsSqlFailoverGroupResource) expandPartnerServers(input []PartnerServerModel) []failovergroups.PartnerInfo {
	var partnerServers []failovergroups.PartnerInfo
	if input == nil {
		return partnerServers
	}

	for _, v := range input {
		partnerServers = append(partnerServers, failovergroups.PartnerInfo{
			Id: v.ID,
		})
	}

	return partnerServers
}
