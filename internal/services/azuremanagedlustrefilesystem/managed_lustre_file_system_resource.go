package azuremanagedlustrefilesystem

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-05-01/amlfilesystems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/azuremanagedlustrefilesystem/validate"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagedLustreFileSystemModel struct {
	Name                string                       `tfschema:"name"`
	ResourceGroupName   string                       `tfschema:"resource_group_name"`
	Location            string                       `tfschema:"location"`
	HsmSetting          []HsmSetting                 `tfschema:"hsm_setting"`
	Identity            []identity.ModelUserAssigned `tfschema:"identity"`
	EncryptionKey       []EncryptionKey              `tfschema:"encryption_key"`
	MaintenanceWindow   []MaintenanceWindow          `tfschema:"maintenance_window"`
	SkuName             string                       `tfschema:"sku_name"`
	StorageCapacityInTb int64                        `tfschema:"storage_capacity_in_tb"`
	SubnetId            string                       `tfschema:"subnet_id"`
	Zones               []string                     `tfschema:"zones"`
	Tags                map[string]string            `tfschema:"tags"`
}

type HsmSetting struct {
	ContainerId        string `tfschema:"container_id"`
	LoggingContainerId string `tfschema:"logging_container_id"`
	ImportPrefix       string `tfschema:"import_prefix"`
}

type EncryptionKey struct {
	KeyUrl        string `tfschema:"key_url"`
	SourceVaultId string `tfschema:"source_vault_id"`
}

type MaintenanceWindow struct {
	DayOfWeek      amlfilesystems.MaintenanceDayOfWeekType `tfschema:"day_of_week"`
	TimeOfDayInUTC string                                  `tfschema:"time_of_day_in_utc"`
}

type ManagedLustreFileSystemResource struct{}

var _ sdk.ResourceWithUpdate = ManagedLustreFileSystemResource{}

var _ sdk.ResourceWithCustomizeDiff = ManagedLustreFileSystemResource{}

func (r ManagedLustreFileSystemResource) ResourceType() string {
	return "azurerm_managed_lustre_file_system"
}

func (r ManagedLustreFileSystemResource) ModelObject() interface{} {
	return &ManagedLustreFileSystemModel{}
}

func (r ManagedLustreFileSystemResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return amlfilesystems.ValidateAmlFilesystemID
}

func (r ManagedLustreFileSystemResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagedLustreFileSystemName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"maintenance_window": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"day_of_week": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(amlfilesystems.PossibleValuesForMaintenanceDayOfWeekType(), false),
					},

					"time_of_day_in_utc": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.TimeOfDayInUTC,
					},
				},
			},
		},

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"AMLFS-Durable-Premium-40",
				"AMLFS-Durable-Premium-125",
				"AMLFS-Durable-Premium-250",
				"AMLFS-Durable-Premium-500",
			}, false),
		},

		"storage_capacity_in_tb": {
			Type:     pluginsdk.TypeInt,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.IntBetween(8, 128),
				validation.IntDivisibleBy(8),
			),
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"zones": commonschema.ZonesMultipleRequiredForceNew(),

		"hsm_setting": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"container_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateStorageContainerID,
					},

					"logging_container_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateStorageContainerID,
					},

					"import_prefix": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validate.ImportPrefix,
					},
				},
			},
		},

		"identity": commonschema.UserAssignedIdentityOptionalForceNew(),

		"encryption_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_url": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyVaultValidate.NestedItemId,
					},

					"source_vault_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateKeyVaultID,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r ManagedLustreFileSystemResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagedLustreFileSystemResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if oldVal, newVal := metadata.ResourceDiff.GetChange("encryption_key"); len(oldVal.([]interface{})) > 0 && len(newVal.([]interface{})) == 0 {
				if err := metadata.ResourceDiff.ForceNew("encryption_key"); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func (r ManagedLustreFileSystemResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagedLustreFileSystemModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AzureManagedLustreFileSystem.AmlFilesystems
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := amlfilesystems.NewAmlFilesystemID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			identity, err := identity.ExpandUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return err
			}

			properties := &amlfilesystems.AmlFilesystem{
				Location: location.Normalize(model.Location),
				Identity: identity,
				Properties: &amlfilesystems.AmlFilesystemProperties{
					Hsm:                expandManagedLustreFileSystemHsmSetting(model.HsmSetting),
					EncryptionSettings: expandManagedLustreFileSystemEncryptionKey(model.EncryptionKey),
					MaintenanceWindow:  expandManagedLustreFileSystemMaintenanceWindowForCreate(model.MaintenanceWindow),
					FilesystemSubnet:   model.SubnetId,
					StorageCapacityTiB: float64(model.StorageCapacityInTb),
				},
				Sku: &amlfilesystems.SkuName{
					Name: pointer.To(model.SkuName),
				},
				Zones: pointer.To(model.Zones),
				Tags:  pointer.To(model.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagedLustreFileSystemResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureManagedLustreFileSystem.AmlFilesystems

			id, err := amlfilesystems.ParseAmlFilesystemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagedLustreFileSystemModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			properties := amlfilesystems.AmlFilesystemUpdate{
				Properties: &amlfilesystems.AmlFilesystemUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("maintenance_window") {
				properties.Properties.MaintenanceWindow = expandManagedLustreFileSystemMaintenanceWindowForUpdate(model.MaintenanceWindow)
			}

			if metadata.ResourceData.HasChange("encryption_key") {
				properties.Properties.EncryptionSettings = expandManagedLustreFileSystemEncryptionKey(model.EncryptionKey)
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = pointer.To(model.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ManagedLustreFileSystemResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureManagedLustreFileSystem.AmlFilesystems

			id, err := amlfilesystems.ParseAmlFilesystemID(metadata.ResourceData.Id())
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

			state := ManagedLustreFileSystemModel{}
			if model := resp.Model; model != nil {
				state.Name = id.AmlFilesystemName
				state.ResourceGroupName = id.ResourceGroupName
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				identity, err := identity.FlattenUserAssignedMapToModel(model.Identity)
				if err != nil {
					return err
				}
				state.Identity = pointer.From(identity)

				if properties := model.Properties; properties != nil {
					state.SubnetId = properties.FilesystemSubnet
					state.StorageCapacityInTb = int64(properties.StorageCapacityTiB)
					state.MaintenanceWindow = flattenManagedLustreFileSystemMaintenanceWindow(properties.MaintenanceWindow)
					state.HsmSetting = flattenManagedLustreFileSystemHsmSetting(properties.Hsm)
					state.Zones = pointer.From(model.Zones)
					state.EncryptionKey = flattenManagedLustreFileSystemEncryptionKey(properties.EncryptionSettings)

					if v := model.Sku; v != nil {
						state.SkuName = pointer.From(v.Name)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagedLustreFileSystemResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureManagedLustreFileSystem.AmlFilesystems

			id, err := amlfilesystems.ParseAmlFilesystemID(metadata.ResourceData.Id())
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

func expandManagedLustreFileSystemMaintenanceWindowForCreate(input []MaintenanceWindow) amlfilesystems.AmlFilesystemPropertiesMaintenanceWindow {
	maintenanceWindow := &input[0]

	return amlfilesystems.AmlFilesystemPropertiesMaintenanceWindow{
		DayOfWeek:    pointer.To(maintenanceWindow.DayOfWeek),
		TimeOfDayUTC: pointer.To(maintenanceWindow.TimeOfDayInUTC),
	}
}

func expandManagedLustreFileSystemMaintenanceWindowForUpdate(input []MaintenanceWindow) *amlfilesystems.AmlFilesystemUpdatePropertiesMaintenanceWindow {
	if len(input) == 0 {
		return nil
	}

	maintenanceWindow := &input[0]

	return &amlfilesystems.AmlFilesystemUpdatePropertiesMaintenanceWindow{
		DayOfWeek:    pointer.To(maintenanceWindow.DayOfWeek),
		TimeOfDayUTC: pointer.To(maintenanceWindow.TimeOfDayInUTC),
	}
}

func flattenManagedLustreFileSystemMaintenanceWindow(input amlfilesystems.AmlFilesystemPropertiesMaintenanceWindow) []MaintenanceWindow {
	var result []MaintenanceWindow

	maintenanceWindow := MaintenanceWindow{
		DayOfWeek:      pointer.From(input.DayOfWeek),
		TimeOfDayInUTC: pointer.From(input.TimeOfDayUTC),
	}

	return append(result, maintenanceWindow)
}

func expandManagedLustreFileSystemEncryptionKey(input []EncryptionKey) *amlfilesystems.AmlFilesystemEncryptionSettings {
	if len(input) == 0 {
		return nil
	}

	encryptionKey := &input[0]

	result := &amlfilesystems.KeyVaultKeyReference{
		KeyUrl: encryptionKey.KeyUrl,
		SourceVault: amlfilesystems.KeyVaultKeyReferenceSourceVault{
			Id: pointer.To(encryptionKey.SourceVaultId),
		},
	}

	return &amlfilesystems.AmlFilesystemEncryptionSettings{
		KeyEncryptionKey: result,
	}
}

func flattenManagedLustreFileSystemEncryptionKey(input *amlfilesystems.AmlFilesystemEncryptionSettings) []EncryptionKey {
	result := make([]EncryptionKey, 0)
	if input == nil || input.KeyEncryptionKey == nil {
		return result
	}

	encryptionKey := EncryptionKey{
		KeyUrl:        input.KeyEncryptionKey.KeyUrl,
		SourceVaultId: pointer.From(input.KeyEncryptionKey.SourceVault.Id),
	}

	return append(result, encryptionKey)
}

func expandManagedLustreFileSystemHsmSetting(input []HsmSetting) *amlfilesystems.AmlFilesystemPropertiesHsm {
	if len(input) == 0 {
		return nil
	}

	hsmSetting := &input[0]

	result := &amlfilesystems.AmlFilesystemHsmSettings{
		Container:        hsmSetting.ContainerId,
		LoggingContainer: hsmSetting.LoggingContainerId,
	}

	if hsmSetting.ImportPrefix != "" {
		result.ImportPrefix = pointer.To(hsmSetting.ImportPrefix)
	}

	return &amlfilesystems.AmlFilesystemPropertiesHsm{
		Settings: result,
	}
}

func flattenManagedLustreFileSystemHsmSetting(input *amlfilesystems.AmlFilesystemPropertiesHsm) []HsmSetting {
	result := make([]HsmSetting, 0)
	if input == nil || input.Settings == nil {
		return result
	}

	hsmSetting := HsmSetting{}

	if v := input.Settings; v != nil {
		hsmSetting.ContainerId = v.Container
		hsmSetting.LoggingContainerId = v.LoggingContainer
		hsmSetting.ImportPrefix = pointer.From(v.ImportPrefix)
	}

	return append(result, hsmSetting)
}
