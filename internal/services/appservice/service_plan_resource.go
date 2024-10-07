// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/appserviceplans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	webValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ServicePlanResource struct{}

var _ sdk.ResourceWithUpdate = ServicePlanResource{}

var _ sdk.ResourceWithStateMigration = ServicePlanResource{}

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
	WorkerCount               int64             `tfschema:"worker_count"`
	MaximumElasticWorkerCount int64             `tfschema:"maximum_elastic_worker_count"`
	ZoneBalancing             bool              `tfschema:"zone_balancing_enabled"`
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

		"resource_group_name": commonschema.ResourceGroupName(),

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

		"zone_balancing_enabled": {
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

			id := commonids.NewAppServicePlanID(subscriptionId, servicePlan.ResourceGroup, servicePlan.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("retreiving %s: %v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			appServicePlan := appserviceplans.AppServicePlan{
				Properties: &appserviceplans.AppServicePlanProperties{
					PerSiteScaling: pointer.To(servicePlan.PerSiteScaling),
					Reserved:       pointer.To(servicePlan.OSType == OSTypeLinux),
					HyperV:         pointer.To(servicePlan.OSType == OSTypeWindowsContainer),
					ZoneRedundant:  pointer.To(servicePlan.ZoneBalancing),
				},
				Sku: &appserviceplans.SkuDescription{
					Name: pointer.To(servicePlan.Sku),
				},
				Location: location.Normalize(servicePlan.Location),
				Tags:     pointer.To(servicePlan.Tags),
			}

			if servicePlan.AppServiceEnvironmentId != "" {
				if !strings.HasPrefix(servicePlan.Sku, "I") {
					return fmt.Errorf("App Service Environment based Service Plans can only be used with Isolated SKUs")
				}
				appServicePlan.Properties.HostingEnvironmentProfile = &appserviceplans.HostingEnvironmentProfile{
					Id: utils.String(servicePlan.AppServiceEnvironmentId),
				}
			}

			if servicePlan.MaximumElasticWorkerCount > 0 {
				if !isServicePlanSupportScaleOut(servicePlan.Sku) {
					return fmt.Errorf("`maximum_elastic_worker_count` can only be specified with Elastic Premium Skus")
				}
				appServicePlan.Properties.MaximumElasticWorkerCount = pointer.To(servicePlan.MaximumElasticWorkerCount)
			}

			if servicePlan.WorkerCount != 0 {
				appServicePlan.Sku.Capacity = pointer.To(servicePlan.WorkerCount)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, appServicePlan); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
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
			id, err := commonids.ParseAppServicePlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			servicePlan, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(servicePlan.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			state := ServicePlanModel{
				Name:          id.ServerFarmName,
				ResourceGroup: id.ResourceGroupName,
			}

			if model := servicePlan.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Kind = pointer.From(model.Kind)

				// sku read
				if sku := model.Sku; sku != nil {
					if sku.Name != nil {
						state.Sku = *sku.Name
						if sku.Capacity != nil {
							state.WorkerCount = *sku.Capacity
						}
					}
				}

				// props read
				if props := model.Properties; props != nil {
					state.OSType = OSTypeWindows
					if props.HyperV != nil && *props.HyperV {
						state.OSType = OSTypeWindowsContainer
					}
					if props.Reserved != nil && *props.Reserved {
						state.OSType = OSTypeLinux
					}

					if ase := props.HostingEnvironmentProfile; ase != nil && ase.Id != nil {
						state.AppServiceEnvironmentId = *ase.Id
					}

					state.PerSiteScaling = utils.NormaliseNilableBool(props.PerSiteScaling)

					state.Reserved = utils.NormaliseNilableBool(props.Reserved)

					state.ZoneBalancing = utils.NormaliseNilableBool(props.ZoneRedundant)

					state.MaximumElasticWorkerCount = pointer.From(props.MaximumElasticWorkerCount)
				}
				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ServicePlanResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := commonids.ParseAppServicePlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.AppService.ServicePlanClient
			metadata.Logger.Infof("deleting %s", id)

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}

			return nil
		},
	}
}

func (r ServicePlanResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateAppServicePlanID
}

func (r ServicePlanResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := commonids.ParseAppServicePlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.AppService.ServicePlanClient

			var state ServicePlanModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			model := *existing.Model

			if metadata.ResourceData.HasChange("per_site_scaling_enabled") {
				model.Properties.PerSiteScaling = pointer.To(state.PerSiteScaling)
			}

			if metadata.ResourceData.HasChange("sku_name") {
				model.Sku.Name = utils.String(state.Sku)
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = pointer.To(state.Tags)
			}

			if metadata.ResourceData.HasChange("worker_count") {
				model.Sku.Capacity = pointer.To(state.WorkerCount)
			}

			if metadata.ResourceData.HasChange("maximum_elastic_worker_count") {
				if metadata.ResourceData.HasChange("maximum_elastic_worker_count") && !isServicePlanSupportScaleOut(state.Sku) {
					return fmt.Errorf("`maximum_elastic_worker_count` can only be specified with Elastic Premium Skus")
				}
				model.Properties.MaximumElasticWorkerCount = pointer.To(state.MaximumElasticWorkerCount)
			}

			if err = client.CreateOrUpdateThenPoll(ctx, *id, model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func isServicePlanSupportScaleOut(plan string) bool {
	support := false
	support = support || strings.HasPrefix(plan, "EP")
	support = support || strings.HasPrefix(plan, "PC")
	support = support || strings.HasPrefix(plan, "WS")
	return support
}

func (r ServicePlanResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.ServicePlanV0toV1{},
		},
	}
}
