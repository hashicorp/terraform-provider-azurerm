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
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-06-01-preview/nginxconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-06-01-preview/nginxdeployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const defaultCapacity = 20 // TODO: remove this in v4.0

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

// Deprecated: remove in next major version
type ConfigureFile struct {
	Content     string `tfschema:"content"`
	VirtualPath string `tfschema:"virtual_path"`
}

// Deprecated: remove in next major version
type Configuration struct {
	ConfigureFile []ConfigureFile `tfschema:"config_file"`
	ProtectedFile []ConfigureFile `tfschema:"protected_file"`
	PackageData   string          `tfschema:"package_data"`
	RootFile      string          `tfschema:"root_file"`
}

type AutoScaleProfile struct {
	Name string `tfschema:"name"`
	Min  int64  `tfschema:"min_capacity"`
	Max  int64  `tfschema:"max_capacity"`
}

type DeploymentModel struct {
	ResourceGroupName      string                                     `tfschema:"resource_group_name"`
	Name                   string                                     `tfschema:"name"`
	NginxVersion           string                                     `tfschema:"nginx_version"`
	Identity               []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Sku                    string                                     `tfschema:"sku"`
	ManagedResourceGroup   string                                     `tfschema:"managed_resource_group"`
	Location               string                                     `tfschema:"location"`
	Capacity               int64                                      `tfschema:"capacity"`
	AutoScaleProfile       []AutoScaleProfile                         `tfschema:"auto_scale_profile"`
	DiagnoseSupportEnabled bool                                       `tfschema:"diagnose_support_enabled"`
	Email                  string                                     `tfschema:"email"`
	IpAddress              string                                     `tfschema:"ip_address"`
	LoggingStorageAccount  []LoggingStorageAccount                    `tfschema:"logging_storage_account"`
	FrontendPublic         []FrontendPublic                           `tfschema:"frontend_public"`
	FrontendPrivate        []FrontendPrivate                          `tfschema:"frontend_private"`
	NetworkInterface       []NetworkInterface                         `tfschema:"network_interface"`
	UpgradeChannel         string                                     `tfschema:"automatic_upgrade_channel"`
	// Deprecated: remove in next major version
	Configuration []Configuration   `tfschema:"configuration,removedInNextMajorVersion"`
	Tags          map[string]string `tfschema:"tags"`
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

		"managed_resource_group": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

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

		"logging_storage_account": {
			Type:     pluginsdk.TypeList,
			Optional: true,
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
		},

		"frontend_public": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			ForceNew:      true,
			MaxItems:      1,
			ConflictsWith: []string{"frontend_private"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_address": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
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
			ForceNew:      true,
			ConflictsWith: []string{"frontend_public"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_address": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},

					"allocation_method": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(nginxdeployment.PossibleValuesForNginxPrivateIPAllocationMethod(), false),
					},

					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},
				},
			},
		},

		"network_interface": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
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

		"tags": commonschema.Tags(),
	}

	if !features.FourPointOhBeta() {
		resource["capacity"].Default = defaultCapacity

		resource["configuration"] = &pluginsdk.Schema{
			Deprecated: "The `configuration` block has been superseded by the `azurerm_nginx_configuration` resource and will be removed in v4.0 of the AzureRM Provider.",
			Type:       pluginsdk.TypeList,
			Optional:   true,
			Computed:   true,
			MaxItems:   1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"config_file": {
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						AtLeastOneOf: []string{"configuration.0.config_file", "configuration.0.package_data"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"content": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsBase64,
								},

								"virtual_path": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"protected_file": {
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						RequiredWith: []string{"configuration.0.config_file"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"content": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsBase64,
								},

								"virtual_path": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"package_data": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						AtLeastOneOf:  []string{"configuration.0.config_file", "configuration.0.package_data"},
						ConflictsWith: []string{"configuration.0.protected_file", "configuration.0.config_file"},
					},

					"root_file": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
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
			req.Name = pointer.FromString(model.Name)
			req.Location = pointer.FromString(model.Location)
			req.Tags = pointer.FromMapOfStringStrings(model.Tags)

			if model.Sku != "" {
				sku := nginxdeployment.ResourceSku{Name: model.Sku}
				req.Sku = &sku
			}

			prop := &nginxdeployment.NginxDeploymentProperties{}
			prop.ManagedResourceGroup = pointer.FromString(model.ManagedResourceGroup)

			if len(model.LoggingStorageAccount) > 0 {
				prop.Logging = &nginxdeployment.NginxLogging{
					StorageAccount: &nginxdeployment.NginxStorageAccount{
						AccountName:   pointer.FromString(model.LoggingStorageAccount[0].Name),
						ContainerName: pointer.FromString(model.LoggingStorageAccount[0].ContainerName),
					},
				}
			}

			prop.EnableDiagnosticsSupport = pointer.FromBool(model.DiagnoseSupportEnabled)
			prop.NetworkProfile = &nginxdeployment.NginxNetworkProfile{
				FrontEndIPConfiguration:       &nginxdeployment.NginxFrontendIPConfiguration{},
				NetworkInterfaceConfiguration: &nginxdeployment.NginxNetworkInterfaceConfiguration{},
			}

			if public := model.FrontendPublic; len(public) > 0 && len(public[0].IpAddress) > 0 {
				var publicIPs []nginxdeployment.NginxPublicIPAddress
				for _, ip := range public[0].IpAddress {
					publicIPs = append(publicIPs, nginxdeployment.NginxPublicIPAddress{
						Id: pointer.FromString(ip),
					})
				}
				prop.NetworkProfile.FrontEndIPConfiguration.PublicIPAddresses = &publicIPs
			}

			if private := model.FrontendPrivate; len(private) > 0 {
				var privateIPs []nginxdeployment.NginxPrivateIPAddress
				for _, ip := range private {
					alloc := nginxdeployment.NginxPrivateIPAllocationMethod(ip.AllocationMethod)
					privateIPs = append(privateIPs, nginxdeployment.NginxPrivateIPAddress{
						PrivateIPAddress:          pointer.FromString(ip.IpAddress),
						PrivateIPAllocationMethod: &alloc,
						SubnetId:                  pointer.FromString(ip.SubnetId),
					})
				}
				prop.NetworkProfile.FrontEndIPConfiguration.PrivateIPAddresses = &privateIPs
			}

			if len(model.NetworkInterface) > 0 {
				prop.NetworkProfile.NetworkInterfaceConfiguration.SubnetId = pointer.FromString(model.NetworkInterface[0].SubnetId)
			}

			isBasicSKU := strings.HasPrefix(model.Sku, "basic")
			if !features.FourPointOhBeta() {
				if isBasicSKU && (model.Capacity != defaultCapacity || len(model.AutoScaleProfile) > 0) {
					return fmt.Errorf("basic SKUs are incompatible with `capacity` or `auto_scale_profiles`")
				}

				if model.Capacity > 0 && !isBasicSKU {
					prop.ScalingProperties = &nginxdeployment.NginxDeploymentScalingProperties{
						Capacity: pointer.FromInt64(model.Capacity),
					}
				}
			} else {
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
					PreferredEmail: pointer.FromString(model.Email),
				}
			}

			if model.UpgradeChannel != "" {
				prop.AutoUpgradeProfile = &nginxdeployment.AutoUpgradeProfile{
					UpgradeChannel: model.UpgradeChannel,
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

			if !features.FourPointOhBeta() {
				if len(model.Configuration) > 0 {
					// update configuration
					configID := nginxconfiguration.NewConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.NginxDeploymentName, defaultConfigurationName)

					configProp := expandConfiguration(model.Configuration[0])
					if err := meta.Client.Nginx.NginxConfiguration.ConfigurationsCreateOrUpdateThenPoll(ctx, configID, configProp); err != nil {
						return fmt.Errorf("update default configuration of %q: %v", configID, err)
					}
				}
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
					output.ManagedResourceGroup = pointer.ToString(props.ManagedResourceGroup)
					output.NginxVersion = pointer.ToString(props.NginxVersion)
					output.DiagnoseSupportEnabled = pointer.ToBool(props.EnableDiagnosticsSupport)

					if props.Logging != nil && props.Logging.StorageAccount != nil {
						output.LoggingStorageAccount = []LoggingStorageAccount{
							{
								Name:          pointer.ToString(props.Logging.StorageAccount.AccountName),
								ContainerName: pointer.ToString(props.Logging.StorageAccount.ContainerName),
							},
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

					flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
					if err != nil {
						return fmt.Errorf("flattening `identity`: %v", err)
					}
					output.Identity = *flattenedIdentity
				}
			}

			if !features.FourPointOhBeta() {
				if v := meta.ResourceData.Get("configuration"); len(v.([]interface{})) != 0 {
					// read configuration
					configResp, err := meta.Client.Nginx.NginxConfiguration.ConfigurationsGet(ctx, nginxconfiguration.NewConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.NginxDeploymentName, defaultConfigurationName))
					if err != nil && !response.WasNotFound(configResp.HttpResponse) {
						return fmt.Errorf("retrieving default configuration of %q: %v", id, err)
					}
					if model := configResp.Model; model != nil {
						if prop := model.Properties; prop != nil {
							var files []ConfigureFile
							if prop.Files != nil {
								for _, file := range *prop.Files {
									files = append(files, ConfigureFile{
										Content:     pointer.From(file.Content),
										VirtualPath: pointer.From(file.VirtualPath),
									})
								}
							}

							var protectedFiles []ConfigureFile
							if prop.ProtectedFiles != nil {
								for _, file := range *prop.ProtectedFiles {
									protectedFiles = append(protectedFiles, ConfigureFile{
										Content:     pointer.From(file.Content),
										VirtualPath: pointer.From(file.VirtualPath),
									})
								}
							}

							output.Configuration = []Configuration{
								{
									ConfigureFile: files,
									ProtectedFile: protectedFiles,
									PackageData:   pointer.From(pointer.From(prop.Package).Data),
									RootFile:      pointer.From(prop.RootFile),
								},
							}
						}
					}
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
			if meta.ResourceData.HasChange("logging_storage_account") && len(model.LoggingStorageAccount) > 0 {
				req.Properties.Logging = &nginxdeployment.NginxLogging{
					StorageAccount: &nginxdeployment.NginxStorageAccount{
						AccountName:   pointer.FromString(model.LoggingStorageAccount[0].Name),
						ContainerName: pointer.FromString(model.LoggingStorageAccount[0].ContainerName),
					},
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
					PreferredEmail: pointer.FromString(model.Email),
				}
			}

			if meta.ResourceData.HasChange("automatic_upgrade_channel") {
				req.Properties.AutoUpgradeProfile = &nginxdeployment.AutoUpgradeProfile{
					UpgradeChannel: model.UpgradeChannel,
				}
			}

			if strings.HasPrefix(model.Sku, "basic") && req.Properties.ScalingProperties != nil {
				return fmt.Errorf("basic SKUs are incompatible with `capacity` or `auto_scale_profiles`")
			}

			if err := client.DeploymentsUpdateThenPoll(ctx, *id, req); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			if !features.FourPointOhBeta() {
				if meta.ResourceData.HasChange("configuration") {
					configID := nginxconfiguration.NewConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.NginxDeploymentName, defaultConfigurationName)

					configProp := expandConfiguration(model.Configuration[0])
					if err := meta.Client.Nginx.NginxConfiguration.ConfigurationsCreateOrUpdateThenPoll(ctx, configID, configProp); err != nil {
						return fmt.Errorf("update default configuration of %q: %v", configID, err)
					}
				}
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

func expandConfiguration(model Configuration) nginxconfiguration.NginxConfiguration {
	result := nginxconfiguration.NginxConfiguration{
		Properties: &nginxconfiguration.NginxConfigurationProperties{},
	}

	if len(model.ConfigureFile) > 0 {
		var files []nginxconfiguration.NginxConfigurationFile
		for _, file := range model.ConfigureFile {
			files = append(files, nginxconfiguration.NginxConfigurationFile{
				Content:     pointer.To(file.Content),
				VirtualPath: pointer.To(file.VirtualPath),
			})
		}
		result.Properties.Files = &files
	}

	if len(model.ProtectedFile) > 0 {
		var files []nginxconfiguration.NginxConfigurationFile
		for _, file := range model.ProtectedFile {
			files = append(files, nginxconfiguration.NginxConfigurationFile{
				Content:     pointer.To(file.Content),
				VirtualPath: pointer.To(file.VirtualPath),
			})
		}
		result.Properties.ProtectedFiles = &files
	}

	if model.PackageData != "" {
		result.Properties.Package = &nginxconfiguration.NginxConfigurationPackage{
			Data: pointer.To(model.PackageData),
		}
	}

	if model.RootFile != "" {
		result.Properties.RootFile = pointer.To(model.RootFile)
	}

	return result
}
