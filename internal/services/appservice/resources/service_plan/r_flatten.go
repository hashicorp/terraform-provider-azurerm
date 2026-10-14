// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package service_plan

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/appserviceplans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) flatten(metadata sdk.ResourceMetaData, id *commonids.AppServicePlanId, model *appserviceplans.AppServicePlan) error {
	state := ServicePlanModel{
		Name:          id.ServerFarmName,
		ResourceGroup: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		state.Kind = pointer.From(model.Kind)

		// sku read
		if sku := model.Sku; sku != nil {
			state.Sku = pointer.From(sku.Name)
			state.WorkerCount = pointer.From(sku.Capacity)
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

			if pointer.From(props.ElasticScaleEnabled) && state.Sku != "" && helpers.PlanIsPremium(state.Sku) {
				state.PremiumPlanAutoScaleEnabled = pointer.From(props.ElasticScaleEnabled)
			}

			state.PerSiteScaling = pointer.From(props.PerSiteScaling)
			state.Reserved = pointer.From(props.Reserved)
			state.ZoneBalancing = pointer.From(props.ZoneRedundant)
			state.MaximumElasticWorkerCount = pointer.From(props.MaximumElasticWorkerCount)
		}
		state.Tags = pointer.From(model.Tags)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
