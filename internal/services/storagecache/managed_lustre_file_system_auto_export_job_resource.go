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
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/autoexportjob"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/autoexportjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagecache/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagedLustreFileSystemAutoExportJobModel struct {
	Name               string            `tfschema:"name"`
	ResourceGroupName  string            `tfschema:"resource_group_name"`
	AmlFileSystemName  string            `tfschema:"aml_file_system_name"`
	Location           string            `tfschema:"location"`
	AutoExportPrefixes []string          `tfschema:"auto_export_prefixes"`
	AdminStatusEnabled bool              `tfschema:"admin_status_enabled"`
	Tags               map[string]string `tfschema:"tags"`
}

type ManagedLustreFileSystemAutoExportJobResource struct{}

var _ sdk.ResourceWithUpdate = ManagedLustreFileSystemAutoExportJobResource{}

func (r ManagedLustreFileSystemAutoExportJobResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagedLustreFileSystemAutoExportJobResource) ModelObject() interface{} {
	return &ManagedLustreFileSystemAutoExportJobModel{}
}

func (r ManagedLustreFileSystemAutoExportJobResource) ResourceType() string {
	return "azurerm_managed_lustre_file_system_auto_export_job"
}

func (r ManagedLustreFileSystemAutoExportJobResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z][-0-9a-zA-Z_]{0,78}[0-9a-zA-Z]$`), "name must be 3-80 characters long and can only contain alphanumeric characters, underscores, and hyphens, and must start and end with an alphanumeric character"),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"aml_file_system_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagedLustreFileSystemName,
		},

		"auto_export_prefixes": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
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

		"tags": commonschema.Tags(),
	}
}

func (r ManagedLustreFileSystemAutoExportJobResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagedLustreFileSystemAutoExportJobModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding model: %w", err)
			}

			autoExportJobsClient := metadata.Client.StorageCache.AutoExportJobs
			autoExportJobClient := metadata.Client.StorageCache.AutoExportJob
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := autoexportjobs.NewAutoExportJobID(subscriptionId, model.ResourceGroupName, model.AmlFileSystemName, model.Name)

			existing, err := autoExportJobsClient.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing Auto Export Job %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			props := autoexportjob.AutoExportJob{
				Location: location.Normalize(model.Location),
				Properties: pointer.To(autoexportjob.AutoExportJobProperties{
					AutoExportPrefixes: pointer.To(model.AutoExportPrefixes),
				}),
				Tags: pointer.To(model.Tags),
			}

			if model.AdminStatusEnabled {
				props.Properties.AdminStatus = pointer.To(autoexportjob.AutoExportJobAdminStatusEnable)
			} else {
				props.Properties.AdminStatus = pointer.To(autoexportjob.AutoExportJobAdminStatusDisable)
			}

			autoExportJobId := autoexportjob.NewAutoExportJobID(subscriptionId, model.ResourceGroupName, model.AmlFileSystemName, model.Name)
			if err := autoExportJobClient.CreateOrUpdateThenPoll(ctx, autoExportJobId, props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(autoExportJobId)
			return nil
		},
	}
}

func (r ManagedLustreFileSystemAutoExportJobResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageCache.AutoExportJobs

			id, err := autoexportjobs.ParseAutoExportJobID(metadata.ResourceData.Id())
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

			state := ManagedLustreFileSystemAutoExportJobModel{}
			if model := resp.Model; model != nil {
				state.Name = id.AutoExportJobName
				state.ResourceGroupName = id.ResourceGroupName
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
				state.AmlFileSystemName = id.AmlFilesystemName

				if props := model.Properties; props != nil {
					state.AutoExportPrefixes = pointer.From(props.AutoExportPrefixes)
					if props.AdminStatus != nil && string(pointer.From(props.AdminStatus)) == string(autoexportjob.AutoExportJobAdminStatusEnable) {
						state.AdminStatusEnabled = true
					} else {
						state.AdminStatusEnabled = false
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagedLustreFileSystemAutoExportJobResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageCache.AutoExportJobs
			id, err := autoexportjobs.ParseAutoExportJobID(metadata.ResourceData.Id())
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

func (r ManagedLustreFileSystemAutoExportJobResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autoexportjob.ValidateAutoExportJobID
}

func (r ManagedLustreFileSystemAutoExportJobResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageCache.AutoExportJob

			id, err := autoexportjob.ParseAutoExportJobID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagedLustreFileSystemAutoExportJobModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %w", err)
			}

			props := autoexportjob.AutoExportJob{
				Properties: pointer.To(autoexportjob.AutoExportJobProperties{}),
			}

			if metadata.ResourceData.HasChange("auto_export_prefixes") {
				props.Properties.AutoExportPrefixes = pointer.To(model.AutoExportPrefixes)
			}

			if metadata.ResourceData.HasChange("admin_status_enabled") {
				if model.AdminStatusEnabled {
					props.Properties.AdminStatus = pointer.To(autoexportjob.AutoExportJobAdminStatusEnable)
				} else {
					props.Properties.AdminStatus = pointer.To(autoexportjob.AutoExportJobAdminStatusDisable)
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
