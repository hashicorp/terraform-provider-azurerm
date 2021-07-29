package appservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AppServicePlanDataSource struct{}

var _ sdk.DataSource = AppServicePlanDataSource{}

func (r AppServicePlanDataSource) ModelObject() interface{} {
	return AppServicePlanModel{}
}

func (r AppServicePlanDataSource) ResourceType() string {
	return "azurerm_service_plan"
}

func (r AppServicePlanDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ServicePlanName,
		},

		"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),
	}
}

func (r AppServicePlanDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": location.SchemaComputed(),

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"app_service_environment_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"os_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"per_site_scaling": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"number_of_workers": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"maximum_elastic_worker_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"kind": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"reserved": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"tags": tags.SchemaDataSource(),
	}
}

func (r AppServicePlanDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.ServicePlanClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var servicePlan AppServicePlanModel
			if err := metadata.Decode(&servicePlan); err != nil {
				return err
			}

			id := parse.NewServicePlanID(subscriptionId, servicePlan.ResourceGroup, servicePlan.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.ServerfarmName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("%s not found", id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			servicePlan.Location = location.NormalizeNilable(existing.Location)

			if sku := existing.Sku; sku != nil {
				if sku.Name != nil {
					servicePlan.Sku = *sku.Name
					if sku.Capacity != nil {
						servicePlan.NumberOfWorkers = int(*sku.Capacity)
					}
				}
			}

			if props := existing.AppServicePlanProperties; props != nil {
				if props.HyperV != nil && *props.HyperV {
					servicePlan.OSType = OSTypeWindowsContainer
				}
				if props.Reserved != nil && *props.Reserved {
					servicePlan.OSType = OSTypeLinux
				}

				if props.HostingEnvironmentProfile != nil {
					servicePlan.AppServiceEnvironmentId = utils.NormalizeNilableString(props.HostingEnvironmentProfile.ID)
				}

				servicePlan.PerSiteScaling = *props.PerSiteScaling
				servicePlan.Reserved = *props.Reserved
				servicePlan.MaximumElasticWorkerCount = int(utils.NormaliseNilableInt32(props.MaximumElasticWorkerCount))
			}
			servicePlan.Tags = tags.ToTypedObject(existing.Tags)

			metadata.SetID(id)

			return metadata.Encode(&servicePlan)
		},
	}
}
