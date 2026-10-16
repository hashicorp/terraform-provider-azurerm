// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_service_environment_v3

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/appserviceenvironments"
)

func expandCreateForAppServiceEnvironmentV3(model AppServiceEnvironmentV3Model, vnetLoc string) appserviceenvironments.AppServiceEnvironmentResource {
	return appserviceenvironments.AppServiceEnvironmentResource{
		Kind:     pointer.To(KindASEV3),
		Location: location.Normalize(vnetLoc),
		Properties: &appserviceenvironments.AppServiceEnvironment{
			DedicatedHostCount:        pointer.To(model.DedicatedHostCount),
			ClusterSettings:           expandClusterSettingsModel(model.ClusterSetting),
			InternalLoadBalancingMode: pointer.ToEnum[appserviceenvironments.LoadBalancingMode](model.InternalLoadBalancingMode),
			VirtualNetwork: appserviceenvironments.VirtualNetworkProfile{
				Id: model.SubnetId,
			},
			ZoneRedundant: pointer.To(model.ZoneRedundant),
		},
		Tags: pointer.To(model.Tags),
	}
}

func expandClusterSettingsModel(input []ClusterSettingModel) *[]appserviceenvironments.NameValuePair {
	clusterSettings := make([]appserviceenvironments.NameValuePair, 0, len(input))
	if input == nil {
		return &clusterSettings
	}

	for _, v := range input {
		clusterSettings = append(clusterSettings, appserviceenvironments.NameValuePair{
			Name:  pointer.To(v.Name),
			Value: pointer.To(v.Value),
		})
	}

	return &clusterSettings
}
