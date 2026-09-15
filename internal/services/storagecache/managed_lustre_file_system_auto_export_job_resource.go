package storagecache

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/autoexportjob"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/autoexportjobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/importjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagedLustreFileSystemAutoExportJobModel struct {
	Name                      string            `tfschema:"name"`
	ManagedLustreFileSystemId string            `tfschema:"managed_lustre_file_system_id"`
	Location                  string            `tfschema:"location"`
	AutoExportPrefixes        []string          `tfschema:"auto_export_prefixes"`
	AdminStatusEnabled        bool              `tfschema:"admin_status_enabled"`
	Tags                      map[string]string `tfschema:"tags"`
}

type ManagedLustreFileSystemAutoExportJobResource struct{}

var _ sdk.ResourceWithUpdate = ManagedLustreFileSystemAutoExportJobResource{}

func (r ManagedLustreFileSystemAutoExportJobResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagedLustreFileSystemAutoExportJobResource) ModelObject() any {
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
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z][-0-9a-zA-Z_]{0,78}[0-9a-zA-Z]$`), "name must be 2-80 characters long and can only contain alphanumeric characters, underscores, and hyphens, and must start and end with an alphanumeric character"),
		},

		"managed_lustre_file_system_id": commonschema.ResourceIDReferenceRequiredForceNew(&autoexportjob.AmlFilesystemId{}),

		"location": commonschema.Location(),

		"auto_export_prefixes": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MinItems: 1,
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
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagedLustreFileSystemAutoExportJobModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding model: %w", err)
			}

			autoExportJobsClient := metadata.Client.StorageCache.AutoExportJobs
			autoExportJobClient := metadata.Client.StorageCache.AutoExportJob

			amlFileSystemId, err := autoexportjob.ParseAmlFilesystemID(model.ManagedLustreFileSystemId)
			if err != nil {
				return err
			}

			id := autoexportjobs.NewAutoExportJobID(amlFileSystemId.SubscriptionId, amlFileSystemId.ResourceGroupName, amlFileSystemId.AmlFilesystemName, model.Name)

			existing, err := autoExportJobsClient.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
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

			// Azure API does not allow creating an auto export job with adminStatus set to 'Disable',
			// so always create with 'Enable' and then update to 'Disable' if needed.
			props.Properties.AdminStatus = pointer.To(autoexportjob.AutoExportJobAdminStatusEnable)

			autoExportJobId := autoexportjob.NewAutoExportJobID(amlFileSystemId.SubscriptionId, amlFileSystemId.ResourceGroupName, amlFileSystemId.AmlFilesystemName, model.Name)

			importJobsClient := metadata.Client.StorageCache.ImportJobs
			importJobsFsId := importjobs.NewAmlFilesystemID(amlFileSystemId.SubscriptionId, amlFileSystemId.ResourceGroupName, amlFileSystemId.AmlFilesystemName)
			if err := waitForImportJobsToComplete(ctx, importJobsClient, importJobsFsId); err != nil {
				return fmt.Errorf("waiting for import jobs to complete on %s: %+v", importJobsFsId, err)
			}

			if err := autoExportJobClient.CreateOrUpdateThenPoll(ctx, autoExportJobId, props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if !model.AdminStatusEnabled {
				disableProps := autoexportjob.AutoExportJobUpdate{
					Properties: pointer.To(autoexportjob.AutoExportJobUpdateProperties{
						AdminStatus: pointer.To(autoexportjob.AutoExportJobAdminStatusDisable),
					}),
				}
				if err := autoExportJobClient.UpdateThenPoll(ctx, autoExportJobId, disableProps); err != nil {
					return fmt.Errorf("disabling admin status for %s: %+v", id, err)
				}
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
				state.ManagedLustreFileSystemId = autoexportjob.NewAmlFilesystemID(id.SubscriptionId, id.ResourceGroupName, id.AmlFilesystemName).ID()
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.AutoExportPrefixes = pointer.From(props.AutoExportPrefixes)
					state.AdminStatusEnabled = pointer.From(props.AdminStatus) == autoexportjobs.AutoExportJobAdminStatusEnable
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagedLustreFileSystemAutoExportJobResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
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
		Timeout: 120 * time.Minute,
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

			props := autoexportjob.AutoExportJobUpdate{}

			if metadata.ResourceData.HasChange("admin_status_enabled") {
				if model.AdminStatusEnabled {
					props.Properties = pointer.To(autoexportjob.AutoExportJobUpdateProperties{
						AdminStatus: pointer.To(autoexportjob.AutoExportJobAdminStatusEnable),
					})
				} else {
					props.Properties = pointer.To(autoexportjob.AutoExportJobUpdateProperties{
						AdminStatus: pointer.To(autoexportjob.AutoExportJobAdminStatusDisable),
					})
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				props.Tags = pointer.To(model.Tags)
			}

			importJobsClient := metadata.Client.StorageCache.ImportJobs
			importJobsFsId := importjobs.NewAmlFilesystemID(id.SubscriptionId, id.ResourceGroupName, id.AmlFilesystemName)
			if err := waitForImportJobsToComplete(ctx, importJobsClient, importJobsFsId); err != nil {
				return fmt.Errorf("waiting for import jobs to complete on %s: %+v", importJobsFsId, err)
			}

			if err := client.UpdateThenPoll(ctx, *id, props); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func waitForImportJobsToComplete(ctx context.Context, client *importjobs.ImportJobsClient, fsId importjobs.AmlFilesystemId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(30 * time.Minute)
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"InProgress"},
		Target:     []string{"Complete"},
		MinTimeout: 30 * time.Second,
		Timeout:    time.Until(deadline),
		Refresh: func() (any, string, error) {
			resp, err := client.ListByAmlFilesystemComplete(ctx, fsId)
			if err != nil {
				return nil, "", fmt.Errorf("listing import jobs for %s: %+v", fsId, err)
			}

			for _, job := range resp.Items {
				if job.Properties != nil {
					if job.Properties.ProvisioningState != nil {
						provisioningState := *job.Properties.ProvisioningState
						if provisioningState == importjobs.ImportJobProvisioningStateTypeCreating ||
							provisioningState == importjobs.ImportJobProvisioningStateTypeUpdating ||
							provisioningState == importjobs.ImportJobProvisioningStateTypeDeleting {
							log.Printf("[DEBUG] Import job %q has provisioning state %q, waiting...", pointer.From(job.Name), string(provisioningState))
							return resp, "InProgress", nil
						}
					}

					if job.Properties.Status != nil && job.Properties.Status.State != nil {
						statusState := *job.Properties.Status.State
						if statusState == importjobs.ImportStatusTypeInProgress ||
							statusState == importjobs.ImportStatusTypeCancelling {
							log.Printf("[DEBUG] Import job %q has status state %q, waiting...", pointer.From(job.Name), string(statusState))
							return resp, "InProgress", nil
						}
					}
				}
			}

			return resp, "Complete", nil
		},
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return err
	}

	return nil
}
