// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package service_plan

import (
	"errors"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/appserviceplans"
)

func expandCreateForServicePlan(servicePlan ServicePlanModel) (appserviceplans.AppServicePlan, error) {
	appServicePlan := appserviceplans.AppServicePlan{
		Properties: &appserviceplans.AppServicePlanProperties{
			PerSiteScaling:      pointer.To(servicePlan.PerSiteScaling),
			Reserved:            pointer.To(servicePlan.OSType == OSTypeLinux),
			HyperV:              pointer.To(servicePlan.OSType == OSTypeWindowsContainer),
			ElasticScaleEnabled: pointer.To(servicePlan.PremiumPlanAutoScaleEnabled),
			ZoneRedundant:       pointer.To(servicePlan.ZoneBalancing),
		},
		Sku: &appserviceplans.SkuDescription{
			Name: pointer.To(servicePlan.Sku),
		},
		Location: location.Normalize(servicePlan.Location),
		Tags:     pointer.To(servicePlan.Tags),
	}

	if servicePlan.AppServiceEnvironmentId != "" {
		if !strings.HasPrefix(servicePlan.Sku, "I") {
			return appServicePlan, errors.New("'App Service Environment' based Service Plans can only be used with Isolated SKUs")
		}
		appServicePlan.Properties.HostingEnvironmentProfile = &appserviceplans.HostingEnvironmentProfile{
			Id: pointer.To(servicePlan.AppServiceEnvironmentId),
		}
	}

	if servicePlan.MaximumElasticWorkerCount > 0 {
		appServicePlan.Properties.MaximumElasticWorkerCount = pointer.To(servicePlan.MaximumElasticWorkerCount)
	}

	if servicePlan.WorkerCount != 0 {
		appServicePlan.Sku.Capacity = pointer.To(servicePlan.WorkerCount)
	}
	return appServicePlan, nil
}
