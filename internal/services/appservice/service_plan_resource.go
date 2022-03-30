package appservice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	webValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ServicePlanResource struct{}

var _ sdk.ResourceWithUpdate = ServicePlanResource{}

type OSType string

const (
	OSTypeLinux            OSType = "Linux"
	OSTypeWindows          OSType = "Windows"
	OSTypeWindowsContainer OSType = "WindowsContainer"
)

type ServicePlanModel struct {
	Name                      string            `tfschema:"name"`
	ResourceGroup             string            `tfschema:"resource_group_name"`
	Location                  string            `tfschema:"location"`
	Kind                      string            `tfschema:"kind"`
	OSType                    OSType            `tfschema:"os_type"`
	Sku                       string            `tfschema:"sku_name"`
	AppServiceEnvironmentId   string            `tfschema:"app_service_environment_id"`
	PerSiteScaling            bool              `tfschema:"per_site_scaling_enabled"`
	Reserved                  bool              `tfschema:"reserved"`
	WorkerCount               int               `tfschema:"worker_count"`
	MaximumElasticWorkerCount int               `tfschema:"maximum_elastic_worker_count"`
	AvailabilityZoneBalancing bool              `tfschema:"availability_zone_balancing_enabled"`
	Tags                      map[string]string `tfschema:"tags"`
}

func (r ServicePlanResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ServicePlanName,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"location": commonschema.Location(),

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice(
				helpers.AllKnownServicePlanSkus(),
				false),
		},

		"os_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(OSTypeLinux),
				string(OSTypeWindows),
				string(OSTypeWindowsContainer),
			}, false),
		},

		"app_service_environment_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: webValidate.AppServiceEnvironmentID,
		},

		"per_site_scaling_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"worker_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"maximum_elastic_worker_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"availability_zone_balancing_enabled": {
			Type:     pluginsdk.TypeBool,
			ForceNew: true,
			Optional: true,
		},

		"tags": tags.Schema(),
	}
}

func (r ServicePlanResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"kind": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"reserved": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
	}
}

func (r ServicePlanResource) ModelObject() interface{} {
	return &ServicePlanModel{}
}

func (r ServicePlanResource) ResourceType() string {
	return "azurerm_service_plan"
}

