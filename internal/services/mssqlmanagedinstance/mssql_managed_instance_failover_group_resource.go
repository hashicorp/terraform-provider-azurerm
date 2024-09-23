// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/instancefailovergroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MsSqlManagedInstanceFailoverGroupModel struct {
	Name                                  string `tfschema:"name"`
	Location                              string `tfschema:"location"`
	ManagedInstanceId                     string `tfschema:"managed_instance_id"`
	PartnerManagedInstanceId              string `tfschema:"partner_managed_instance_id"`
	ReadOnlyEndpointFailoverPolicyEnabled bool   `tfschema:"readonly_endpoint_failover_policy_enabled"`

	ReadWriteEndpointFailurePolicy []MsSqlManagedInstanceReadWriteEndpointFailurePolicyModel `tfschema:"read_write_endpoint_failover_policy"`

	PartnerRegion []MsSqlManagedInstancePartnerRegionModel `tfschema:"partner_region"`
	Role          string                                   `tfschema:"role"`
}

type MsSqlManagedInstanceReadWriteEndpointFailurePolicyModel struct {
	GraceMinutes int64  `tfschema:"grace_minutes"`
	Mode         string `tfschema:"mode"`
}

type MsSqlManagedInstancePartnerRegionModel struct {
	Location string `tfschema:"location"`
	Role     string `tfschema:"role"`
}

var _ sdk.Resource = MsSqlManagedInstanceFailoverGroupResource{}
var _ sdk.ResourceWithUpdate = MsSqlManagedInstanceFailoverGroupResource{}

type MsSqlManagedInstanceFailoverGroupResource struct{}

func (r MsSqlManagedInstanceFailoverGroupResource) ResourceType() string {
	return "azurerm_mssql_managed_instance_failover_group"
}

func (r MsSqlManagedInstanceFailoverGroupResource) ModelObject() interface{} {
	return &MsSqlManagedInstanceFailoverGroupModel{}
}

func (r MsSqlManagedInstanceFailoverGroupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ManagedInstanceFailoverGroupID
}

func (r MsSqlManagedInstanceFailoverGroupResource) Arguments() map[string]*pluginsdk.Schema {
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
			ValidateFunc: validate.ManagedInstanceID,
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
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(instancefailovergroups.ReadWriteEndpointFailoverPolicyAutomatic),
							string(instancefailovergroups.ReadWriteEndpointFailoverPolicyManual),
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
	}
}

