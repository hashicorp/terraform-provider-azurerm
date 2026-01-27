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
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/autoimportjob"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/autoimportjobs"
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
	EnableDeletions           bool              `tfschema:"enable_deletions"`
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
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z][-0-9a-zA-Z_]{0,78}[0-9a-zA-Z]$`), "name must be 3-80 characters long and can only contain alphanumeric characters, underscores, and hyphens, and must start and end with an alphanumeric character"),
		},

		"managed_lustre_file_system_id": commonschema.ResourceIDReferenceRequiredForceNew(&autoimportjob.AmlFilesystemId{}),

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

		"conflict_resolution_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(autoimportjob.ConflictResolutionModeFail),
			ValidateFunc: validation.StringInSlice([]string{
				string(autoimportjob.ConflictResolutionModeFail),
				string(autoimportjob.ConflictResolutionModeOverwriteAlways),
				string(autoimportjob.ConflictResolutionModeOverwriteIfDirty),
				string(autoimportjob.ConflictResolutionModeSkip),
			}, false),
		},

		"enable_deletions": {
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

		"admin_status_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
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

			autoImportJobsClient := metadata.Client.StorageCache.AutoImportJobs
			autoImportJobClient := metadata.Client.StorageCache.AutoImportJob

			amlFileSystemId, err := autoimportjob.ParseAmlFilesystemID(model.ManagedLustreFileSystemId)
			if err != nil {
				return err
			}

			id := autoimportjobs.NewAutoImportJobID(amlFileSystemId.SubscriptionId, amlFileSystemId.ResourceGroupName, amlFileSystemId.AmlFilesystemName, model.Name)

			existing, err := autoImportJobsClient.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			conflictResolutionMode := autoimportjob.ConflictResolutionMode(model.ConflictResolutionMode)
			props := autoimportjob.AutoImportJob{
				Location: location.Normalize(model.Location),
				Properties: pointer.To(autoimportjob.AutoImportJobProperties{
					AutoImportPrefixes:     pointer.To(model.AutoImportPrefixes),
					ConflictResolutionMode: pointer.To(conflictResolutionMode),
					EnableDeletions:        pointer.To(model.EnableDeletions),
					MaximumErrors:          pointer.To(model.MaximumErrors),
				}),
				Tags: pointer.To(model.Tags),
			}

			if model.AdminStatusEnabled {
				props.Properties.AdminStatus = pointer.To(autoimportjob.AutoImportJobAdminStatusEnable)
			} else {
				props.Properties.AdminStatus = pointer.To(autoimportjob.AutoImportJobAdminStatusDisable)
			}

			autoImportJobId := autoimportjob.NewAutoImportJobID(amlFileSystemId.SubscriptionId, amlFileSystemId.ResourceGroupName, amlFileSystemId.AmlFilesystemName, model.Name)
			if err := autoImportJobClient.CreateOrUpdateThenPoll(ctx, autoImportJobId, props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(autoImportJobId)
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
				state.ManagedLustreFileSystemId = autoimportjob.NewAmlFilesystemID(id.SubscriptionId, id.ResourceGroupName, id.AmlFilesystemName).ID()
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.AutoImportPrefixes = pointer.From(props.AutoImportPrefixes)
					state.AdminStatusEnabled = pointer.From(props.AdminStatus) == autoimportjobs.AutoImportJobAdminStatusEnable
					state.ConflictResolutionMode = string(pointer.From(props.ConflictResolutionMode))
					state.EnableDeletions = pointer.From(props.EnableDeletions)
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
	return autoimportjob.ValidateAutoImportJobID
}

func (r ManagedLustreFileSystemAutoImportJobResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageCache.AutoImportJob

			id, err := autoimportjob.ParseAutoImportJobID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagedLustreFileSystemAutoImportJobModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %w", err)
			}

			props := autoimportjob.AutoImportJob{
				Properties: pointer.To(autoimportjob.AutoImportJobProperties{}),
			}

			if metadata.ResourceData.HasChange("admin_status_enabled") {
				if model.AdminStatusEnabled {
					props.Properties.AdminStatus = pointer.To(autoimportjob.AutoImportJobAdminStatusEnable)
				} else {
					props.Properties.AdminStatus = pointer.To(autoimportjob.AutoImportJobAdminStatusDisable)
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				props.Tags = pointer.To(model.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, props); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
