package web

import (
	"context"
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/sdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AppServiceEnvironmentV3DataSource struct{}

var _ sdk.DataSource = AppServiceEnvironmentV3DataSource{}

func (r AppServiceEnvironmentV3DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AppServiceEnvironmentName,
		},

		"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),
	}
}

func (r AppServiceEnvironmentV3DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"subnet_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"cluster_setting": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"value": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"tags": tags.SchemaDataSource(),
	}
}

func (r AppServiceEnvironmentV3DataSource) ModelObject() interface{} {
	return AppServiceEnvironmentV3Model{}
}

func (r AppServiceEnvironmentV3DataSource) ResourceType() string {
	return "azurerm_app_service_environment_v3"
}

func (r AppServiceEnvironmentV3DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Web.AppServiceEnvironmentsClient
			id, err := parse.AppServiceEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.HostingEnvironmentName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := AppServiceEnvironmentV3Model{
				Name:          id.HostingEnvironmentName,
				ResourceGroup: id.ResourceGroup,
				Location:      location.NormalizeNilable(existing.Location),
			}

			if props := existing.AppServiceEnvironment; props != nil {
				if props.VirtualNetwork != nil {
					model.SubnetId = utils.NormalizeNilableString(props.VirtualNetwork.ID)
				}

				model.PricingTier = utils.NormalizeNilableString(props.MultiSize)

				model.ClusterSetting = flattenClusterSettingsModel(props.ClusterSettings)
			}

			model.Tags = tags.Flatten(existing.Tags)

			return metadata.Encode(&model)
		},
	}
}
