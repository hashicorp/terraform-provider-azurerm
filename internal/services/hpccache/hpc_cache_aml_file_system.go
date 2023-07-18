package hpccache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-05-01/amlfilesystems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type HPCCacheAMLFileSystemModel struct {
	Name              string                  `tfschema:"name"`
	ResourceGroupName string                  `tfschema:"resource_group_name"`
	Location          string                  `tfschema:"location"`
	SubnetId          string                  `tfschema:"subnet_id"`
	Identity          []AMLFileSystemIdentity `tfschema:"identity"`
	SkuName           string                  `tfschema:"sku_name"`
	Tags              map[string]string       `tfschema:"tags"`
	Zones             []string                `tfschema:"zones"`
}

type AMLFileSystemIdentity struct {
	PrincipalId            string                                   `tfschema:"principal_id"`
	TenantId               string                                   `tfschema:"tenant_id"`
	Type                   amlfilesystems.AmlFilesystemIdentityType `tfschema:"type"`
	UserAssignedIdentities string                                   `tfschema:"user_assigned_identities"`
}

type SkuName struct {
	Name string `tfschema:"name"`
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
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"identity": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"principal_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"tenant_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(amlfilesystems.AmlFilesystemIdentityTypeUserAssigned),
							string(amlfilesystems.AmlFilesystemIdentityTypeNone),
						}, false),
					},

					"user_assigned_identities": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						ValidateFunc:     validation.StringIsJSON,
						DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
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

			properties := &amlfilesystems.AmlFilesystem{
				Location: location.Normalize(model.Location),
				Properties: &amlfilesystems.AmlFilesystemProperties{
					FilesystemSubnet: model.SubnetId,
				},
				Tags:  &model.Tags,
				Zones: &model.Zones,
			}

			identity, err := expandAMLFileSystemIdentity(model.Identity)
			if err != nil {
				return err
			}
			properties.Identity = identity

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
			}

			if v := model.Sku; v != nil {
				state.SkuName = pointer.From(v.Name)
			}

			if model.Zones != nil {
				state.Zones = *model.Zones
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
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

func expandAMLFileSystemIdentity(input []AMLFileSystemIdentity) (*amlfilesystems.AmlFilesystemIdentity, error) {
	if len(input) == 0 {
		return nil, nil
	}

	amlFileSystemIdentity := &input[0]
	result := amlfilesystems.AmlFilesystemIdentity{
		Type: &amlFileSystemIdentity.Type,
	}

	var userAssignedIdentitiesValue map[string]amlfilesystems.UserAssignedIdentitiesProperties
	err := json.Unmarshal([]byte(amlFileSystemIdentity.UserAssignedIdentities), &userAssignedIdentitiesValue)
	if err != nil {
		return nil, err
	}
	result.UserAssignedIdentities = &userAssignedIdentitiesValue

	return &result, nil
}

func flattenAMLFileSystemIdentity(input *amlfilesystems.AmlFilesystemIdentity) ([]AMLFileSystemIdentity, error) {
	var amlFileSystemIdentities []AMLFileSystemIdentity
	if input == nil {
		return amlFileSystemIdentities, nil
	}

	amlFileSystemIdentity := AMLFileSystemIdentity{}

	if input.PrincipalId != nil {
		amlFileSystemIdentity.PrincipalId = *input.PrincipalId
	}

	if input.TenantId != nil {
		amlFileSystemIdentity.TenantId = *input.TenantId
	}

	if input.Type != nil {
		amlFileSystemIdentity.Type = *input.Type
	}

	if input.UserAssignedIdentities != nil && *input.UserAssignedIdentities != nil {
		userAssignedIdentitiesValue, err := json.Marshal(*input.UserAssignedIdentities)
		if err != nil {
			return nil, err
		}
		amlFileSystemIdentity.UserAssignedIdentities = string(userAssignedIdentitiesValue)
	}

	return append(amlFileSystemIdentities, amlFileSystemIdentity), nil
}
