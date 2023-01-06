package mobilenetwork

import (
	"context"
	"fmt"
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

type SimGroupDataSource struct{}

var _ sdk.DataSource = SimGroupDataSource{}

func (r SimGroupDataSource) ResourceType() string {
	return "azurerm_mobile_network_sim_group"
}

func (r SimGroupDataSource) ModelObject() interface{} {
	return &SimGroupModel{}
}

func (r SimGroupDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return simgroup.ValidateSimGroupID
}

func (r SimGroupDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),
	}
}

func (r SimGroupDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{

		"encryption_key_url": { // needs UserAssignedIdentity
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

		"location": commonschema.LocationComputed(),

		"mobile_network_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r SimGroupDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var metaModel SimGroupModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.SIMGroupClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := simgroup.NewSimGroupID(subscriptionId, metaModel.ResourceGroupName, metaModel.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
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

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
