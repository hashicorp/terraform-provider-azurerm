// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storagemover

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2023-03-01/agents"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2023-03-01/storagemovers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computevalidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type StorageMoverAgentResourceModel struct {
	Name                string `tfschema:"name"`
	StorageMoverId      string `tfschema:"storage_mover_id"`
	ArcVirtualMachineId string `tfschema:"arc_virtual_machine_id"`
	ArcVmUuid           string `tfschema:"arc_virtual_machine_uuid"`
	Description         string `tfschema:"description"`
}

type StorageMoverAgentResource struct{}

var _ sdk.ResourceWithUpdate = StorageMoverAgentResource{}

func (r StorageMoverAgentResource) ResourceType() string {
	return "azurerm_storage_mover_agent"
}

func (r StorageMoverAgentResource) ModelObject() interface{} {
	return &StorageMoverAgentResourceModel{}
}

func (r StorageMoverAgentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return agents.ValidateAgentID
}

func (r StorageMoverAgentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"arc_virtual_machine_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: computevalidate.HybridMachineID,
		},

		"arc_virtual_machine_uuid": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
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

func (r StorageMoverAgentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StorageMoverAgentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model StorageMoverAgentResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.StorageMover.AgentsClient
			storageMoverId, err := storagemovers.ParseStorageMoverID(model.StorageMoverId)
			if err != nil {
				return err
			}

			id := agents.NewAgentID(storageMoverId.SubscriptionId, storageMoverId.ResourceGroupName, storageMoverId.StorageMoverName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := agents.Agent{
				Properties: agents.AgentProperties{
					ArcResourceId: model.ArcVirtualMachineId,
					ArcVMUuid:     model.ArcVmUuid,
				},
			}

			if model.Description != "" {
				properties.Properties.Description = &model.Description
			}

			if _, err := client.CreateOrUpdate(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r StorageMoverAgentResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageMover.AgentsClient

			id, err := agents.ParseAgentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model StorageMoverAgentResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if metadata.ResourceData.HasChange("description") {
				properties.Properties.Description = &model.Description
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r StorageMoverAgentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageMover.AgentsClient

			id, err := agents.ParseAgentID(metadata.ResourceData.Id())
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

			state := StorageMoverAgentResourceModel{
				Name:           id.AgentName,
				StorageMoverId: storagemovers.NewStorageMoverID(id.SubscriptionId, id.ResourceGroupName, id.StorageMoverName).ID(),
			}

			if model := resp.Model; model != nil {
				state.ArcVmUuid = model.Properties.ArcVMUuid
				state.ArcVirtualMachineId = model.Properties.ArcResourceId

				state.Description = pointer.From(model.Properties.Description)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r StorageMoverAgentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageMover.AgentsClient

			id, err := agents.ParseAgentID(metadata.ResourceData.Id())
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
