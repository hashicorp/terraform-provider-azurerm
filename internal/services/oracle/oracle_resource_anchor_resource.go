// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/resourceanchors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = ResourceAnchorResource{}

type ResourceAnchorResource struct{}

type ResourceAnchorResourceModel struct {
	Location            string            `tfschema:"location"`
	Name                string            `tfschema:"name"`
	ResourceGroupName   string            `tfschema:"resource_group_name"`
	LinkedCompartmentID string            `tfschema:"linked_compartment_id"`
	Tags                map[string]string `tfschema:"tags"`
}

func (ResourceAnchorResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ResourceAnchorName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"tags": commonschema.Tags(),
	}
}

func (ResourceAnchorResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"linked_compartment_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (ResourceAnchorResource) ModelObject() interface{} {
	return &ResourceAnchorResourceModel{}
}

func (ResourceAnchorResource) ResourceType() string {
	return "azurerm_oracle_resource_anchor"
}

func (r ResourceAnchorResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.ResourceAnchors
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ResourceAnchorResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := resourceanchors.NewResourceAnchorID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := resourceanchors.ResourceAnchor{
				Location:   model.Location,
				Tags:       pointer.To(model.Tags),
				Properties: &resourceanchors.ResourceAnchorProperties{},
			}
			if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (ResourceAnchorResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.ResourceAnchors
			id, err := resourceanchors.ParseResourceAnchorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := ResourceAnchorResourceModel{
				Name:              id.ResourceAnchorName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = model.Location
				state.Tags = pointer.From(model.Tags)
				if props := model.Properties; props != nil {
					state.LinkedCompartmentID = pointer.From(props.LinkedCompartmentId)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (ResourceAnchorResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.ResourceAnchors
			id, err := resourceanchors.ParseResourceAnchorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ResourceAnchorResourceModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			update := resourceanchors.ResourceAnchorUpdate{}

			if metadata.ResourceData.HasChange("tags") {
				update.Tags = pointer.To(model.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, update); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (ResourceAnchorResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.ResourceAnchors

			id, err := resourceanchors.ParseResourceAnchorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (ResourceAnchorResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return resourceanchors.ValidateResourceAnchorID
}
