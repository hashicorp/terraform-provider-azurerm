// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storagemover

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2023-03-01/projects"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2023-03-01/storagemovers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type StorageMoverProjectResourceModel struct {
	Name           string `tfschema:"name"`
	StorageMoverId string `tfschema:"storage_mover_id"`
	Description    string `tfschema:"description"`
}

type StorageMoverProjectResource struct{}

var _ sdk.ResourceWithUpdate = StorageMoverProjectResource{}

func (r StorageMoverProjectResource) ResourceType() string {
	return "azurerm_storage_mover_project"
}

func (r StorageMoverProjectResource) ModelObject() interface{} {
	return &StorageMoverProjectResourceModel{}
}

func (r StorageMoverProjectResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return projects.ValidateProjectID
}

func (r StorageMoverProjectResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[0-9a-zA-Z][-_0-9a-zA-Z]{0,63}$`),
				`The name must be between 1 and 64 characters in length, begin with a letter or number, and may contain letters, numbers, dashes and underscore.`,
			),
		},

		"storage_mover_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: storagemovers.ValidateStorageMoverID,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r StorageMoverProjectResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StorageMoverProjectResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model StorageMoverProjectResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.StorageMover.ProjectsClient
			storageMoverId, err := storagemovers.ParseStorageMoverID(model.StorageMoverId)
			if err != nil {
				return err
			}

			id := projects.NewProjectID(storageMoverId.SubscriptionId, storageMoverId.ResourceGroupName, storageMoverId.StorageMoverName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &projects.Project{
				Properties: &projects.ProjectProperties{},
			}

			if model.Description != "" {
				properties.Properties.Description = &model.Description
			}

			if _, err := client.CreateOrUpdate(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r StorageMoverProjectResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageMover.ProjectsClient

			id, err := projects.ParseProjectID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model StorageMoverProjectResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			if metadata.ResourceData.HasChange("description") && properties.Properties != nil {
				properties.Properties.Description = pointer.To(model.Description)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r StorageMoverProjectResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageMover.ProjectsClient

			id, err := projects.ParseProjectID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := StorageMoverProjectResourceModel{
				Name:           id.ProjectName,
				StorageMoverId: storagemovers.NewStorageMoverID(id.SubscriptionId, id.ResourceGroupName, id.StorageMoverName).ID(),
			}

			if resp.Model != nil {
				if properties := resp.Model.Properties; properties != nil {
					state.Description = pointer.From(properties.Description)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r StorageMoverProjectResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageMover.ProjectsClient

			id, err := projects.ParseProjectID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
