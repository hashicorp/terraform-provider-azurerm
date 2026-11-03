// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nginx_deployment

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxdeployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: time.Minute * 30,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Nginx.NginxDeployment

			id, err := nginxdeployment.ParseNginxDeploymentID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			var model DeploymentModel
			if err := meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding NginxDeploymentModel %s: %v", id, err)
			}

			var req nginxdeployment.NginxDeploymentUpdateParameters
			if meta.ResourceData.HasChange("sku") {
				req.Sku = &nginxdeployment.ResourceSku{Name: model.Sku}
			}

			if meta.ResourceData.HasChange("tags") {
				req.Tags = pointer.To(model.Tags)
			}

			if meta.ResourceData.HasChange("identity") {
				if req.Identity, err = identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity); err != nil {
					return fmt.Errorf("expanding identities: %+v", err)
				}
			}

			req.Properties = &nginxdeployment.NginxDeploymentUpdateProperties{}

			if meta.ResourceData.HasChange("capacity") && model.Capacity > 0 {
				req.Properties.ScalingProperties = &nginxdeployment.NginxDeploymentScalingProperties{
					Capacity: pointer.To(model.Capacity),
				}
			}

			if meta.ResourceData.HasChange("auto_scale_profile") && len(model.AutoScaleProfile) > 0 {
				var autoScaleProfiles []nginxdeployment.ScaleProfile
				for _, profile := range model.AutoScaleProfile {
					autoScaleProfiles = append(autoScaleProfiles, nginxdeployment.ScaleProfile{
						Name: profile.Name,
						Capacity: nginxdeployment.ScaleProfileCapacity{
							Min: profile.Min,
							Max: profile.Max,
						},
					})
				}
				req.Properties.ScalingProperties = &nginxdeployment.NginxDeploymentScalingProperties{
					AutoScaleSettings: &nginxdeployment.NginxDeploymentScalingPropertiesAutoScaleSettings{
						Profiles: autoScaleProfiles,
					},
				}
			}

			if meta.ResourceData.HasChange("email") {
				req.Properties.UserProfile = &nginxdeployment.NginxDeploymentUserProfile{
					PreferredEmail: pointer.To(model.Email),
				}
			}

			if meta.ResourceData.HasChange("automatic_upgrade_channel") {
				req.Properties.AutoUpgradeProfile = &nginxdeployment.AutoUpgradeProfile{
					UpgradeChannel: model.UpgradeChannel,
				}
			}

			if meta.ResourceData.HasChanges("frontend_public", "frontend_private", "network_interface") {
				req.Properties.NetworkProfile = expandNetworkProfile(model.FrontendPublic, model.FrontendPrivate, model.NetworkInterface)
			}

			if strings.HasPrefix(model.Sku, "basic") && req.Properties.ScalingProperties != nil {
				return fmt.Errorf("basic SKUs are incompatible with `capacity` or `auto_scale_profiles`")
			}

			if meta.ResourceData.HasChange("web_application_firewall") {
				activationState := nginxdeployment.ActivationStateDisabled
				if model.WebApplicationFirewall[0].ActivationStateEnabled {
					activationState = nginxdeployment.ActivationStateEnabled
				}
				req.Properties.NginxAppProtect = &nginxdeployment.NginxDeploymentUpdatePropertiesNginxAppProtect{
					WebApplicationFirewallSettings: &nginxdeployment.WebApplicationFirewallSettings{
						ActivationState: &activationState,
					},
				}
			}

			if err := client.DeploymentsUpdateThenPoll(ctx, *id, req); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}
