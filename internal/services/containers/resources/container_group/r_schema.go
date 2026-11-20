// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_group

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerinstance/2025-09-01/containerinstance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func containerGroupSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"ip_address_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(containerinstance.ContainerGroupIPAddressTypePublic),
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(containerinstance.ContainerGroupIPAddressTypePublic),
				string(containerinstance.ContainerGroupIPAddressTypePrivate),
				"None",
			}, false),
		},

		"os_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(containerinstance.PossibleValuesForOperatingSystemTypes(), false),
		},

		"image_registry_credential": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"server": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"user_assigned_identity_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
						Description:  "The User Assigned Identity to use for Container Registry access.",
					},

					"username": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"password": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"network_profile_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			// TODO: 6.0 - remove this, was meant to be removed in 4.0...
			Deprecated: "the 'network_profile_id' has been removed from the latest versions of the container instance API and has been deprecated. It no longer functions and will be removed from the 6.0 AzureRM provider. Please use the 'subnet_ids' field instead",
		},

		// lintignore:S018
		"subnet_ids": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: commonids.ValidateSubnetID,
			},
			Set:           pluginsdk.HashString,
			ConflictsWith: []string{"dns_name_label"},
		},

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"tags": commonschema.Tags(),

		"sku": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(containerinstance.ContainerGroupSkuStandard),
			ValidateFunc: validation.StringInSlice(containerinstance.PossibleValuesForContainerGroupSku(), false),
		},

		"restart_policy": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(containerinstance.ContainerGroupRestartPolicyAlways),
			ValidateFunc: validation.StringInSlice(containerinstance.PossibleValuesForContainerGroupRestartPolicy(), false),
		},

		"dns_name_label": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"dns_name_label_reuse_policy": {
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Optional:     true,
			Default:      string(containerinstance.DnsNameLabelReusePolicyUnsecure),
			ValidateFunc: validation.StringInSlice(containerinstance.PossibleValuesForDnsNameLabelReusePolicy(), false),
		},

		"exposed_port": {
			Type:       pluginsdk.TypeSet,
			Optional:   true,
			Computed:   true, // azignore:AZS007 - pre-existing violation
			ForceNew:   true,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Set:        resourceContainerGroupPortsHash,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"port": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsPortNumber,
					},

					"protocol": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						Default:      string(containerinstance.ContainerGroupNetworkProtocolTCP),
						ValidateFunc: validation.StringInSlice(containerinstance.PossibleValuesForContainerGroupNetworkProtocol(), false),
					},
				},
			},
		},

		"init_container": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"image": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"environment_variables": {
						Type:     pluginsdk.TypeMap,
						ForceNew: true,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"secure_environment_variables": {
						Type:      pluginsdk.TypeMap,
						Optional:  true,
						ForceNew:  true,
						Sensitive: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"commands": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"volume": containerVolumeSchema(),

					"security": containerSecurityContextSchema(),
				},
			},
		},
		"container": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"image": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"cpu": {
						Type:     pluginsdk.TypeFloat,
						Required: true,
						ForceNew: true,
					},

					"memory": {
						Type:     pluginsdk.TypeFloat,
						Required: true,
						ForceNew: true,
					},

					"config_map": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"cpu_limit": {
						Type:         pluginsdk.TypeFloat,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.FloatAtLeast(0.0),
					},

					"memory_limit": {
						Type:         pluginsdk.TypeFloat,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.FloatAtLeast(0.0),
					},

					"ports": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						ForceNew: true,
						Set:      resourceContainerGroupPortsHash,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"port": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.IsPortNumber,
								},

								"protocol": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ForceNew: true,
									Default:  string(containerinstance.ContainerGroupNetworkProtocolTCP),
									ValidateFunc: validation.StringInSlice([]string{
										string(containerinstance.ContainerGroupNetworkProtocolTCP),
										string(containerinstance.ContainerNetworkProtocolUDP),
									}, false),
								},
							},
						},
					},

					"environment_variables": {
						Type:     pluginsdk.TypeMap,
						ForceNew: true,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"secure_environment_variables": {
						Type:      pluginsdk.TypeMap,
						Optional:  true,
						ForceNew:  true,
						Sensitive: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"commands": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"volume": containerVolumeSchema(),

					"security": containerSecurityContextSchema(),

					"liveness_probe": SchemaContainerGroupProbe(),

					"readiness_probe": SchemaContainerGroupProbe(),
				},
			},
		},

		"diagnostics": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"log_analytics": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"workspace_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IsUUID,
								},

								"workspace_key": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									Sensitive:    true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"log_type": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringInSlice(containerinstance.PossibleValuesForLogAnalyticsLogType(), false),
								},

								"metadata": {
									Type:     pluginsdk.TypeMap,
									Optional: true,
									ForceNew: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},

		"ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"dns_config": {
			Optional: true,
			MaxItems: 1,
			ForceNew: true,
			Type:     pluginsdk.TypeList,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"nameservers": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"search_domains": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"options": {
						Type:     pluginsdk.TypeSet,
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

		"key_vault_key_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
		},

		"key_vault_user_assigned_identity_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"priority": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(containerinstance.PossibleValuesForContainerGroupPriority(), false),
		},
	}
}

func containerVolumeSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"mount_path": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"read_only": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
					Default:  false,
				},

				"share_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"storage_account_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"storage_account_key": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"empty_dir": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
					Default:  false,
				},

				"git_repo": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"url": {
								Type:     pluginsdk.TypeString,
								Required: true,
								ForceNew: true,
							},

							"directory": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								ForceNew: true,
							},

							"revision": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								ForceNew: true,
							},
						},
					},
				},

				"secret": {
					Type:      pluginsdk.TypeMap,
					ForceNew:  true,
					Optional:  true,
					Sensitive: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},
			},
		},
	}
}

func containerSecurityContextSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"privilege_enabled": {
					Type:     pluginsdk.TypeBool,
					ForceNew: true,
					Required: true,
				},
			},
		},
	}
}
