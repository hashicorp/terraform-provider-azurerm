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
)

var _ sdk.DataSource = DevCenterProjectPoolDataSource{}

type DevCenterProjectPoolDataSource struct{}

type DevCenterProjectPoolDataSourceModel struct {
	Name                               string            `tfschema:"name"`
	Location                           string            `tfschema:"location"`
	DevCenterProjectId                 string            `tfschema:"dev_center_project_id"`
	DevBoxDefinitionName               string            `tfschema:"dev_box_definition_name"`
	LocalAdministratorEnabled          bool              `tfschema:"local_administrator_enabled"`
	DevCenterAttachedNetworkName       string            `tfschema:"dev_center_attached_network_name"`
	StopOnDisconnectGracePeriodMinutes int64             `tfschema:"stop_on_disconnect_grace_period_minutes"`
	Tags                               map[string]string `tfschema:"tags"`
}

func (DevCenterProjectPoolDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.DevCenterProjectPoolName,
		},

		"dev_center_project_id": commonschema.ResourceIDReferenceRequired(&pools.ProjectId{}),
	}
}

func (DevCenterProjectPoolDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dev_box_definition_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"dev_center_attached_network_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"local_administrator_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"location": commonschema.LocationComputed(),

		"stop_on_disconnect_grace_period_minutes": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (DevCenterProjectPoolDataSource) ModelObject() interface{} {
	return &DevCenterProjectPoolDataSourceModel{}
}

func (DevCenterProjectPoolDataSource) ResourceType() string {
	return "azurerm_dev_center_project_pool"
}

func (r DevCenterProjectPoolDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Pools
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DevCenterProjectPoolDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			devCenterProjectId, err := projects.ParseProjectID(state.DevCenterProjectId)
			if err != nil {
				return err
			}

			id := pools.NewPoolID(subscriptionId, devCenterProjectId.ResourceGroupName, devCenterProjectId.ProjectName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				state.Name = id.PoolName
				state.DevCenterProjectId = projects.NewProjectID(id.SubscriptionId, id.ResourceGroupName, id.ProjectName).ID()
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.DevBoxDefinitionName = pointer.From(props.DevBoxDefinitionName)
					state.LocalAdministratorEnabled = pointer.From(props.LocalAdministrator) == pools.LocalAdminStatusEnabled
					state.DevCenterAttachedNetworkName = pointer.From(props.NetworkConnectionName)
					state.StopOnDisconnectGracePeriodMinutes = flattenDevCenterProjectPoolStopOnDisconnectForDataSource(props.StopOnDisconnect)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenDevCenterProjectPoolStopOnDisconnectForDataSource(input *pools.StopOnDisconnectConfiguration) int64 {
	var gracePeriodMinutes int64
	if input == nil || pointer.From(input.Status) == pools.StopOnDisconnectEnableStatusDisabled {
		return gracePeriodMinutes
	}

	return pointer.From(input.GracePeriodMinutes)
}