func (r MsSqlManagedInstanceFailoverGroupResource) Attributes() map[string]*pluginsdk.Schema {
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

func (r MsSqlManagedInstanceFailoverGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstanceFailoverGroupsClient

			var model MsSqlManagedInstanceFailoverGroupModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			managedInstanceId, err := commonids.ParseSqlManagedInstanceID(model.ManagedInstanceId)
			if err != nil {
				return err
			}

			id := instancefailovergroups.NewInstanceFailoverGroupID(managedInstanceId.SubscriptionId,
				managedInstanceId.ResourceGroupName, model.Location, model.Name)

			partnerId, err := parse.ManagedInstanceID(model.PartnerManagedInstanceId)
			if err != nil {
				return err
			}

			instancesClient := metadata.Client.MSSQLManagedInstance.ManagedInstancesClientForSubscription(partnerId.SubscriptionId)
			partner, err := instancesClient.Get(ctx, partnerId.ResourceGroup, partnerId.Name, "")
			if err != nil || partner.Location == nil || *partner.Location == "" {
				return fmt.Errorf("checking for existence and region of Partner of %q: %+v", id, err)
			}

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			readOnlyFailoverPolicy := instancefailovergroups.ReadOnlyEndpointFailoverPolicyDisabled
			if model.ReadOnlyEndpointFailoverPolicyEnabled {
				readOnlyFailoverPolicy = instancefailovergroups.ReadOnlyEndpointFailoverPolicyEnabled
			}

			parameters := instancefailovergroups.InstanceFailoverGroup{
				Properties: &instancefailovergroups.InstanceFailoverGroupProperties{
					ReadOnlyEndpoint: &instancefailovergroups.InstanceFailoverGroupReadOnlyEndpoint{
						FailoverPolicy: &readOnlyFailoverPolicy,
					},
					ReadWriteEndpoint: instancefailovergroups.InstanceFailoverGroupReadWriteEndpoint{},
					PartnerRegions: []instancefailovergroups.PartnerRegionInfo{
						{
							Location: partner.Location,
						},
					},
					ManagedInstancePairs: []instancefailovergroups.ManagedInstancePairInfo{
						{
							PrimaryManagedInstanceId: utils.String(managedInstanceId.ID()),
							PartnerManagedInstanceId: utils.String(partnerId.ID()),
						},
					},
				},
			}

			if rwPolicy := model.ReadWriteEndpointFailurePolicy; len(rwPolicy) > 0 {
				parameters.Properties.ReadWriteEndpoint.FailoverPolicy = instancefailovergroups.ReadWriteEndpointFailoverPolicy(rwPolicy[0].Mode)
				if rwPolicy[0].Mode == string(instancefailovergroups.ReadWriteEndpointFailoverPolicyAutomatic) {
					parameters.Properties.ReadWriteEndpoint.FailoverWithDataLossGracePeriodMinutes = &rwPolicy[0].GraceMinutes
				}
			}

			metadata.Logger.Infof("Creating %s", id)

			err = client.CreateOrUpdateThenPoll(ctx, id, parameters)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r MsSqlManagedInstanceFailoverGroupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstanceFailoverGroupsClient

			id, err := instancefailovergroups.ParseInstanceFailoverGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state MsSqlManagedInstanceFailoverGroupModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			managedInstanceId, err := commonids.ParseSqlManagedInstanceID(state.ManagedInstanceId)
			if err != nil {
				return err
			}

			partnerId, err := parse.ManagedInstanceID(state.PartnerManagedInstanceId)
			if err != nil {
				return err
			}

			instancesClient := metadata.Client.MSSQLManagedInstance.ManagedInstancesClientForSubscription(partnerId.SubscriptionId)
			partner, err := instancesClient.Get(ctx, partnerId.ResourceGroup, partnerId.Name, "")
			if err != nil || partner.Location == nil || *partner.Location == "" {
				return fmt.Errorf("checking for existence and region of Partner of %q: %+v", id, err)
			}

			readOnlyFailoverPolicy := instancefailovergroups.ReadOnlyEndpointFailoverPolicyDisabled
			if state.ReadOnlyEndpointFailoverPolicyEnabled {
				readOnlyFailoverPolicy = instancefailovergroups.ReadOnlyEndpointFailoverPolicyEnabled
			}

			parameters := instancefailovergroups.InstanceFailoverGroup{
				Properties: &instancefailovergroups.InstanceFailoverGroupProperties{
					ReadOnlyEndpoint: &instancefailovergroups.InstanceFailoverGroupReadOnlyEndpoint{
						FailoverPolicy: &readOnlyFailoverPolicy,
					},
					ReadWriteEndpoint: instancefailovergroups.InstanceFailoverGroupReadWriteEndpoint{},
					PartnerRegions: []instancefailovergroups.PartnerRegionInfo{
						{
							Location: partner.Location,
						},
					},
					ManagedInstancePairs: []instancefailovergroups.ManagedInstancePairInfo{
						{
							PrimaryManagedInstanceId: utils.String(managedInstanceId.ID()),
							PartnerManagedInstanceId: utils.String(partnerId.ID()),
						},
					},
				},
			}

			if rwPolicy := state.ReadWriteEndpointFailurePolicy; len(rwPolicy) > 0 {
				parameters.Properties.ReadWriteEndpoint.FailoverPolicy = instancefailovergroups.ReadWriteEndpointFailoverPolicy(rwPolicy[0].Mode)
				if rwPolicy[0].Mode == string(instancefailovergroups.ReadWriteEndpointFailoverPolicyAutomatic) {
					parameters.Properties.ReadWriteEndpoint.FailoverWithDataLossGracePeriodMinutes = &rwPolicy[0].GraceMinutes
				}
			}

			metadata.Logger.Infof("Updating %s", id)

			err = client.CreateOrUpdateThenPoll(ctx, *id, parameters)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r MsSqlManagedInstanceFailoverGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstanceFailoverGroupsClient

			id, err := instancefailovergroups.ParseInstanceFailoverGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state MsSqlManagedInstanceFailoverGroupModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			model := MsSqlManagedInstanceFailoverGroupModel{
				Name:     id.InstanceFailoverGroupName,
				Location: id.LocationName,
			}

			if result.Model != nil {
				if props := result.Model.Properties; props != nil {
					model.Role = string(pointer.From(props.ReplicationRole))

					if instancePairs := props.ManagedInstancePairs; len(instancePairs) == 1 {
						if primaryId := instancePairs[0].PrimaryManagedInstanceId; primaryId != nil {
							id, err := parse.ManagedInstanceIDInsensitively(*primaryId)
							if err != nil {
								return fmt.Errorf("parsing `PrimaryManagedInstanceID` from response: %v", err)
							}

							model.ManagedInstanceId = id.ID()
						}

						if partnerId := instancePairs[0].PartnerManagedInstanceId; partnerId != nil {
							id, err := parse.ManagedInstanceIDInsensitively(*partnerId)
							if err != nil {
								return fmt.Errorf("parsing `PrimaryManagedInstanceID` from response: %v", err)
							}

							model.PartnerManagedInstanceId = id.ID()
						}
					}

					for _, partnerRegion := range props.PartnerRegions {
						var location string
						if partnerRegion.Location != nil {
							location = *partnerRegion.Location
						}

						model.PartnerRegion = append(model.PartnerRegion, MsSqlManagedInstancePartnerRegionModel{
							Location: location,
							Role:     string(pointer.From(partnerRegion.ReplicationRole)),
						})
					}

					if readOnlyEndpoint := props.ReadOnlyEndpoint; readOnlyEndpoint != nil {
						if *readOnlyEndpoint.FailoverPolicy == instancefailovergroups.ReadOnlyEndpointFailoverPolicyEnabled {
							model.ReadOnlyEndpointFailoverPolicyEnabled = true
						}
					}

					model.ReadWriteEndpointFailurePolicy = []MsSqlManagedInstanceReadWriteEndpointFailurePolicyModel{
						{
							Mode:         string(props.ReadWriteEndpoint.FailoverPolicy),
							GraceMinutes: pointer.From(props.ReadWriteEndpoint.FailoverWithDataLossGracePeriodMinutes),
						},
					}
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r MsSqlManagedInstanceFailoverGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstanceFailoverGroupsClient

			id, err := instancefailovergroups.ParseInstanceFailoverGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			err = client.DeleteThenPoll(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
