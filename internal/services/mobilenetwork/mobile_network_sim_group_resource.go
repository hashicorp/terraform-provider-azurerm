package mobilenetwork

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/simgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SimGroupModel struct {
	Name              string                       `tfschema:"name"`
	ResourceGroupName string                       `tfschema:"resource_group_name"`
	EncryptionKeyUrl  string                       `tfschema:"encryption_key_url"`
	Identity          []identity.ModelUserAssigned `tfschema:"identity"`
	Location          string                       `tfschema:"location"`
	MobileNetworkId   string                       `tfschema:"mobile_network_id"`
	Tags              map[string]string            `tfschema:"tags"`
}

type SimGroupResource struct{}

var _ sdk.ResourceWithUpdate = SimGroupResource{}

func (r SimGroupResource) ResourceType() string {
	return "azurerm_mobile_network_sim_group"
}

func (r SimGroupResource) ModelObject() interface{} {
	return &SimGroupModel{}
}

func (r SimGroupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return simgroup.ValidateSimGroupID
}

func (r SimGroupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"encryption_key_url": { // needs UserAssignedIdentity
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"location": commonschema.Location(),

		"mobile_network_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r SimGroupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r SimGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SimGroupModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.SIMGroupClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := simgroup.NewSimGroupID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			identityValue, err := expandSimGroupIdentity(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			properties := &simgroup.SimGroup{
				Identity: identityValue,
				Location: location.Normalize(model.Location),
				Properties: simgroup.SimGroupPropertiesFormat{
					MobileNetwork: &simgroup.MobileNetworkResourceId{
						Id: model.MobileNetworkId,
					},
				},
				Tags: &model.Tags,
			}

			if model.EncryptionKeyUrl != "" {
				properties.Properties.EncryptionKey = &simgroup.KeyVaultKey{
					KeyUrl: &model.EncryptionKeyUrl,
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

func (r SimGroupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.SIMGroupClient

			id, err := simgroup.ParseSimGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SimGroupModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("identity") {
				identityValue, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				properties.Identity = identityValue
			}

			if metadata.ResourceData.HasChange("encryption_key") {
				properties.Properties.EncryptionKey = &simgroup.KeyVaultKey{
					KeyUrl: &model.EncryptionKeyUrl,
				}
			}

			if metadata.ResourceData.HasChange("mobile_network") {
				properties.Properties.MobileNetwork = &simgroup.MobileNetworkResourceId{
					Id: model.MobileNetworkId,
				}
			}

			properties.SystemData = nil

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r SimGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.SIMGroupClient

			id, err := simgroup.ParseSimGroupID(metadata.ResourceData.Id())
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

			state := SimGroupModel{
				Name:              id.SimGroupName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			identityValue, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}

			if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			properties := &model.Properties

			if properties.EncryptionKey != nil && properties.EncryptionKey.KeyUrl != nil {
				state.EncryptionKeyUrl = *properties.EncryptionKey.KeyUrl
			}

			if properties.MobileNetwork != nil {
				state.MobileNetworkId = properties.MobileNetwork.Id
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SimGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.SIMGroupClient

			id, err := simgroup.ParseSimGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err := resourceMobileNetworkChildWaitForDeletion(ctx, id.ID(), func() (*http.Response, error) {
				resp, err := client.Get(ctx, *id)
				return resp.HttpResponse, err
			}); err != nil {
				return err
			}

			return nil
		},
	}
}

func expandSimGroupIdentity(input []identity.ModelUserAssigned) (*identity.LegacySystemAndUserAssignedMap, error) {
	if len(input) == 0 {
		return nil, nil
	}

	identityValue, err := identity.ExpandUserAssignedMapFromModel(input)
	if err != nil {
		return nil, fmt.Errorf("expanding `identity`: %+v", err)
	}

	output := identity.LegacySystemAndUserAssignedMap{
		Type:        identityValue.Type,
		IdentityIds: identityValue.IdentityIds,
	}

	return &output, nil
}
