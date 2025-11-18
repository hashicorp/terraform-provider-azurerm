// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dashboard

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2025-08-01/manageddashboards"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DashboardMonitorResourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`
}

type DashboardMonitorResource struct{}

var _ sdk.Resource = DashboardMonitorResource{}

func (r DashboardMonitorResource) ResourceType() string {
	return "azurerm_dashboard_monitor"
}

func (r DashboardMonitorResource) ModelObject() interface{} {
	return &DashboardMonitorResourceModel{}
}

func (r DashboardMonitorResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return manageddashboards.ValidateDashboardID
}

func (r DashboardMonitorResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z][a-z0-9A-Z-]{0,28}[a-z0-9A-Z]$`),
				"The name must be 2-30 characters long, start with a letter, end with a letter or number, and can only contain letters, numbers, and hyphens.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (r DashboardMonitorResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DashboardMonitorResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model DashboardMonitorResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Dashboard.ManagedDashboardsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := manageddashboards.NewDashboardID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.DashboardsGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := manageddashboards.ManagedDashboard{
				Location:   location.Normalize(model.Location),
				Tags:       &model.Tags,
				Properties: &manageddashboards.ManagedDashboardProperties{},
			}

			if err := client.CreateThenPoll(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DashboardMonitorResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dashboard.ManagedDashboardsClient

			id, err := manageddashboards.ParseDashboardID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model DashboardMonitorResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("tags") {
				updateParams := manageddashboards.ManagedDashboardUpdateParameters{
					Tags: &model.Tags,
				}

				if _, err := client.Update(ctx, *id, updateParams); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

func (r DashboardMonitorResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dashboard.ManagedDashboardsClient

			id, err := manageddashboards.ParseDashboardID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.DashboardsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := DashboardMonitorResourceModel{
				Name:              id.DashboardName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DashboardMonitorResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dashboard.ManagedDashboardsClient

			id, err := manageddashboards.ParseDashboardID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
