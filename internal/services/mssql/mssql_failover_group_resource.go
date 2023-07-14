// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
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
	GraceMinutes int32  `tfschema:"grace_minutes"`
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
							string(sql.ReadWriteEndpointFailoverPolicyAutomatic),
							string(sql.ReadWriteEndpointFailoverPolicyManual),
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
				if rwPolicy[0].Mode == string(sql.ReadWriteEndpointFailoverPolicyAutomatic) && rwPolicy[0].GraceMinutes < 60 {
					return fmt.Errorf("`grace_minutes` should be %d or greater when `mode` is %q", 60, sql.ReadWriteEndpointFailoverPolicyAutomatic)
				}
				if rwPolicy[0].Mode == string(sql.ReadWriteEndpointFailoverPolicyManual) && rwPolicy[0].GraceMinutes > 0 {
					return fmt.Errorf("`grace_minutes` should not be specified when `mode` is %q", sql.ReadWriteEndpointFailoverPolicyManual)
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

			serverId, err := parse.ServerID(model.ServerId)
			if err != nil {
				return err
			}

			if _, err = serversClient.Get(ctx, serverId.ResourceGroup, serverId.Name, ""); err != nil {
				return fmt.Errorf("retrieving %s: %+v", serverId, err)
			}

			id := parse.NewFailoverGroupID(subscriptionId, serverId.ResourceGroup, serverId.Name, model.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			readOnlyFailoverPolicy := sql.ReadOnlyEndpointFailoverPolicyDisabled
			if model.ReadonlyEndpointFailurePolicyEnabled {
				readOnlyFailoverPolicy = sql.ReadOnlyEndpointFailoverPolicyEnabled
			}

			properties := sql.FailoverGroup{
				FailoverGroupProperties: &sql.FailoverGroupProperties{
					Databases: &model.Databases,
					ReadOnlyEndpoint: &sql.FailoverGroupReadOnlyEndpoint{
						FailoverPolicy: readOnlyFailoverPolicy,
					},
					ReadWriteEndpoint: &sql.FailoverGroupReadWriteEndpoint{},
					PartnerServers:    r.expandPartnerServers(model.PartnerServers),
				},
				Tags: tags.FromTypedObject(model.Tags),
			}

			if rwPolicy := model.ReadWriteEndpointFailurePolicy; len(rwPolicy) > 0 {
				properties.FailoverGroupProperties.ReadWriteEndpoint.FailoverPolicy = sql.ReadWriteEndpointFailoverPolicy(rwPolicy[0].Mode)
				if rwPolicy[0].Mode == string(sql.ReadWriteEndpointFailoverPolicyAutomatic) {
					properties.FailoverGroupProperties.ReadWriteEndpoint.FailoverWithDataLossGracePeriodMinutes = utils.Int32(rwPolicy[0].GraceMinutes)
				}
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.Name, properties)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
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

			id, err := parse.FailoverGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Info("Decoding state...")
			var state MsSqlFailoverGroupModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			metadata.Logger.Infof("updating %s", id)

			readOnlyFailoverPolicy := sql.ReadOnlyEndpointFailoverPolicyDisabled
			if state.ReadonlyEndpointFailurePolicyEnabled {
				readOnlyFailoverPolicy = sql.ReadOnlyEndpointFailoverPolicyEnabled
			}

			properties := sql.FailoverGroup{
				FailoverGroupProperties: &sql.FailoverGroupProperties{
					Databases: &state.Databases,
					ReadOnlyEndpoint: &sql.FailoverGroupReadOnlyEndpoint{
						FailoverPolicy: readOnlyFailoverPolicy,
					},
					ReadWriteEndpoint: &sql.FailoverGroupReadWriteEndpoint{
						FailoverPolicy: sql.ReadWriteEndpointFailoverPolicy(state.ReadWriteEndpointFailurePolicy[0].Mode),
					},
					PartnerServers: r.expandPartnerServers(state.PartnerServers),
				},
				Tags: tags.FromTypedObject(state.Tags),
			}

			if state.ReadWriteEndpointFailurePolicy[0].Mode == string(sql.ReadWriteEndpointFailoverPolicyAutomatic) {
				properties.FailoverGroupProperties.ReadWriteEndpoint.FailoverWithDataLossGracePeriodMinutes = utils.Int32(state.ReadWriteEndpointFailurePolicy[0].GraceMinutes)
			}

			// client.Update doesn't support changing the PartnerServers
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.Name, properties)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", id, err)
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

			id, err := parse.FailoverGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			serverId := parse.NewServerID(subscriptionId, id.ResourceGroup, id.ServerName)

			model := MsSqlFailoverGroupModel{
				Name:     id.Name,
				ServerId: serverId.ID(),
				Tags:     tags.ToTypedObject(existing.Tags),
			}

			if props := existing.FailoverGroupProperties; props != nil {
				if props.Databases != nil {
					model.Databases = *props.Databases
				}

				model.PartnerServers = r.flattenPartnerServers(props.PartnerServers)

				if props.ReadOnlyEndpoint != nil && props.ReadOnlyEndpoint.FailoverPolicy == sql.ReadOnlyEndpointFailoverPolicyEnabled {
					model.ReadonlyEndpointFailurePolicyEnabled = true
				}

				if endpoint := props.ReadWriteEndpoint; endpoint != nil {
					model.ReadWriteEndpointFailurePolicy = []ReadWriteEndpointFailurePolicyModel{{
						Mode: string(endpoint.FailoverPolicy),
					}}

					if endpoint.FailoverWithDataLossGracePeriodMinutes != nil {
						model.ReadWriteEndpointFailurePolicy[0].GraceMinutes = *endpoint.FailoverWithDataLossGracePeriodMinutes
					}
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

			id, err := parse.FailoverGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if existing, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name); err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.ServerName, id.Name)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r MsSqlFailoverGroupResource) flattenPartnerServers(input *[]sql.PartnerInfo) []PartnerServerModel {
	output := make([]PartnerServerModel, 0)
	if input == nil {
		return output
	}

	for _, partner := range *input {
		model := PartnerServerModel{
			Location: location.NormalizeNilable(partner.Location),
			Role:     string(partner.ReplicationRole),
		}
		if partner.ID != nil {
			model.ID = *partner.ID
		}
		output = append(output, model)
	}

	return output
}

func (r MsSqlFailoverGroupResource) expandPartnerServers(input []PartnerServerModel) *[]sql.PartnerInfo {
	var partnerServers []sql.PartnerInfo
	if input == nil {
		return &partnerServers
	}

	for _, v := range input {
		partnerServers = append(partnerServers, sql.PartnerInfo{
			ID: utils.String(v.ID),
		})
	}

	return &partnerServers
}
