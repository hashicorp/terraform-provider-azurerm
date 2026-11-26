// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nginx_deployment

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxdeployment"
)

func expandNetworkProfile(public []FrontendPublic, private []FrontendPrivate, networkInterface []NetworkInterface) *nginxdeployment.NginxNetworkProfile {
	out := nginxdeployment.NginxNetworkProfile{
		FrontEndIPConfiguration:       &nginxdeployment.NginxFrontendIPConfiguration{},
		NetworkInterfaceConfiguration: &nginxdeployment.NginxNetworkInterfaceConfiguration{},
	}

	if len(public) > 0 && len(public[0].IpAddress) > 0 {
		var publicIPs []nginxdeployment.NginxPublicIPAddress
		for _, ip := range public[0].IpAddress {
			publicIPs = append(publicIPs, nginxdeployment.NginxPublicIPAddress{
				Id: pointer.To(ip),
			})
		}
		out.FrontEndIPConfiguration.PublicIPAddresses = &publicIPs
	}

	if len(private) > 0 {
		var privateIPs []nginxdeployment.NginxPrivateIPAddress
		for _, ip := range private {
			privateIPs = append(privateIPs, nginxdeployment.NginxPrivateIPAddress{
				PrivateIPAddress:          pointer.To(ip.IpAddress),
				PrivateIPAllocationMethod: pointer.ToEnum[nginxdeployment.NginxPrivateIPAllocationMethod](ip.AllocationMethod),
				SubnetId:                  pointer.To(ip.SubnetId),
			})
		}
		out.FrontEndIPConfiguration.PrivateIPAddresses = &privateIPs
	}

	if len(networkInterface) > 0 {
		out.NetworkInterfaceConfiguration.SubnetId = pointer.To(networkInterface[0].SubnetId)
	}

	return &out
}

func expandCreateForNginxDeployment(model DeploymentModel) (nginxdeployment.NginxDeployment, error) {
	req := nginxdeployment.NginxDeployment{}
	req.Name = pointer.To(model.Name)
	req.Location = pointer.To(model.Location)
	req.Tags = pointer.To(model.Tags)

	if model.Sku != "" {
		req.Sku = pointer.To(nginxdeployment.ResourceSku{Name: model.Sku})
	}

	prop := &nginxdeployment.NginxDeploymentProperties{}

	prop.NetworkProfile = expandNetworkProfile(model.FrontendPublic, model.FrontendPrivate, model.NetworkInterface)

	isBasicSKU := strings.HasPrefix(model.Sku, "basic")
	hasScaling := model.Capacity > 0 || len(model.AutoScaleProfile) > 0
	if isBasicSKU && hasScaling {
		return req, fmt.Errorf("basic SKUs are incompatible with `capacity` or `auto_scale_profiles`")
	}
	if !isBasicSKU && !hasScaling {
		return req, fmt.Errorf("scaling is required for `sku` '%s', please provide `capacity` or `auto_scale_profiles`", model.Sku)
	}

	if model.Capacity > 0 {
		prop.ScalingProperties = &nginxdeployment.NginxDeploymentScalingProperties{
			Capacity: pointer.To(model.Capacity),
		}
	}

	if autoScaleProfile := model.AutoScaleProfile; len(autoScaleProfile) > 0 {
		var autoScaleProfiles []nginxdeployment.ScaleProfile
		for _, profile := range autoScaleProfile {
			autoScaleProfiles = append(autoScaleProfiles, nginxdeployment.ScaleProfile{
				Name: profile.Name,
				Capacity: nginxdeployment.ScaleProfileCapacity{
					Min: profile.Min,
					Max: profile.Max,
				},
			})
		}
		prop.ScalingProperties = &nginxdeployment.NginxDeploymentScalingProperties{
			AutoScaleSettings: &nginxdeployment.NginxDeploymentScalingPropertiesAutoScaleSettings{
				Profiles: autoScaleProfiles,
			},
		}
	}

	if model.Email != "" {
		prop.UserProfile = &nginxdeployment.NginxDeploymentUserProfile{
			PreferredEmail: pointer.To(model.Email),
		}
	}

	if model.UpgradeChannel != "" {
		prop.AutoUpgradeProfile = &nginxdeployment.AutoUpgradeProfile{
			UpgradeChannel: model.UpgradeChannel,
		}
	}

	if len(model.WebApplicationFirewall) > 0 {
		activationState := nginxdeployment.ActivationStateDisabled
		if model.WebApplicationFirewall[0].ActivationStateEnabled {
			activationState = nginxdeployment.ActivationStateEnabled
		}

		prop.NginxAppProtect = &nginxdeployment.NginxDeploymentPropertiesNginxAppProtect{
			WebApplicationFirewallSettings: nginxdeployment.WebApplicationFirewallSettings{
				ActivationState: &activationState,
			},
		}
	}

	req.Properties = prop

	expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
	if err != nil {
		return req, fmt.Errorf("expanding identities: %+v", err)
	}

	req.Identity = expandedIdentity

	return req, nil
}
