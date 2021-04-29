package web

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/response"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/sdk"
	networkParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

const KindASEV3 = "ASEV3"

type ClusterSettingModel struct {
	Name  string `tfschema:"name"`
	Value string `tfschema:"value"`
}

type AppServiceEnvironmentV3Model struct {
	Name           string                 `tfschema:"name"`
	ResourceGroup  string                 `tfschema:"resource_group_name"`
	SubnetId       string                 `tfschema:"subnet_id"`
	ClusterSetting []ClusterSettingModel  `tfschema:"cluster_setting"`
	PricingTier    string                 `tfschema:"pricing_tier"`
	Location       string                 `tfschema:"location"`
	Tags           map[string]interface{} `tfschema:"tags"`
}

// (@jackofallops) - Two important properties are missing from the SDK / Swagger that will need to be added later
// these are `dedicated_host_count` https://docs.microsoft.com/en-gb/azure/app-service/environment/creation#dedicated-hosts
// and `upgrade_preference` https://docs.microsoft.com/en-us/azure/app-service/environment/using#upgrade-preference

type AppServiceEnvironmentV3Resource struct{}

var _ sdk.Resource = AppServiceEnvironmentV3Resource{}
var _ sdk.ResourceWithUpdate = AppServiceEnvironmentV3Resource{}

func (r AppServiceEnvironmentV3Resource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AppServiceEnvironmentName,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"subnet_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: networkValidate.SubnetID,
		},

		"cluster_setting": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"value": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},

		"tags": tags.ForceNewSchema(),
	}
}

func (r AppServiceEnvironmentV3Resource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"pricing_tier": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"location": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (r AppServiceEnvironmentV3Resource) ModelObject() interface{} {
	return AppServiceEnvironmentV3Model{}
}

func (r AppServiceEnvironmentV3Resource) ResourceType() string {
	return "azurerm_app_service_environment_v3"
}

func (r AppServiceEnvironmentV3Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 6 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Web.AppServiceEnvironmentsClient
			networksClient := metadata.Client.Network.VnetClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AppServiceEnvironmentV3Model
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			subnet, err := networkParse.SubnetID(model.SubnetId)
			if err != nil {
				return err
			}

			vnet, err := networksClient.Get(ctx, subnet.ResourceGroup, subnet.VirtualNetworkName, "")
			if err != nil {
				return fmt.Errorf("retrieving Virtual Network %q (Resource Group %q): %+v", subnet.VirtualNetworkName, subnet.ResourceGroup, err)
			}

			vnetLoc := location.NormalizeNilable(vnet.Location)
			if vnetLoc == "" {
				return fmt.Errorf("determining Location from Virtual Network %q (Resource Group %q): `location` was missing", subnet.VirtualNetworkName, subnet.ResourceGroup)
			}

			id := parse.NewAppServiceEnvironmentID(subscriptionId, model.ResourceGroup, model.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.HostingEnvironmentName)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			envelope := web.AppServiceEnvironmentResource{
				Kind:     utils.String(KindASEV3),
				Location: utils.String(vnetLoc),
				AppServiceEnvironment: &web.AppServiceEnvironment{
					Name: utils.String(id.HostingEnvironmentName),
					VirtualNetwork: &web.VirtualNetworkProfile{
						ID:     utils.String(model.SubnetId),
						Subnet: utils.String(subnet.Name),
					},
					ClusterSettings: expandClusterSettingsModel(model.ClusterSetting),
				},
				Tags: tags.Expand(model.Tags),
			}

			if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.HostingEnvironmentName, envelope); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			createWait := resource.StateChangeConf{
				Pending: []string{
					string(web.ProvisioningStateInProgress),
				},
				Target: []string{
					string(web.ProvisioningStateSucceeded),
				},
				MinTimeout: 1 * time.Minute,
				Refresh:    appServiceEnvironmentRefresh(ctx, client, id.ResourceGroup, id.HostingEnvironmentName),
			}

			timeout, _ := ctx.Deadline()
			createWait.Timeout = time.Until(timeout)

			if _, err := createWait.WaitForState(); err != nil {
				return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AppServiceEnvironmentV3Resource) Read() sdk.ResourceFunc {
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

func (r AppServiceEnvironmentV3Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 6 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Web.AppServiceEnvironmentsClient

			id, err := parse.AppServiceEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.HostingEnvironmentName, utils.Bool(false))
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				// This future can return a 404 for the polling check if the ASE is successfully deleted but this raises an error in the SDK
				if !response.WasNotFound(future.Response()) {
					return fmt.Errorf("waiting for removal of %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func (r AppServiceEnvironmentV3Resource) IDValidationFunc() schema.SchemaValidateFunc {
	return validate.AppServiceEnvironmentID
}

func (r AppServiceEnvironmentV3Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 6 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.AppServiceEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Info("Decoding state...")
			var state AppServiceEnvironmentV3Model
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			metadata.Logger.Infof("updating %s", id)
			client := metadata.Client.Web.AppServiceEnvironmentsClient

			patch := web.AppServiceEnvironmentPatchResource{
				AppServiceEnvironment: &web.AppServiceEnvironment{},
			}

			if metadata.ResourceData.HasChange("cluster_setting") {
				patch.AppServiceEnvironment.ClusterSettings = expandClusterSettingsModel(state.ClusterSetting)
			}

			if _, err = client.Update(ctx, id.ResourceGroup, id.HostingEnvironmentName, patch); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func flattenClusterSettingsModel(input *[]web.NameValuePair) []ClusterSettingModel {
	var output []ClusterSettingModel
	if input == nil || len(*input) == 0 {
		return output
	}

	for _, v := range *input {
		if v.Name == nil {
			continue
		}

		output = append(output, ClusterSettingModel{
			Name:  *v.Name,
			Value: utils.NormalizeNilableString(v.Value),
		})
	}
	return output
}

func expandClusterSettingsModel(input []ClusterSettingModel) *[]web.NameValuePair {
	var clusterSettings []web.NameValuePair
	if input == nil {
		return &clusterSettings
	}

	for _, v := range input {
		clusterSettings = append(clusterSettings, web.NameValuePair{
			Name:  utils.String(v.Name),
			Value: utils.String(v.Value),
		})
	}
	return &clusterSettings
}
