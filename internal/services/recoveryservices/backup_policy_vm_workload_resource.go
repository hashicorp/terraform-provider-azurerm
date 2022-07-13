package recoveryservices

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2021-12-01/backup"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type BackupProtectionPolicyVMWorkloadModel struct {
	Name              string     `tfschema:"name"`
	ResourceGroupName string     `tfschema:"resource_group_name"`
	RecoveryVaultName string     `tfschema:"recovery_vault_name"`
	Settings          []Settings `tfschema:"settings"`
	WorkloadType      string     `tfschema:"workload_type"`
}

type Settings struct {
	CompressionEnabled    *bool   `tfschema:"compression_enabled"`
	SqlCompressionEnabled *bool   `tfschema:"sql_compression_enabled"`
	TimeZone              *string `tfschema:"time_zone"`
}

type BackupProtectionPolicyVMWorkloadResource struct{}

var _ sdk.ResourceWithUpdate = BackupProtectionPolicyVMWorkloadResource{}

func (r BackupProtectionPolicyVMWorkloadResource) ResourceType() string {
	return "azurerm_backup_policy_vm_workload"
}

func (r BackupProtectionPolicyVMWorkloadResource) ModelObject() interface{} {
	return &BackupProtectionPolicyVMWorkloadModel{}
}

func (r BackupProtectionPolicyVMWorkloadResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.BackupPolicyID
}

func (r BackupProtectionPolicyVMWorkloadResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.BackupPolicyName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"recovery_vault_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RecoveryServicesVaultName,
		},

		"workload_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(backup.WorkloadTypeSQLDataBase),
				string(backup.WorkloadTypeSAPHanaDatabase),
			}, false),
		},

		"settings": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"compression_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"sql_compression_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"time_zone": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}

func (r BackupProtectionPolicyVMWorkloadResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r BackupProtectionPolicyVMWorkloadResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model BackupProtectionPolicyVMWorkloadModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.RecoveryServices.ProtectionPoliciesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := parse.NewBackupPolicyID(subscriptionId, model.ResourceGroupName, model.RecoveryVaultName, model.Name)

			existing, err := client.Get(ctx, id.VaultName, id.ResourceGroup, id.Name)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return tf.ImportAsExistsError("azurerm_backup_policy_vm_workload", id.ID())
			}

			properties := &backup.ProtectionPolicyResource{
				Properties: &backup.AzureVMWorkloadProtectionPolicy{
					BackupManagementType: backup.ManagementTypeBasicProtectionPolicyBackupManagementTypeAzureWorkload,
					Settings:             expandBackupProtectionPolicyVMWorkloadSettings(model.Settings),
					WorkLoadType:         backup.WorkloadType(model.WorkloadType),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id.VaultName, id.ResourceGroup, id.Name, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r BackupProtectionPolicyVMWorkloadResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ProtectionPoliciesClient

			id, err := parse.BackupPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model BackupProtectionPolicyVMWorkloadModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, id.VaultName, id.ResourceGroup, id.Name)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if props := existing.Properties; props != nil {
				vmWorkload, _ := props.AsAzureVMWorkloadProtectionPolicy()

				if metadata.ResourceData.HasChange("settings") {
					vmWorkload.Settings = expandBackupProtectionPolicyVMWorkloadSettings(model.Settings)
				}
			}

			if _, err := client.CreateOrUpdate(ctx, id.VaultName, id.ResourceGroup, id.Name, existing); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r BackupProtectionPolicyVMWorkloadResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ProtectionPoliciesClient

			id, err := parse.BackupPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.VaultName, id.ResourceGroup, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := BackupProtectionPolicyVMWorkloadModel{
				Name:              id.Name,
				ResourceGroupName: id.ResourceGroup,
				RecoveryVaultName: id.VaultName,
			}

			if props := resp.Properties; props != nil {
				vmWorkload, _ := props.AsAzureVMWorkloadProtectionPolicy()
				state.WorkloadType = string(vmWorkload.WorkLoadType)
				state.Settings = flattenBackupProtectionPolicyVMWorkloadSettings(vmWorkload.Settings)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r BackupProtectionPolicyVMWorkloadResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ProtectionPoliciesClient

			id, err := parse.BackupPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.VaultName, id.ResourceGroup, id.Name)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandBackupProtectionPolicyVMWorkloadSettings(input []Settings) *backup.Settings {
	if len(input) == 0 {
		return &backup.Settings{}
	}

	result := &backup.Settings{}
	settings := input[0]

	if settings.CompressionEnabled != nil {
		result.IsCompression = settings.CompressionEnabled
	}

	if settings.SqlCompressionEnabled != nil {
		result.Issqlcompression = settings.SqlCompressionEnabled
	}

	if settings.TimeZone != nil {
		result.TimeZone = settings.TimeZone
	}

	return result
}

func flattenBackupProtectionPolicyVMWorkloadSettings(input *backup.Settings) []Settings {
	if input == nil {
		return make([]Settings, 0)
	}

	result := make([]Settings, 0)

	result = append(result, Settings{
		CompressionEnabled:    input.IsCompression,
		SqlCompressionEnabled: input.Issqlcompression,
		TimeZone:              input.TimeZone,
	})

	return result
}
