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
			Description:  "The name of the Container Group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"ip_address_type": {
			Description: "The IP address type of the container. `Public`, `Private` or `None`. If set to `Private`, `subnet_ids` also needs to be set.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     string(containerinstance.ContainerGroupIPAddressTypePublic),
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(containerinstance.ContainerGroupIPAddressTypePublic),
				string(containerinstance.ContainerGroupIPAddressTypePrivate),
				"None",
			}, false),
		},

		"os_type": {
			Description:  "The OS for the Container Group.",
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
						Description:  "The address to use to connect to the registry without protocol (`https`/`http`). For example: `myacr.acr.io`.",
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
						Description:  "The identity ID for the private registry.",
					},

					"username": {
						Description:  "The username with which to connect to the registry.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"password": {
						Description:  "The password with which to connect to the registry.",
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
			Description: "The subnet resource IDs for a Container Group.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			ForceNew:    true,
			MaxItems:    1,
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
			Description:  "The SKU of the Container Group.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(containerinstance.ContainerGroupSkuStandard),
			ValidateFunc: validation.StringInSlice(containerinstance.PossibleValuesForContainerGroupSku(), false),
		},

		"restart_policy": {
			Description:  "Restart policy for the Container Group.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(containerinstance.ContainerGroupRestartPolicyAlways),
			ValidateFunc: validation.StringInSlice(containerinstance.PossibleValuesForContainerGroupRestartPolicy(), false),
		},

		"dns_name_label": {
			Description: "The DNS label/name for the Container Group's IP.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
		},

		"dns_name_label_reuse_policy": {
			Description:  "The value representing the security enum. `Noreuse`, `ResourceGroupReuse`, `SubscriptionReuse`, `TenantReuse` or `Unsecure`.",
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
						Description:  "The port number the container will expose.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsPortNumber,
					},

					"protocol": {
						Description:  "The network protocol associated with port.",
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
			Description: "The definition of an init container that is part of the group.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			ForceNew:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the Container.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"image": {
						Description:  "The container image name.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"environment_variables": {
						Description: "A list of environment variables to be set on the container. Specified as a map of name/value pairs.",
						Type:        pluginsdk.TypeMap,
						ForceNew:    true,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"secure_environment_variables": {
						Description: "A list of sensitive environment variables to be set on the container. Specified as a map of name/value pairs.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						ForceNew:    true,
						Sensitive:   true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"commands": {
						Description: "A list of commands which should be run on the container.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						ForceNew:    true,
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
			Description: "The definition of a container that is part of the group.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			ForceNew:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the Container.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"image": {
						Description:  "The container image name.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"cpu": {
						Description: "The required number of CPU cores of the containers.",
						Type:        pluginsdk.TypeFloat,
						Required:    true,
						ForceNew:    true,
					},

					"memory": {
						Description: "The required memory of the containers in GB.",
						Type:        pluginsdk.TypeFloat,
						Required:    true,
						ForceNew:    true,
					},

					"cpu_limit": {
						Description:  "The upper limit of the number of CPU cores of the containers.",
						Type:         pluginsdk.TypeFloat,
						Optional:     true,
						ValidateFunc: validation.FloatAtLeast(0.0),
					},

					"memory_limit": {
						Description:  "The upper limit of the memory of the containers in GB.",
						Type:         pluginsdk.TypeFloat,
						Optional:     true,
						ValidateFunc: validation.FloatAtLeast(0.0),
					},

					"ports": {
						Description: "A set of public ports for the container.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						ForceNew:    true,
						Set:         resourceContainerGroupPortsHash,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"port": {
									Description:  "The port number the container will expose.",
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.IsPortNumber,
								},

								"protocol": {
									Description: "The network protocol associated with port.",
									Type:        pluginsdk.TypeString,
									Optional:    true,
									ForceNew:    true,
									Default:     string(containerinstance.ContainerGroupNetworkProtocolTCP),
									ValidateFunc: validation.StringInSlice([]string{
										string(containerinstance.ContainerGroupNetworkProtocolTCP),
										string(containerinstance.ContainerNetworkProtocolUDP),
									}, false),
								},
							},
						},
					},

					"environment_variables": {
						Description: "A list of environment variables to be set on the container. Specified as a map of name/value pairs.",
						Type:        pluginsdk.TypeMap,
						ForceNew:    true,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"secure_environment_variables": {
						Description: "A list of sensitive environment variables to be set on the container. Specified as a map of name/value pairs.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						ForceNew:    true,
						Sensitive:   true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"commands": {
						Description: "A list of commands which should be run on the container.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						ForceNew:    true,
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
									Description:  "The Workspace ID of the Log Analytics Workspace.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IsUUID,
								},

								"workspace_key": {
									Description:  "The Workspace Key of the Log Analytics Workspace.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									Sensitive:    true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"log_type": {
									Description:  "The log type which should be used.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringInSlice(containerinstance.PossibleValuesForLogAnalyticsLogType(), false),
								},

								"metadata": {
									Description: "Any metadata required for Log Analytics.",
									Type:        pluginsdk.TypeMap,
									Optional:    true,
									ForceNew:    true,
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
			Description: "The IP address allocated to the Container Group.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"fqdn": {
			Description: "The FQDN of the Container Group derived from `dns_name_label`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"dns_config": {
			Optional: true,
			MaxItems: 1,
			ForceNew: true,
			Type:     pluginsdk.TypeList,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"nameservers": {
						Description: "A list of nameservers the containers will search out to resolve requests.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						ForceNew:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"search_domains": {
						Description: "A list of search domains that DNS requests will search along.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						ForceNew:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"options": {
						Description: "A list of [resolver configuration options](https://man7.org/linux/man-pages/man5/resolv.conf.5.html).",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						ForceNew:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},

		"key_vault_key_id": {
			Description:  "The Key Vault Key URI for CMK encryption.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
		},

		"key_vault_user_assigned_identity_id": {
			Description:  "The User Assigned Identity that has access to the Key Vault Key. If not specified, the RP principal named \"Azure Container Instance Service\" will be used instead. Make sure the identity has the proper `key_permissions` set, at least with `Get`, `UnwrapKey`, `WrapKey` and `GetRotationPolicy`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"priority": {
			Description:  "The priority of the Container Group.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(containerinstance.PossibleValuesForContainerGroupPriority(), false),
		},
	}
}

func containerVolumeSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Description: "The definition of a volume mount for this container.",
		Type:        pluginsdk.TypeList,
		Optional:    true,
		ForceNew:    true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Description:  "The name of the volume mount.",
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"mount_path": {
					Description:  "The path on which this volume is to be mounted.",
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"read_only": {
					Description: "Whether the volume is to be mounted as read only.",
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					ForceNew:    true,
					Default:     false,
				},

				"share_name": {
					Description:  "The Azure Storage Share that is to be mounted as a volume.",
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"storage_account_name": {
					Description:  "The Azure Storage Account from which the volume is to be mounted.",
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
					Description: "Boolean as to whether the mounted volume should be an empty directory.",
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					ForceNew:    true,
					Default:     false,
				},

				"git_repo": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"url": {
								Description: "The Git repository to be cloned.",
								Type:        pluginsdk.TypeString,
								Required:    true,
								ForceNew:    true,
							},

							"directory": {
								Description: "The directory into which the repository should be cloned.",
								Type:        pluginsdk.TypeString,
								Optional:    true,
								ForceNew:    true,
							},

							"revision": {
								Description: "The commit hash of the revision to be cloned. If unspecified, the HEAD revision is cloned.",
								Type:        pluginsdk.TypeString,
								Optional:    true,
								ForceNew:    true,
							},
						},
					},
				},

				"secret": {
					Description: "A map of secrets that will be mounted as files in the volume.",
					Type:        pluginsdk.TypeMap,
					ForceNew:    true,
					Optional:    true,
					Sensitive:   true,
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
		Description: "The definition of the security context for this container.",
		Type:        pluginsdk.TypeList,
		Optional:    true,
		ForceNew:    true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"privilege_enabled": {
					Description: "Whether the container's permission is elevated to privileged.",
					Type:        pluginsdk.TypeBool,
					ForceNew:    true,
					Required:    true,
				},
			},
		},
	}
}
