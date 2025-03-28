// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/pools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projects"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = DevCenterProjectPoolResource{}
	_ sdk.ResourceWithUpdate = DevCenterProjectPoolResource{}
)

type DevCenterProjectPoolResource struct{}

func (r DevCenterProjectPoolResource) ModelObject() interface{} {
	return &DevCenterProjectPoolResourceModel{}
}

type DevCenterProjectPoolResourceModel struct {
	Name                               string            `tfschema:"name"`
	Location                           string            `tfschema:"location"`
	DevCenterProjectId                 string            `tfschema:"dev_center_project_id"`
	DevBoxDefinitionName               string            `tfschema:"dev_box_definition_name"`
	LocalAdministratorEnabled          bool              `tfschema:"local_administrator_enabled"`
	DevCenterAttachedNetworkName       string            `tfschema:"dev_center_attached_network_name"`
	StopOnDisconnectGracePeriodMinutes int64             `tfschema:"stop_on_disconnect_grace_period_minutes"`
	Tags                               map[string]string `tfschema:"tags"`
}

func (r DevCenterProjectPoolResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return pools.ValidatePoolID
}

func (r DevCenterProjectPoolResource) ResourceType() string {
	return "azurerm_dev_center_project_pool"
}

func (r DevCenterProjectPoolResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DevCenterProjectPoolName,
		},

		"location": commonschema.Location(),

		"dev_center_project_id": commonschema.ResourceIDReferenceRequiredForceNew(&pools.ProjectId{}),

		"dev_box_definition_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.DevCenterDevBoxDefinitionName,
		},

		"local_administrator_enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},

		"dev_center_attached_network_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"stop_on_disconnect_grace_period_minutes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(60, 480),
		},

		"tags": commonschema.Tags(),
	}
}

func (r DevCenterProjectPoolResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DevCenterProjectPoolResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Pools
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DevCenterProjectPoolResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			devCenterProjectId, err := projects.ParseProjectID(model.DevCenterProjectId)
			if err != nil {
				return err
			}

			id := pools.NewPoolID(subscriptionId, devCenterProjectId.ResourceGroupName, devCenterProjectId.ProjectName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := pools.Pool{
				Location: location.Normalize(model.Location),
				Properties: &pools.PoolProperties{
					DevBoxDefinitionName:  pointer.To(model.DevBoxDefinitionName),
					NetworkConnectionName: pointer.To(model.DevCenterAttachedNetworkName),
					LicenseType:           pointer.To(pools.LicenseTypeWindowsClient),
					StopOnDisconnect:      expandDevCenterProjectPoolStopOnDisconnect(model.StopOnDisconnectGracePeriodMinutes),
				},
				Tags: pointer.To(model.Tags),
			}

			if model.LocalAdministratorEnabled {
				parameters.Properties.LocalAdministrator = pointer.To(pools.LocalAdminStatusEnabled)
			} else {
				parameters.Properties.LocalAdministrator = pointer.To(pools.LocalAdminStatusDisabled)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DevCenterProjectPoolResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Pools

			id, err := pools.ParsePoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := DevCenterProjectPoolResourceModel{
				Name:               id.PoolName,
				DevCenterProjectId: projects.NewProjectID(id.SubscriptionId, id.ResourceGroupName, id.ProjectName).ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.DevBoxDefinitionName = pointer.From(props.DevBoxDefinitionName)
					state.LocalAdministratorEnabled = pointer.From(props.LocalAdministrator) == pools.LocalAdminStatusEnabled
					state.DevCenterAttachedNetworkName = pointer.From(props.NetworkConnectionName)
					state.StopOnDisconnectGracePeriodMinutes = flattenDevCenterProjectPoolStopOnDisconnect(props.StopOnDisconnect)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DevCenterProjectPoolResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Pools

			id, err := pools.ParsePoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model DevCenterProjectPoolResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := pools.PoolUpdate{
				Properties: &pools.PoolUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("dev_box_definition_name") {
				parameters.Properties.DevBoxDefinitionName = pointer.To(model.DevBoxDefinitionName)
			}

			if metadata.ResourceData.HasChange("local_administrator_enabled") {
				if model.LocalAdministratorEnabled {
					parameters.Properties.LocalAdministrator = pointer.To(pools.LocalAdminStatusEnabled)
				} else {
					parameters.Properties.LocalAdministrator = pointer.To(pools.LocalAdminStatusDisabled)
				}
			}

			if metadata.ResourceData.HasChange("dev_center_attached_network_name") {
				parameters.Properties.NetworkConnectionName = pointer.To(model.DevCenterAttachedNetworkName)
			}

			if metadata.ResourceData.HasChange("stop_on_disconnect_grace_period_minutes") {
				parameters.Properties.StopOnDisconnect = expandDevCenterProjectPoolStopOnDisconnect(model.StopOnDisconnectGracePeriodMinutes)
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = pointer.To(model.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r DevCenterProjectPoolResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Pools

			id, err := pools.ParsePoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandDevCenterProjectPoolStopOnDisconnect(input int64) *pools.StopOnDisconnectConfiguration {
	if input == 0 {
		return &pools.StopOnDisconnectConfiguration{
			GracePeriodMinutes: pointer.To(int64(0)),
			Status:             pointer.To(pools.StopOnDisconnectEnableStatusDisabled),
		}
	}

	return &pools.StopOnDisconnectConfiguration{
		GracePeriodMinutes: pointer.To(input),
		Status:             pointer.To(pools.StopOnDisconnectEnableStatusEnabled),
	}
}

func flattenDevCenterProjectPoolStopOnDisconnect(input *pools.StopOnDisconnectConfiguration) int64 {
	var gracePeriodMinutes int64
	if input == nil || pointer.From(input.Status) == pools.StopOnDisconnectEnableStatusDisabled {
		return gracePeriodMinutes
	}

	return pointer.From(input.GracePeriodMinutes)
}