func (r ServicePlanResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var servicePlan ServicePlanModel
			if err := metadata.Decode(&servicePlan); err != nil {
				return err
			}

			client := metadata.Client.AppService.ServicePlanClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := parse.NewServicePlanID(subscriptionId, servicePlan.ResourceGroup, servicePlan.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.ServerfarmName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("retreiving %s: %v", id, err)
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			appServicePlan := web.AppServicePlan{
				AppServicePlanProperties: &web.AppServicePlanProperties{
					PerSiteScaling: utils.Bool(servicePlan.PerSiteScaling),
					Reserved:       utils.Bool(servicePlan.OSType == OSTypeLinux),
					HyperV:         utils.Bool(servicePlan.OSType == OSTypeWindowsContainer),
					ZoneRedundant:  utils.Bool(servicePlan.AvailabilityZoneBalancing),
				},
				Sku: &web.SkuDescription{
					Name: utils.String(servicePlan.Sku),
				},
				Location: utils.String(location.Normalize(servicePlan.Location)),
				Tags:     tags.FromTypedObject(servicePlan.Tags),
			}

			if servicePlan.AppServiceEnvironmentId != "" {
				if !strings.HasPrefix(servicePlan.Sku, "I") {
					return fmt.Errorf("App Service Environment based Service Plans can only be used with Isolated SKUs")
				}
				appServicePlan.AppServicePlanProperties.HostingEnvironmentProfile = &web.HostingEnvironmentProfile{
					ID: utils.String(servicePlan.AppServiceEnvironmentId),
				}
			}

			if servicePlan.MaximumElasticWorkerCount > 0 {
				if !strings.HasPrefix(servicePlan.Sku, "EP") && !strings.HasPrefix(servicePlan.Sku, "PC") {
					return fmt.Errorf("`maximum_elastic_worker_count` can only be specified with Elastic Premium Skus")
				}
				appServicePlan.AppServicePlanProperties.MaximumElasticWorkerCount = utils.Int32(int32(servicePlan.MaximumElasticWorkerCount))
			}

			if servicePlan.WorkerCount != 0 {
				appServicePlan.Sku.Capacity = utils.Int32(int32(servicePlan.WorkerCount))
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerfarmName, appServicePlan)
			if err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waitng for creation of %s: %v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ServicePlanResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.ServicePlanClient
			id, err := parse.ServicePlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			servicePlan, err := client.Get(ctx, id.ResourceGroup, id.ServerfarmName)
			if err != nil {
				if utils.ResponseWasNotFound(servicePlan.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			state := ServicePlanModel{
				Name:          id.ServerfarmName,
				ResourceGroup: id.ResourceGroup,
				Location:      location.NormalizeNilable(servicePlan.Location),
				Kind:          utils.NormalizeNilableString(servicePlan.Kind),
			}

			// sku read
			if sku := servicePlan.Sku; sku != nil {
				if sku.Name != nil {
					state.Sku = *sku.Name
					if sku.Capacity != nil {
						state.WorkerCount = int(*sku.Capacity)
					}
				}
			}

			// props read
			if props := servicePlan.AppServicePlanProperties; props != nil {
				state.OSType = OSTypeWindows
				if props.HyperV != nil && *props.HyperV {
					state.OSType = OSTypeWindowsContainer
				}
				if props.Reserved != nil && *props.Reserved {
					state.OSType = OSTypeLinux
				}

				if ase := props.HostingEnvironmentProfile; ase != nil && ase.ID != nil {
					state.AppServiceEnvironmentId = *ase.ID
				}

				state.PerSiteScaling = utils.NormaliseNilableBool(props.PerSiteScaling)

				state.Reserved = utils.NormaliseNilableBool(props.Reserved)

				state.AvailabilityZoneBalancing = utils.NormaliseNilableBool(props.ZoneRedundant)

				state.MaximumElasticWorkerCount = int(utils.NormaliseNilableInt32(props.MaximumElasticWorkerCount))
			}
			state.Tags = tags.ToTypedObject(servicePlan.Tags)

			return metadata.Encode(&state)
		},
	}
}

func (r ServicePlanResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.ServicePlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.AppService.ServicePlanClient
			metadata.Logger.Infof("deleting %s", id)

			if _, err := client.Delete(ctx, id.ResourceGroup, id.ServerfarmName); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}

			return nil
		},
	}
}

func (r ServicePlanResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ServicePlanID
}

func (r ServicePlanResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.ServicePlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.AppService.ServicePlanClient

			var state ServicePlanModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.ServerfarmName)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("per_site_scaling_enabled") {
				existing.AppServicePlanProperties.PerSiteScaling = utils.Bool(state.PerSiteScaling)
			}

			if metadata.ResourceData.HasChange("sku_name") {
				existing.Sku.Name = utils.String(state.Sku)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Tags = tags.FromTypedObject(state.Tags)
			}

			if metadata.ResourceData.HasChange("worker_count") {
				existing.Sku.Capacity = utils.Int32(int32(state.WorkerCount))
			}

			if metadata.ResourceData.HasChange("maximum_elastic_worker_count") {
				if metadata.ResourceData.HasChange("maximum_elastic_worker_count") && !strings.HasPrefix(state.Sku, "EP") && !strings.HasPrefix(state.Sku, "PC") {
					return fmt.Errorf("`maximum_elastic_worker_count` can only be specified with Elastic Premium Skus")
				}
				existing.AppServicePlanProperties.MaximumElasticWorkerCount = utils.Int32(int32(state.MaximumElasticWorkerCount))
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerfarmName, existing)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update to %s: %+v", id, err)
			}

			return nil
		},
	}
}
