package web

import (
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/sdk"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
)

type ClusterSettingModel struct {
	Name  string `tfschema:"name"`
	Value string `tfschema:"value"`
}

type AppServiceEnvironmentV3Model struct {
	Name                      string `tfschema:"name"`
	ResourceGroup             string
	SubnetId                  string
	ClusterSetting            ClusterSettingModel
	InternalLoadBalancingMode string
	FrontEndScaleFactor       string
}

type AppServiceEnvironmentV3Resource struct{}

var _ sdk.Resource = AppServiceEnvironmentV3Resource{}
var _ sdk.ResourceWithUpdate = AppServiceEnvironmentV3Resource{}

func (a AppServiceEnvironmentV3Resource) Arguments() map[string]*schema.Schema {
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

		"internal_load_balancing_mode": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(web.LoadBalancingModeNone),
			ValidateFunc: validation.StringInSlice([]string{
				string(web.LoadBalancingModeNone),
				string(web.LoadBalancingModePublishing),
				string(web.LoadBalancingModeWeb),
				string(web.LoadBalancingModeWebPublishing),
				// (@jackofallops) breaking change in SDK - Enum for internal_load_balancing_mode changed from Web, Publishing to Web,Publishing
				string(LoadBalancingModeWebPublishing),
			}, false),
			DiffSuppressFunc: loadBalancingModeDiffSuppress,
		},

		"front_end_scale_factor": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      15,
			ValidateFunc: validation.IntBetween(5, 15),
		},

		"tags": tags.ForceNewSchema(),
	}
}

func (a AppServiceEnvironmentV3Resource) Attributes() map[string]*schema.Schema {
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

func (a AppServiceEnvironmentV3Resource) ModelObject() interface{} {
	panic("implement me")
}

func (a AppServiceEnvironmentV3Resource) ResourceType() string {
	panic("implement me")
}

func (a AppServiceEnvironmentV3Resource) Create() sdk.ResourceFunc {
	panic("implement me")
}

func (a AppServiceEnvironmentV3Resource) Read() sdk.ResourceFunc {
	panic("implement me")
}

func (a AppServiceEnvironmentV3Resource) Delete() sdk.ResourceFunc {
	panic("implement me")
}

func (a AppServiceEnvironmentV3Resource) IDValidationFunc() schema.SchemaValidateFunc {
	panic("implement me")
}

func (a AppServiceEnvironmentV3Resource) Update() sdk.ResourceFunc {
	panic("implement me")
}
