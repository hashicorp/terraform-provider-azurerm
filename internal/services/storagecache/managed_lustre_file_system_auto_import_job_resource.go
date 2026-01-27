// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storagecache

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2025-07-01/autoimportjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagedLustreFileSystemAutoImportJobModel struct {
	Name                      string            `tfschema:"name"`
	ManagedLustreFileSystemId string            `tfschema:"managed_lustre_file_system_id"`
	Location                  string            `tfschema:"location"`
	AutoImportPrefixes        []string          `tfschema:"auto_import_prefixes"`
	ConflictResolutionMode    string            `tfschema:"conflict_resolution_mode"`
	DeletionsEnabled          bool              `tfschema:"deletions_enabled"`
	MaximumErrors             int64             `tfschema:"maximum_errors"`
	AdminStatusEnabled        bool              `tfschema:"admin_status_enabled"`
	Tags                      map[string]string `tfschema:"tags"`
}

type ManagedLustreFileSystemAutoImportJobResource struct{}

var _ sdk.ResourceWithUpdate = ManagedLustreFileSystemAutoImportJobResource{}

func (r ManagedLustreFileSystemAutoImportJobResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagedLustreFileSystemAutoImportJobResource) ModelObject() interface{} {
	return &ManagedLustreFileSystemAutoImportJobModel{}
}

func (r ManagedLustreFileSystemAutoImportJobResource) ResourceType() string {
	return "azurerm_managed_lustre_file_system_auto_import_job"
}

func (r ManagedLustreFileSystemAutoImportJobResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z][-0-9a-zA-Z_]{0,78}[0-9a-zA-Z]$`), "The `name` must be 2-80 characters long and can only contain alphanumeric characters, underscores, and hyphens, and must start and end with an alphanumeric character"),
		},

		"managed_lustre_file_system_id": commonschema.ResourceIDReferenceRequiredForceNew(&autoimportjobs.AmlFilesystemId{}),

		"location": commonschema.Location(),

		"auto_import_prefixes": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"admin_status_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"conflict_resolution_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(autoimportjobs.ConflictResolutionModeFail),
			ValidateFunc: validation.StringInSlice([]string{
				string(autoimportjobs.ConflictResolutionModeFail),
				string(autoimportjobs.ConflictResolutionModeOverwriteAlways),
				string(autoimportjobs.ConflictResolutionModeOverwriteIfDirty),
				string(autoimportjobs.ConflictResolutionModeSkip),
			}, false),
		},

		"deletions_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"maximum_errors": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      0,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"tags": commonschema.Tags(),
	}
}

func (r ManagedLustreFileSystemAutoImportJobResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagedLustreFileSystemAutoImportJobModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding model: %w", err)
			}

			client := metadata.Client.StorageCache.AutoImportJobs

			amlFileSystemId, err := autoimportjobs.ParseAmlFilesystemID(model.ManagedLustreFileSystemId)
			if err != nil {
				return err
			}

			id := autoimportjobs.NewAutoImportJobID(amlFileSystemId.SubscriptionId, amlFileSystemId.ResourceGroupName, amlFileSystemId.AmlFilesystemName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			props := autoimportjobs.AutoImportJob{
				Location: location.Normalize(model.Location),
				Properties: pointer.To(autoimportjobs.AutoImportJobProperties{
					AutoImportPrefixes:     pointer.To(model.AutoImportPrefixes),
					ConflictResolutionMode: pointer.To(autoimportjobs.ConflictResolutionMode(model.ConflictResolutionMode)),
					EnableDeletions:        pointer.To(model.DeletionsEnabled),
					MaximumErrors:          pointer.To(model.MaximumErrors),
				}),
				Tags: pointer.To(model.Tags),
			}

			if model.AdminStatusEnabled {
				props.Properties.AdminStatus = pointer.To(autoimportjobs.AdminStatusEnable)
			} else {
				props.Properties.AdminStatus = pointer.To(autoimportjobs.AdminStatusDisable)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagedLustreFileSystemAutoImportJobResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageCache.AutoImportJobs

			id, err := autoimportjobs.ParseAutoImportJobID(metadata.ResourceData.Id())
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

			state := ManagedLustreFileSystemAutoImportJobModel{}
			if model := resp.Model; model != nil {
				state.Name = id.AutoImportJobName
				state.ManagedLustreFileSystemId = autoimportjobs.NewAmlFilesystemID(id.SubscriptionId, id.ResourceGroupName, id.AmlFilesystemName).ID()
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.AutoImportPrefixes = pointer.From(props.AutoImportPrefixes)
					state.AdminStatusEnabled = pointer.From(props.AdminStatus) == autoimportjobs.AdminStatusEnable
					state.ConflictResolutionMode = string(pointer.From(props.ConflictResolutionMode))
					state.DeletionsEnabled = pointer.From(props.EnableDeletions)
					state.MaximumErrors = pointer.From(props.MaximumErrors)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagedLustreFileSystemAutoImportJobResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageCache.AutoImportJobs
			id, err := autoimportjobs.ParseAutoImportJobID(metadata.ResourceData.Id())
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

func (r ManagedLustreFileSystemAutoImportJobResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autoimportjobs.ValidateAutoImportJobID
}

func (r ManagedLustreFileSystemAutoImportJobResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageCache.AutoImportJobs

			id, err := autoimportjobs.ParseAutoImportJobID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagedLustreFileSystemAutoImportJobModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %w", err)
			}

			props := autoimportjobs.AutoImportJobUpdate{
				Properties: pointer.To(autoimportjobs.AutoImportJobUpdateProperties{}),
			}

			if metadata.ResourceData.HasChange("admin_status_enabled") {
				if model.AdminStatusEnabled {
					props.Properties.AdminStatus = pointer.To(autoimportjobs.AdminStatusEnable)
				} else {
					props.Properties.AdminStatus = pointer.To(autoimportjobs.AdminStatusDisable)
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				props.Tags = pointer.To(model.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, props); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
