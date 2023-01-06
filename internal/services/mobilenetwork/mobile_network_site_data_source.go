package mobilenetwork

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/mobilenetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/site"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SiteDataSource struct{}

var _ sdk.DataSource = SiteDataSource{}

func (r SiteDataSource) ResourceType() string {
	return "azurerm_mobile_network_site"
}

func (r SiteDataSource) ModelObject() interface{} {
	return &SiteModel{}
}

func (r SiteDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return site.ValidateSiteID
}

func (r SiteDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"mobile_network_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: mobilenetwork.ValidateMobileNetworkID,
		},
	}
}

func (r SiteDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"tags": commonschema.TagsDataSource(),

		"network_function_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func (r SiteDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var metaModel SiteModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.SiteClient
			mobileNetworkId, err := mobilenetwork.ParseMobileNetworkID(metaModel.MobileNetworkMobileNetworkId)
			if err != nil {
				return err
			}

			id := site.NewSiteID(mobileNetworkId.SubscriptionId, mobileNetworkId.ResourceGroupName, mobileNetworkId.MobileNetworkName, metaModel.Name)

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

			state := SiteModel{
				Name:                         id.SiteName,
				MobileNetworkMobileNetworkId: mobilenetwork.NewMobileNetworkID(id.SubscriptionId, id.ResourceGroupName, id.MobileNetworkName).ID(),
				Location:                     location.Normalize(model.Location),
			}

			if properties := model.Properties; properties != nil {
				networkFunctionsValue, err := flattenSubResourceModel(properties.NetworkFunctions)
				if err != nil {
					return err
				}

				state.NetworkFunctions = networkFunctionsValue
			}
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
