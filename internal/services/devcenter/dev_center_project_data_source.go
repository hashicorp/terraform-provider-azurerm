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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projects"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = DevCenterProjectDataSource{}

type DevCenterProjectDataSource struct{}

type DevCenterProjectDataSourceModel struct {
	Description            string                                     `tfschema:"description"`
	DevCenterId            string                                     `tfschema:"dev_center_id"`
	DevCenterUri           string                                     `tfschema:"dev_center_uri"`
	Identity               []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Location               string                                     `tfschema:"location"`
	MaximumDevBoxesPerUser int64                                      `tfschema:"maximum_dev_boxes_per_user"`
	Name                   string                                     `tfschema:"name"`
	ResourceGroupName      string                                     `tfschema:"resource_group_name"`
	Tags                   map[string]interface{}                     `tfschema:"tags"`
}

func (DevCenterProjectDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (DevCenterProjectDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"dev_center_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"dev_center_uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

		"location": commonschema.LocationComputed(),

		"maximum_dev_boxes_per_user": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (DevCenterProjectDataSource) ModelObject() interface{} {
	return &DevCenterProjectDataSourceModel{}
}

func (DevCenterProjectDataSource) ResourceType() string {
	return "azurerm_dev_center_project"
}

func (r DevCenterProjectDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Projects
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DevCenterProjectDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := projects.NewProjectID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				state.Name = id.ProjectName
				state.ResourceGroupName = id.ResourceGroupName
				state.Location = location.Normalize(model.Location)
				state.Tags = tags.Flatten(model.Tags)

				identity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %v", err)
				}
				state.Identity = pointer.From(identity)

				if props := model.Properties; props != nil {
					state.Description = pointer.From(props.Description)
					state.DevCenterId = pointer.From(props.DevCenterId)
					state.DevCenterUri = pointer.From(props.DevCenterUri)
					state.MaximumDevBoxesPerUser = pointer.From(props.MaxDevBoxesPerUser)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
