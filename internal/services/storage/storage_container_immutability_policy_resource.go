// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/blobcontainers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type StorageContainerImmutabilityPolicyResource struct{}

var _ sdk.ResourceWithUpdate = StorageContainerImmutabilityPolicyResource{}

type ContainerImmutabilityPolicyModel struct {
	StorageContainerResourceManagerId string `tfschema:"storage_container_resource_manager_id"`
	ImmutabilityPeriodInDays          int64  `tfschema:"immutability_period_in_days"`
	Locked                            bool   `tfschema:"locked"`
	ProtectedAppendWritesAllEnabled   bool   `tfschema:"protected_append_writes_all_enabled"`
	ProtectedAppendWritesEnabled      bool   `tfschema:"protected_append_writes_enabled"`
}

func (r StorageContainerImmutabilityPolicyResource) ResourceType() string {
	return "azurerm_storage_container_immutability_policy"
}

func (r StorageContainerImmutabilityPolicyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.StorageContainerImmutabilityPolicyID
}

func (r StorageContainerImmutabilityPolicyResource) ModelObject() interface{} {
	return &ContainerImmutabilityPolicyModel{}
}

func (r StorageContainerImmutabilityPolicyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"storage_container_resource_manager_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageContainerID,
		},

		"immutability_period_in_days": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 146000),
		},

		"locked": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"protected_append_writes_all_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"protected_append_writes_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func (r StorageContainerImmutabilityPolicyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StorageContainerImmutabilityPolicyResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			diff := metadata.ResourceDiff

			protectedAppendWritesAllEnabled := diff.Get("protected_append_writes_all_enabled").(bool)
			protectedAppendWritesEnabled := diff.Get("protected_append_writes_enabled").(bool)

			if protectedAppendWritesAllEnabled && protectedAppendWritesEnabled {
				return fmt.Errorf("`protected_append_writes_all_enabled` and `protected_append_writes_enabled` cannot be set at the same time")
			}

			lockedOld, lockedNew := diff.GetChange("locked")

			if lockedOld.(bool) && !lockedNew.(bool) {
				return fmt.Errorf("unable to set `locked = false` - once an immutability policy locked it cannot be unlocked")
			}

			if lockedOld.(bool) {
				if diff.HasChange("immutability_period_in_days") {
					if periodOld, periodNew := diff.GetChange("immutability_period_in_days"); periodOld.(int) < periodNew.(int) {
						return fmt.Errorf("`immutability_period_in_days` cannot be decreased once an immutability policy has been locked")
					}
				}

				if diff.HasChange("protected_append_writes_all_enabled") {
					return fmt.Errorf("`protected_append_writes_all_enabled` cannot be changed once an immutability policy has been locked")
				}

				if diff.HasChange("protected_append_writes_enabled") {
					return fmt.Errorf("`protected_append_writes_enabled` cannot be changed once an immutability policy has been locked")
				}
			}

			return nil
		},
	}
}

func (r StorageContainerImmutabilityPolicyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.ResourceManager.BlobContainers

			var model ContainerImmutabilityPolicyModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			containerId, err := commonids.ParseStorageContainerID(model.StorageContainerResourceManagerId)
			if err != nil {
				return err
			}

			id := parse.NewStorageContainerImmutabilityPolicyID(containerId.SubscriptionId, containerId.ResourceGroupName, containerId.StorageAccountName, "default", containerId.ContainerName, "default")

			existing, err := client.GetImmutabilityPolicy(ctx, *containerId, blobcontainers.DefaultGetImmutabilityPolicyOperationOptions())
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) && !r.isDeleted(existing.Model) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			input := blobcontainers.ImmutabilityPolicy{
				Properties: blobcontainers.ImmutabilityPolicyProperty{
					AllowProtectedAppendWrites:            pointer.To(model.ProtectedAppendWritesEnabled),
					AllowProtectedAppendWritesAll:         pointer.To(model.ProtectedAppendWritesAllEnabled),
					ImmutabilityPeriodSinceCreationInDays: pointer.To(model.ImmutabilityPeriodInDays),
				},
			}

			resp, err := client.CreateOrUpdateImmutabilityPolicy(ctx, *containerId, input, blobcontainers.DefaultCreateOrUpdateImmutabilityPolicyOperationOptions())
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			// Lock the policy if requested - note that this is a one-way operation that prevents subsequent changes or
			// deletion to the policy, the container it applies to, and the storage account where it resides.
			if model.Locked {
				if resp.Model == nil {
					return fmt.Errorf("preparing to lock %s: model was nil", id)
				}

				options := blobcontainers.LockImmutabilityPolicyOperationOptions{
					IfMatch: resp.Model.Etag,
				}

				if _, err = client.LockImmutabilityPolicy(ctx, *containerId, options); err != nil {
					return fmt.Errorf("locking %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func (r StorageContainerImmutabilityPolicyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.ResourceManager.BlobContainers

			id, err := parse.StorageContainerImmutabilityPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ContainerImmutabilityPolicyModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			containerId, err := commonids.ParseStorageContainerID(model.StorageContainerResourceManagerId)
			if err != nil {
				return err
			}

			resp, err := client.GetImmutabilityPolicy(ctx, *containerId, blobcontainers.DefaultGetImmutabilityPolicyOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) || r.isDeleted(resp.Model) {
					return nil
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			input := blobcontainers.ImmutabilityPolicy{
				Properties: blobcontainers.ImmutabilityPolicyProperty{
					AllowProtectedAppendWrites:            pointer.To(model.ProtectedAppendWritesEnabled),
					AllowProtectedAppendWritesAll:         pointer.To(model.ProtectedAppendWritesAllEnabled),
					ImmutabilityPeriodSinceCreationInDays: pointer.To(model.ImmutabilityPeriodInDays),
				},
			}

			options := blobcontainers.CreateOrUpdateImmutabilityPolicyOperationOptions{
				IfMatch: resp.Model.Etag,
			}

			updateResp, err := client.CreateOrUpdateImmutabilityPolicy(ctx, *containerId, input, options)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			// Lock the policy if requested - note that this is a one-way operation that prevents subsequent changes or
			// deletion to the policy, the container it applies to, and the storage account where it resides.
			if model.Locked {
				if updateResp.Model == nil {
					return fmt.Errorf("preparing to lock %s: model was nil", id)
				}

				lockOptions := blobcontainers.LockImmutabilityPolicyOperationOptions{
					IfMatch: updateResp.Model.Etag,
				}

				if _, err = client.LockImmutabilityPolicy(ctx, *containerId, lockOptions); err != nil {
					return fmt.Errorf("locking %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func (r StorageContainerImmutabilityPolicyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.ResourceManager.BlobContainers

			id, err := parse.StorageContainerImmutabilityPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			containerId := commonids.NewStorageContainerID(id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.ContainerName)

			resp, err := client.GetImmutabilityPolicy(ctx, containerId, blobcontainers.DefaultGetImmutabilityPolicyOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) || r.isDeleted(resp.Model) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := ContainerImmutabilityPolicyModel{
				StorageContainerResourceManagerId: containerId.ID(),
			}

			if resp.Model != nil {
				props := resp.Model.Properties
				if props.AllowProtectedAppendWrites != nil {
					state.ProtectedAppendWritesEnabled = *props.AllowProtectedAppendWrites
				}
				if props.AllowProtectedAppendWritesAll != nil {
					state.ProtectedAppendWritesAllEnabled = *props.AllowProtectedAppendWritesAll
				}
				if props.ImmutabilityPeriodSinceCreationInDays != nil {
					state.ImmutabilityPeriodInDays = *props.ImmutabilityPeriodSinceCreationInDays
				}
				if props.State != nil {
					state.Locked = *props.State == blobcontainers.ImmutabilityPolicyStateLocked
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r StorageContainerImmutabilityPolicyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.ResourceManager.BlobContainers

			id, err := parse.StorageContainerImmutabilityPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			containerId := commonids.NewStorageContainerID(id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.ContainerName)

			resp, err := client.GetImmutabilityPolicy(ctx, containerId, blobcontainers.DefaultGetImmutabilityPolicyOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) || r.isDeleted(resp.Model) {
					return nil
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			options := blobcontainers.DeleteImmutabilityPolicyOperationOptions{
				IfMatch: resp.Model.Etag,
			}

			if _, err := client.DeleteImmutabilityPolicy(ctx, containerId, options); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r StorageContainerImmutabilityPolicyResource) isDeleted(input *blobcontainers.ImmutabilityPolicy) bool {
	if input == nil {
		return true
	}
	if input.Properties.State != nil && strings.EqualFold(string(*input.Properties.State), "Deleted") {
		return true
	}
	return false
}
