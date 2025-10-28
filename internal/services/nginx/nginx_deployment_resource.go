// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nginx

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxdeployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type FrontendPrivate struct {
	IpAddress        string `tfschema:"ip_address"`
	AllocationMethod string `tfschema:"allocation_method"`
	SubnetId         string `tfschema:"subnet_id"`
}

type FrontendPublic struct {
	IpAddress []string `tfschema:"ip_address"`
}

type LoggingStorageAccount struct {
	Name          string `tfschema:"name"`
	ContainerName string `tfschema:"container_name"`
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

type DeploymentModel struct {
	ResourceGroupName      string                                     `tfschema:"resource_group_name"`
	Name                   string                                     `tfschema:"name"`
	NginxVersion           string                                     `tfschema:"nginx_version"`
	Identity               []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Sku                    string                                     `tfschema:"sku"`
	ManagedResourceGroup   string                                     `tfschema:"managed_resource_group,removedInNextMajorVersion"`
	Location               string                                     `tfschema:"location"`
	Capacity               int64                                      `tfschema:"capacity"`
	AutoScaleProfile       []AutoScaleProfile                         `tfschema:"auto_scale_profile"`
	DiagnoseSupportEnabled bool                                       `tfschema:"diagnose_support_enabled"`
	Email                  string                                     `tfschema:"email"`
	IpAddress              string                                     `tfschema:"ip_address"`
	LoggingStorageAccount  []LoggingStorageAccount                    `tfschema:"logging_storage_account,removedInNextMajorVersion"`
	FrontendPublic         []FrontendPublic                           `tfschema:"frontend_public"`
	FrontendPrivate        []FrontendPrivate                          `tfschema:"frontend_private"`
	NetworkInterface       []NetworkInterface                         `tfschema:"network_interface"`
	UpgradeChannel         string                                     `tfschema:"automatic_upgrade_channel"`
	WebApplicationFirewall []WebApplicationFirewall                   `tfschema:"web_application_firewall"`
	DataplaneAPIEndpoint   string                                     `tfschema:"dataplane_api_endpoint"`
	Tags                   map[string]string                          `tfschema:"tags"`
}

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
			alloc := nginxdeployment.NginxPrivateIPAllocationMethod(ip.AllocationMethod)
			privateIPs = append(privateIPs, nginxdeployment.NginxPrivateIPAddress{
				PrivateIPAddress:          pointer.To(ip.IpAddress),
				PrivateIPAllocationMethod: &alloc,
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

type DeploymentResource struct{}

var _ sdk.ResourceWithUpdate = (*DeploymentResource)(nil)

func (m DeploymentResource) Arguments() map[string]*pluginsdk.Schema {
	resource := map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"sku": {
			// docs: <https://docs.nginx.com/nginxaas/azure/billing/overview/>
			// we will not be forcing validation of SKU as there are internal SKUs
			// used for testing and for F5 NGINX private offers
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"location": commonschema.Location(),

		"capacity": {
			Type:          pluginsdk.TypeInt,
			Optional:      true,
			ConflictsWith: []string{"auto_scale_profile"},
			ValidateFunc:  validation.IntPositive,
		},

		"auto_scale_profile": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			ConflictsWith: []string{"capacity"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"min_capacity": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntPositive,
					},

					"max_capacity": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntPositive,
					},
				},
			},
		},

		"diagnose_support_enabled": {
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			ValidateFunc: nil,
		},

		"email": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"frontend_public": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"frontend_private"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_address": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},

		"frontend_private": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			ConflictsWith: []string{"frontend_public"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_address": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"allocation_method": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(nginxdeployment.PossibleValuesForNginxPrivateIPAllocationMethod(), false),
					},

					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"network_interface": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"automatic_upgrade_channel": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "stable",
			ValidateFunc: validation.StringInSlice(
				[]string{
					"stable",
					"preview",
				}, false),
		},

		"web_application_firewall": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"activation_state_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"status": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"attack_signatures_package": webApplicationFirewallPackageComputed(),
								"bot_signatures_package":    webApplicationFirewallPackageComputed(),
								"threat_campaigns_package":  webApplicationFirewallPackageComputed(),
								"component_versions":        webApplicationFirewallComponentVersionsComputed(),
							},
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}

	if !features.FivePointOh() {
		resource["managed_resource_group"] = &pluginsdk.Schema{
			Deprecated:   "The `managed_resource_group` field isn't supported by the API anymore and has been deprecated and will be removed in v5.0 of the AzureRM Provider.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		}

		resource["logging_storage_account"] = &pluginsdk.Schema{
			Deprecated: "The `logging_storage_account` block has been deprecated and will be removed in v5.0 of the AzureRM Provider. To enable logs, use the `azurerm_monitor_diagnostic_setting` resource instead.",
			Type:       pluginsdk.TypeList,
			Optional:   true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"container_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		}
	}
	return resource
}

func (m DeploymentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"nginx_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"dataplane_api_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (m DeploymentResource) ModelObject() interface{} {
	return &DeploymentModel{}
}

func (m DeploymentResource) ResourceType() string {
	return "azurerm_nginx_deployment"
}

func (m DeploymentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Nginx.NginxDeployment

			var model DeploymentModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			subscriptionID := meta.Client.Account.SubscriptionId
			id := nginxdeployment.NewNginxDeploymentID(subscriptionID, model.ResourceGroupName, model.Name)
			existing, err := client.DeploymentsGet(ctx, id)

			if !response.WasNotFound(existing.HttpResponse) {
				if err != nil {
					return fmt.Errorf("retrieving %s: %v", id, err)
				}
				return meta.ResourceRequiresImport(m.ResourceType(), id)
			}

			req := nginxdeployment.NginxDeployment{}
			req.Name = pointer.To(model.Name)
			req.Location = pointer.To(model.Location)
			req.Tags = pointer.FromMapOfStringStrings(model.Tags)

			if model.Sku != "" {
				sku := nginxdeployment.ResourceSku{Name: model.Sku}
				req.Sku = &sku
			}

			prop := &nginxdeployment.NginxDeploymentProperties{}

			if !features.FivePointOh() {
				if len(model.LoggingStorageAccount) > 0 {
					prop.Logging = &nginxdeployment.NginxLogging{
						StorageAccount: &nginxdeployment.NginxStorageAccount{
							AccountName:   pointer.FromString(model.LoggingStorageAccount[0].Name),
							ContainerName: pointer.FromString(model.LoggingStorageAccount[0].ContainerName),
						},
					}
				}
			}

			prop.EnableDiagnosticsSupport = pointer.FromBool(model.DiagnoseSupportEnabled)
			prop.NetworkProfile = expandNetworkProfile(model.FrontendPublic, model.FrontendPrivate, model.NetworkInterface)

			isBasicSKU := strings.HasPrefix(model.Sku, "basic")
			hasScaling := (model.Capacity > 0 || len(model.AutoScaleProfile) > 0)
			if isBasicSKU && hasScaling {
				return fmt.Errorf("basic SKUs are incompatible with `capacity` or `auto_scale_profiles`")
			}
			if !isBasicSKU && !hasScaling {
				return fmt.Errorf("scaling is required for `sku` '%s', please provide `capacity` or `auto_scale_profiles`", model.Sku)
			}

			if model.Capacity > 0 {
				prop.ScalingProperties = &nginxdeployment.NginxDeploymentScalingProperties{
					Capacity: pointer.FromInt64(model.Capacity),
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

			req.Identity, err = identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding identities: %+v", err)
			}

			err = client.DeploymentsCreateOrUpdateThenPoll(ctx, id, req)
			if err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			meta.SetID(id)
			return nil
		},
	}
}

func (m DeploymentResource) Read() sdk.ResourceFunc {
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
				output.Location = pointer.ToString(model.Location)
				output.Tags = pointer.ToMapOfStringStrings(model.Tags)
				if model.Sku != nil {
					output.Sku = model.Sku.Name
				}

				if props := model.Properties; props != nil {
					output.IpAddress = pointer.ToString(props.IPAddress)
					output.NginxVersion = pointer.ToString(props.NginxVersion)
					output.DataplaneAPIEndpoint = pointer.ToString(props.DataplaneApiEndpoint)
					output.DiagnoseSupportEnabled = pointer.ToBool(props.EnableDiagnosticsSupport)

					if !features.FivePointOh() {
						if props.Logging != nil && props.Logging.StorageAccount != nil {
							output.LoggingStorageAccount = []LoggingStorageAccount{
								{
									Name:          pointer.ToString(props.Logging.StorageAccount.AccountName),
									ContainerName: pointer.ToString(props.Logging.StorageAccount.ContainerName),
								},
							}
						}
					}

					if profile := props.NetworkProfile; profile != nil {
						if frontend := profile.FrontEndIPConfiguration; frontend != nil {
							if publicIps := frontend.PublicIPAddresses; publicIps != nil && len(*publicIps) > 0 {
								output.FrontendPublic = append(output.FrontendPublic, FrontendPublic{})
								for _, ip := range *publicIps {
									output.FrontendPublic[0].IpAddress = append(output.FrontendPublic[0].IpAddress, pointer.ToString(ip.Id))
								}
							}

							if privateIPs := frontend.PrivateIPAddresses; privateIPs != nil && len(*privateIPs) > 0 {
								for _, ip := range *privateIPs {
									method := ""
									if ip.PrivateIPAllocationMethod != nil {
										method = string(*ip.PrivateIPAllocationMethod)
									}

									output.FrontendPrivate = append(output.FrontendPrivate, FrontendPrivate{
										IpAddress:        pointer.ToString(ip.PrivateIPAddress),
										AllocationMethod: method,
										SubnetId:         pointer.ToString(ip.SubnetId),
									})
								}
							}
						}

						if netIf := profile.NetworkInterfaceConfiguration; netIf != nil {
							output.NetworkInterface = []NetworkInterface{
								{SubnetId: pointer.ToString(netIf.SubnetId)},
							}
						}
					}

					if scaling := props.ScalingProperties; scaling != nil {
						if capacity := scaling.Capacity; capacity != nil {
							output.Capacity = pointer.ToInt64(props.ScalingProperties.Capacity)
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
						output.Email = pointer.ToString(props.UserProfile.PreferredEmail)
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

func (m DeploymentResource) Update() sdk.ResourceFunc {
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
				req.Tags = pointer.FromMapOfStringStrings(model.Tags)
			}

			if meta.ResourceData.HasChange("identity") {
				if req.Identity, err = identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity); err != nil {
					return fmt.Errorf("expanding identities: %+v", err)
				}
			}

			req.Properties = &nginxdeployment.NginxDeploymentUpdateProperties{}
			if !features.FivePointOh() {
				if meta.ResourceData.HasChange("logging_storage_account") && len(model.LoggingStorageAccount) > 0 {
					req.Properties.Logging = &nginxdeployment.NginxLogging{
						StorageAccount: &nginxdeployment.NginxStorageAccount{
							AccountName:   pointer.FromString(model.LoggingStorageAccount[0].Name),
							ContainerName: pointer.FromString(model.LoggingStorageAccount[0].ContainerName),
						},
					}
				}
			}

			if meta.ResourceData.HasChange("diagnose_support_enabled") {
				req.Properties.EnableDiagnosticsSupport = pointer.FromBool(model.DiagnoseSupportEnabled)
			}

			if meta.ResourceData.HasChange("capacity") && model.Capacity > 0 {
				req.Properties.ScalingProperties = &nginxdeployment.NginxDeploymentScalingProperties{
					Capacity: pointer.FromInt64(model.Capacity),
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

func (m DeploymentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Nginx.NginxDeployment
			id, err := nginxdeployment.ParseNginxDeploymentID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			meta.Logger.Infof("deleting %s", id)

			if err := client.DeploymentsDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}

			return nil
		},
	}
}

func (m DeploymentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return nginxdeployment.ValidateNginxDeploymentID
}
