// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nginx_deployment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxdeployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := nginxdeployment.ParseNginxDeploymentID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			client := meta.Client.Nginx.NginxDeployment
			result, err := client.DeploymentsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return meta.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			output := DeploymentModel{
				Name:              id.NginxDeploymentName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := result.Model; model != nil {
				output.Location = pointer.From(model.Location)
				output.Tags = pointer.From(model.Tags)
				if model.Sku != nil {
					output.Sku = model.Sku.Name
				}

				if props := model.Properties; props != nil {
					output.IpAddress = pointer.From(props.IPAddress)
					output.NginxVersion = pointer.From(props.NginxVersion)
					output.DataplaneAPIEndpoint = pointer.From(props.DataplaneApiEndpoint)

					if profile := props.NetworkProfile; profile != nil {
						if frontend := profile.FrontEndIPConfiguration; frontend != nil {
							if publicIps := frontend.PublicIPAddresses; publicIps != nil && len(*publicIps) > 0 {
								output.FrontendPublic = append(output.FrontendPublic, FrontendPublic{})
								for _, ip := range *publicIps {
									output.FrontendPublic[0].IpAddress = append(output.FrontendPublic[0].IpAddress, pointer.From(ip.Id))
								}
							}

							if privateIPs := frontend.PrivateIPAddresses; privateIPs != nil && len(*privateIPs) > 0 {
								for _, ip := range *privateIPs {
									method := ""
									if ip.PrivateIPAllocationMethod != nil {
										method = string(*ip.PrivateIPAllocationMethod)
									}

									output.FrontendPrivate = append(output.FrontendPrivate, FrontendPrivate{
										IpAddress:        pointer.From(ip.PrivateIPAddress),
										AllocationMethod: method,
										SubnetId:         pointer.From(ip.SubnetId),
									})
								}
							}
						}

						if netIf := profile.NetworkInterfaceConfiguration; netIf != nil {
							output.NetworkInterface = []NetworkInterface{
								{SubnetId: pointer.From(netIf.SubnetId)},
							}
						}
					}

					if scaling := props.ScalingProperties; scaling != nil {
						if capacity := scaling.Capacity; capacity != nil {
							output.Capacity = pointer.From(props.ScalingProperties.Capacity)
						}
						if autoScaleProfiles := scaling.AutoScaleSettings; autoScaleProfiles != nil {
							profiles := autoScaleProfiles.Profiles
							for _, profile := range profiles {
								output.AutoScaleProfile = append(output.AutoScaleProfile, AutoScaleProfile{
									Name: profile.Name,
									Min:  profile.Capacity.Min,
									Max:  profile.Capacity.Max,
								})
							}
						}
					}

					if userProfile := props.UserProfile; userProfile != nil && userProfile.PreferredEmail != nil {
						output.Email = pointer.From(props.UserProfile.PreferredEmail)
					}

					if props.AutoUpgradeProfile != nil {
						output.UpgradeChannel = props.AutoUpgradeProfile.UpgradeChannel
					}

					if nap := props.NginxAppProtect; nap != nil {
						waf := WebApplicationFirewall{}
						if state := nap.WebApplicationFirewallSettings.ActivationState; state != nil {
							switch *state {
							case nginxdeployment.ActivationStateEnabled:
								waf.ActivationStateEnabled = true
							default:
								waf.ActivationStateEnabled = false
							}
						}
						if status := nap.WebApplicationFirewallStatus; status != nil {
							wafStatus := WebApplicationFirewallStatus{}
							if attackSignature := status.AttackSignaturesPackage; attackSignature != nil {
								wafStatus.AttackSignaturesPackage = []WebApplicationFirewallPackage{
									{
										RevisionDatetime: attackSignature.RevisionDatetime,
										Version:          attackSignature.Version,
									},
								}
							}
							if botSignature := status.BotSignaturesPackage; botSignature != nil {
								wafStatus.BotSignaturesPackage = []WebApplicationFirewallPackage{
									{
										RevisionDatetime: botSignature.RevisionDatetime,
										Version:          botSignature.Version,
									},
								}
							}
							if threatCampaign := status.ThreatCampaignsPackage; threatCampaign != nil {
								wafStatus.ThreatCampaignsPackage = []WebApplicationFirewallPackage{
									{
										RevisionDatetime: threatCampaign.RevisionDatetime,
										Version:          threatCampaign.Version,
									},
								}
							}
							if componentVersions := status.ComponentVersions; componentVersions != nil {
								wafStatus.ComponentVersions = []WebApplicationFirewallComponentVersions{
									{
										WafEngineVersion: componentVersions.WafEngineVersion,
										WafNginxVersion:  componentVersions.WafNginxVersion,
									},
								}
							}
							waf.Status = []WebApplicationFirewallStatus{wafStatus}
							output.WebApplicationFirewall = []WebApplicationFirewall{waf}
						}
					}

					flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
					if err != nil {
						return fmt.Errorf("flattening `identity`: %v", err)
					}
					output.Identity = *flattenedIdentity
				}
			}

			return meta.Encode(&output)
		},
	}
}

type FrontendPrivate struct {
	IpAddress        string `tfschema:"ip_address"`
	AllocationMethod string `tfschema:"allocation_method"`
	SubnetId         string `tfschema:"subnet_id"`
}

type FrontendPublic struct {
	IpAddress []string `tfschema:"ip_address"`
}

type NetworkInterface struct {
	SubnetId string `tfschema:"subnet_id"`
}

type AutoScaleProfile struct {
	Name string `tfschema:"name"`
	Min  int64  `tfschema:"min_capacity"`
	Max  int64  `tfschema:"max_capacity"`
}

type WebApplicationFirewall struct {
	ActivationStateEnabled bool                           `tfschema:"activation_state_enabled"`
	Status                 []WebApplicationFirewallStatus `tfschema:"status"`
}

type WebApplicationFirewallPackage struct {
	RevisionDatetime string `tfschema:"revision_datetime"`
	Version          string `tfschema:"version"`
}

type WebApplicationFirewallComponentVersions struct {
	WafEngineVersion string `tfschema:"waf_engine_version"`
	WafNginxVersion  string `tfschema:"waf_nginx_version"`
}

type WebApplicationFirewallStatus struct {
	AttackSignaturesPackage []WebApplicationFirewallPackage           `tfschema:"attack_signatures_package"`
	BotSignaturesPackage    []WebApplicationFirewallPackage           `tfschema:"bot_signatures_package"`
	ComponentVersions       []WebApplicationFirewallComponentVersions `tfschema:"component_versions"`
	ThreatCampaignsPackage  []WebApplicationFirewallPackage           `tfschema:"threat_campaigns_package"`
}
