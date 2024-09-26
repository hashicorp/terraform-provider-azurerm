// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nginx

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-06-01-preview/nginxdeployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DeploymentDataSourceModel struct {
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
	Tags                   map[string]string                          `tfschema:"tags"`
}

type DeploymentDataSource struct{}

var _ sdk.DataSource = DeploymentDataSource{}

func (m DeploymentDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (m DeploymentDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"nginx_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

		"sku": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"managed_resource_group": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"location": commonschema.LocationComputed(),

		"capacity": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"auto_scale_profile": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"min_capacity": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"max_capacity": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"diagnose_support_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"email": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"logging_storage_account": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"container_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"frontend_public": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_address": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"frontend_private": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_address": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"allocation_method": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"network_interface": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"automatic_upgrade_channel": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (m DeploymentDataSource) ModelObject() interface{} {
	return &DeploymentDataSourceModel{}
}

func (m DeploymentDataSource) ResourceType() string {
	return "azurerm_nginx_deployment"
}

func (m DeploymentDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Nginx.NginxDeployment
			subscriptionId := metadata.Client.Account.SubscriptionId
			var model DeploymentDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}
			id := nginxdeployment.NewNginxDeploymentID(subscriptionId, model.ResourceGroupName, model.Name)
			result, err := client.DeploymentsGet(ctx, id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			output := DeploymentDataSourceModel{
				Name:              id.NginxDeploymentName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := result.Model; model != nil {
				output.Location = pointer.ToString(model.Location)
				if tags := model.Tags; tags != nil {
					output.Tags = pointer.ToMapOfStringStrings(model.Tags)
				}
				if model.Sku != nil {
					output.Sku = model.Sku.Name
				}
				flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %v", err)
				}
				output.Identity = *flattenedIdentity
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
				}
			}

			metadata.SetID(id)
			return metadata.Encode(&output)
		},
	}
}
