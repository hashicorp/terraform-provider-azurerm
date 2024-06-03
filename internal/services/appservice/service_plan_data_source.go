// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ServicePlanDataSource struct{}

var _ sdk.DataSource = ServicePlanDataSource{}

type ServicePlanDataSourceModel struct {
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

func (r ServicePlanDataSource) ModelObject() interface{} {
	return &ServicePlanDataSourceModel{}
}

func (r ServicePlanDataSource) ResourceType() string {
	return "azurerm_service_plan"
}

func (r ServicePlanDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ServicePlanName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r ServicePlanDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

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

		"per_site_scaling_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"worker_count": {
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

		"zone_balancing_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"tags": tags.SchemaDataSource(),
	}
}

func (r ServicePlanDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.ServicePlanClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var servicePlan ServicePlanDataSourceModel
			if err := metadata.Decode(&servicePlan); err != nil {
				return err
			}

			id := commonids.NewAppServicePlanID(subscriptionId, servicePlan.ResourceGroup, servicePlan.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := existing.Model; model != nil {
				servicePlan.Location = location.Normalize(model.Location)
				servicePlan.Kind = pointer.From(model.Kind)

				if sku := model.Sku; sku != nil {
					if sku.Name != nil {
						servicePlan.Sku = *sku.Name
						if sku.Capacity != nil {
							servicePlan.WorkerCount = *sku.Capacity
						}
					}
				}

				if props := model.Properties; props != nil {
					servicePlan.OSType = OSTypeWindows
					if props.HyperV != nil && *props.HyperV {
						servicePlan.OSType = OSTypeWindowsContainer
					}
					if props.Reserved != nil && *props.Reserved {
						servicePlan.OSType = OSTypeLinux
					}

					if props.HostingEnvironmentProfile != nil && props.HostingEnvironmentProfile.Id != nil {
						servicePlan.AppServiceEnvironmentId = utils.NormalizeNilableString(props.HostingEnvironmentProfile.Id)
					}

					servicePlan.PerSiteScaling = utils.NormaliseNilableBool(props.PerSiteScaling)

					servicePlan.Reserved = utils.NormaliseNilableBool(props.Reserved)

					servicePlan.ZoneBalancing = utils.NormaliseNilableBool(props.ZoneRedundant)

					servicePlan.MaximumElasticWorkerCount = pointer.From(props.MaximumElasticWorkerCount)
				}
				servicePlan.Tags = pointer.From(model.Tags)
			}

			metadata.SetID(id)

			return metadata.Encode(&servicePlan)
		},
	}
}
