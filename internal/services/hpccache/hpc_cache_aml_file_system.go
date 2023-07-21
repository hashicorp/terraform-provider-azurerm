package hpccache

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hpccache/validate"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type HPCCacheAMLFileSystemModel struct {
	Name                string                       `tfschema:"name"`
	ResourceGroupName   string                       `tfschema:"resource_group_name"`
	Location            string                       `tfschema:"location"`
	HsmSetting          []HsmSetting                 `tfschema:"hsm_setting"`
	Identity            []identity.ModelUserAssigned `tfschema:"identity"`
	KeyEncryptionKey    []KeyEncryptionKey           `tfschema:"key_encryption_key"`
	MaintenanceWindow   []MaintenanceWindow          `tfschema:"maintenance_window"`
	SkuName             string                       `tfschema:"sku_name"`
	StorageCapacityInTb float64                      `tfschema:"storage_capacity_in_tb"`
	SubnetId            string                       `tfschema:"subnet_id"`
	Zones               []string                     `tfschema:"zones"`
	Tags                map[string]string            `tfschema:"tags"`
}

type HsmSetting struct {
	Container        string `tfschema:"container"`
	ImportPrefix     string `tfschema:"import_prefix"`
	LoggingContainer string `tfschema:"logging_container"`
}

type KeyEncryptionKey struct {
	KeyUrl        string `tfschema:"key_url"`
	SourceVaultId string `tfschema:"source_vault_id"`
}

type MaintenanceWindow struct {
	DayOfWeek      amlfilesystems.MaintenanceDayOfWeekType `tfschema:"day_of_week"`
	TimeOfDayInUTC string                                  `tfschema:"time_of_day_in_utc"`
}

type HPCCacheAMLFileSystemResource struct{}

var _ sdk.ResourceWithUpdate = HPCCacheAMLFileSystemResource{}

func (r HPCCacheAMLFileSystemResource) ResourceType() string {
	return "azurerm_hpc_cache_aml_file_system"
}

func (r HPCCacheAMLFileSystemResource) ModelObject() interface{} {
	return &HPCCacheAMLFileSystemModel{}
}

func (r HPCCacheAMLFileSystemResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return amlfilesystems.ValidateAmlFilesystemID
}

func (r HPCCacheAMLFileSystemResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AMLFileSystemName,
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
						Optional:     true,
						ValidateFunc: validation.StringInSlice(amlfilesystems.PossibleValuesForMaintenanceDayOfWeekType(), false),
					},

					"time_of_day_in_utc": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.TimeOfDayInUTC,
					},
				},
			},
		},

		"storage_capacity_in_tb": {
			Type:     pluginsdk.TypeFloat,
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

		"hsm_setting": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"container": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: storageValidate.StorageContainerResourceManagerID,
					},

					"logging_container": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: storageValidate.StorageContainerResourceManagerID,
					},

					"import_prefix": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						Default:      "/",
						ValidateFunc: validate.ImportPrefix,
					},
				},
			},
		},

		"identity": commonschema.UserAssignedIdentityOptionalForceNew(),

		"key_encryption_key": {
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

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"tags": commonschema.Tags(),
	}
}

func (r HPCCacheAMLFileSystemResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r HPCCacheAMLFileSystemResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model HPCCacheAMLFileSystemModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.HPCCache.AMLFileSystemsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := amlfilesystems.NewAmlFilesystemID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			identity, err := expandAMLFileSystemIdentity(model.Identity)
			if err != nil {
				return err
			}

			properties := &amlfilesystems.AmlFilesystem{
				Location: location.Normalize(model.Location),
				Identity: identity,
				Properties: &amlfilesystems.AmlFilesystemProperties{
					MaintenanceWindow:  expandAMLFileSystemMaintenanceWindowForCreate(model.MaintenanceWindow),
					FilesystemSubnet:   model.SubnetId,
					StorageCapacityTiB: model.StorageCapacityInTb,
				},
				Zones: pointer.To(model.Zones),
				Tags:  pointer.To(model.Tags),
			}

			if v := model.HsmSetting; v != nil {
				properties.Properties.Hsm = &amlfilesystems.AmlFilesystemPropertiesHsm{
					Settings: expandAMLFileSystemHsmSetting(model.HsmSetting),
				}
			}

			if v := model.KeyEncryptionKey; v != nil {
				properties.Properties.EncryptionSettings = &amlfilesystems.AmlFilesystemEncryptionSettings{
					KeyEncryptionKey: expandAMLFileSystemKeyEncryptionKey(v),
				}
			}

			if v := model.SkuName; v != "" {
				properties.Sku = &amlfilesystems.SkuName{
					Name: pointer.To(v),
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r HPCCacheAMLFileSystemResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HPCCache.AMLFileSystemsClient

			id, err := amlfilesystems.ParseAmlFilesystemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model HPCCacheAMLFileSystemModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			properties := amlfilesystems.AmlFilesystemUpdate{}

			if metadata.ResourceData.HasChange("maintenance_window") {
				properties.Properties.MaintenanceWindow = expandAMLFileSystemMaintenanceWindowForUpdate(model.MaintenanceWindow)
			}

			if metadata.ResourceData.HasChange("key_encryption_key") {
				properties.Properties.EncryptionSettings = &amlfilesystems.AmlFilesystemEncryptionSettings{
					KeyEncryptionKey: expandAMLFileSystemKeyEncryptionKey(model.KeyEncryptionKey),
				}
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

func (r HPCCacheAMLFileSystemResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HPCCache.AMLFileSystemsClient

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

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := HPCCacheAMLFileSystemModel{
				Name:              id.AmlFilesystemName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			identity, err := flattenAMLFileSystemIdentity(model.Identity)
			if err != nil {
				return err
			}
			state.Identity = identity

			if properties := model.Properties; properties != nil {
				state.SubnetId = properties.FilesystemSubnet
				state.StorageCapacityInTb = properties.StorageCapacityTiB
				state.MaintenanceWindow = flattenAMLFileSystemMaintenanceWindow(properties.MaintenanceWindow)
				state.HsmSetting = flattenAMLFileSystemHsmSetting(properties.Hsm)

				if v := properties.EncryptionSettings; v != nil && v.KeyEncryptionKey != nil {
					state.KeyEncryptionKey = flattenAMLFileSystemKeyEncryptionKey(v.KeyEncryptionKey)
				}
			}

			if v := model.Sku; v != nil {
				state.SkuName = pointer.From(v.Name)
			}

			if model.Zones != nil {
				state.Zones = pointer.From(model.Zones)
			}

			if model.Tags != nil {
				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r HPCCacheAMLFileSystemResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HPCCache.AMLFileSystemsClient

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

func expandAMLFileSystemIdentity(input []identity.ModelUserAssigned) (*amlfilesystems.AmlFilesystemIdentity, error) {
	identityValue, err := identity.ExpandUserAssignedMapFromModel(input)
	if err != nil {
		return nil, fmt.Errorf("expanding `identity`: %+v", err)
	}

	output := amlfilesystems.AmlFilesystemIdentity{
		Type: pointer.To(amlfilesystems.AmlFilesystemIdentityType(string(identityValue.Type))),
	}

	if identityValue.Type == identity.TypeUserAssigned {
		output.UserAssignedIdentities = pointer.To(make(map[string]amlfilesystems.UserAssignedIdentitiesProperties))
		for k := range identityValue.IdentityIds {
			(*output.UserAssignedIdentities)[k] = amlfilesystems.UserAssignedIdentitiesProperties{}
		}
	}

	return &output, nil
}

func flattenAMLFileSystemIdentity(input *amlfilesystems.AmlFilesystemIdentity) ([]identity.ModelUserAssigned, error) {
	if input == nil {
		return nil, nil
	}

	identityIds := make(map[string]identity.UserAssignedIdentityDetails, 0)
	for k, v := range *input.UserAssignedIdentities {
		identityIds[k] = identity.UserAssignedIdentityDetails{
			ClientId:    v.ClientId,
			PrincipalId: v.PrincipalId,
		}
	}

	identityValue := identity.UserAssignedMap{
		Type:        identity.Type(string(pointer.From(input.Type))),
		IdentityIds: identityIds,
	}

	output, err := identity.FlattenUserAssignedMapToModel(&identityValue)
	if err != nil {
		return nil, fmt.Errorf("expanding `identity`: %+v", err)
	}

	return *output, nil
}

func expandAMLFileSystemMaintenanceWindowForCreate(input []MaintenanceWindow) amlfilesystems.AmlFilesystemPropertiesMaintenanceWindow {
	result := amlfilesystems.AmlFilesystemPropertiesMaintenanceWindow{}
	maintenanceWindow := &input[0]

	if v := maintenanceWindow.DayOfWeek; v != "" {
		result.DayOfWeek = pointer.To(v)
	}

	if v := maintenanceWindow.TimeOfDayInUTC; v != "" {
		result.TimeOfDayUTC = pointer.To(v)
	}

	return result
}

func expandAMLFileSystemMaintenanceWindowForUpdate(input []MaintenanceWindow) *amlfilesystems.AmlFilesystemUpdatePropertiesMaintenanceWindow {
	if len(input) == 0 {
		return nil
	}

	maintenanceWindow := &input[0]
	result := amlfilesystems.AmlFilesystemUpdatePropertiesMaintenanceWindow{}

	if v := maintenanceWindow.DayOfWeek; v != "" {
		result.DayOfWeek = pointer.To(v)
	}

	if v := maintenanceWindow.TimeOfDayInUTC; v != "" {
		result.TimeOfDayUTC = pointer.To(v)
	}

	return &result
}

func flattenAMLFileSystemMaintenanceWindow(input amlfilesystems.AmlFilesystemPropertiesMaintenanceWindow) []MaintenanceWindow {
	var result []MaintenanceWindow
	maintenanceWindow := MaintenanceWindow{}

	if input.DayOfWeek != nil {
		maintenanceWindow.DayOfWeek = pointer.From(input.DayOfWeek)
	}

	if input.TimeOfDayUTC != nil {
		maintenanceWindow.TimeOfDayInUTC = pointer.From(input.TimeOfDayUTC)
	}

	return append(result, maintenanceWindow)
}

func expandAMLFileSystemKeyEncryptionKey(input []KeyEncryptionKey) *amlfilesystems.KeyVaultKeyReference {
	if len(input) == 0 {
		return nil
	}

	keyEncryptionKey := &input[0]

	output := amlfilesystems.KeyVaultKeyReference{
		KeyUrl: keyEncryptionKey.KeyUrl,
		SourceVault: amlfilesystems.KeyVaultKeyReferenceSourceVault{
			Id: pointer.To(keyEncryptionKey.SourceVaultId),
		},
	}

	return &output
}

func flattenAMLFileSystemKeyEncryptionKey(input *amlfilesystems.KeyVaultKeyReference) []KeyEncryptionKey {
	if input == nil {
		return nil
	}

	var result []KeyEncryptionKey

	keyEncryptionKey := KeyEncryptionKey{
		KeyUrl:        input.KeyUrl,
		SourceVaultId: pointer.From(input.SourceVault.Id),
	}

	return append(result, keyEncryptionKey)
}

func expandAMLFileSystemHsmSetting(input []HsmSetting) *amlfilesystems.AmlFilesystemHsmSettings {
	if len(input) == 0 {
		return nil
	}

	hsmSetting := &input[0]

	result := amlfilesystems.AmlFilesystemHsmSettings{
		Container:        hsmSetting.Container,
		LoggingContainer: hsmSetting.LoggingContainer,
	}

	if hsmSetting.ImportPrefix != "" {
		result.ImportPrefix = pointer.To(hsmSetting.ImportPrefix)
	}

	return &result
}

func flattenAMLFileSystemHsmSetting(input *amlfilesystems.AmlFilesystemPropertiesHsm) []HsmSetting {
	if input == nil {
		return nil
	}

	var result []HsmSetting
	hsmSetting := HsmSetting{}

	if v := input.Settings; v != nil {
		hsmSetting.Container = v.Container
		hsmSetting.LoggingContainer = v.LoggingContainer

		if v.ImportPrefix != nil {
			hsmSetting.ImportPrefix = pointer.From(v.ImportPrefix)
		}
	}

	return append(result, hsmSetting)
}
