// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerinstance/2023-05-01/containerinstance"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceContainerGroup() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceContainerGroupCreate,
		Read:   resourceContainerGroupRead,
		Delete: resourceContainerGroupDelete,
		Update: resourceContainerGroupUpdate,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := containerinstance.ParseContainerGroupID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(containerinstance.OperatingSystemTypesWindows),
					string(containerinstance.OperatingSystemTypesLinux),
				}, false),
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
				Type:       pluginsdk.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "the 'network_profile_id' has been removed from the latest versions of the container instance API and has been deprecated. It no longer functions and will be removed from the 4.0 AzureRM provider. Please use the 'subnet_ids' field instead",
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
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(containerinstance.ContainerGroupRestartPolicyAlways),
				ValidateFunc: validation.StringInSlice([]string{
					string(containerinstance.ContainerGroupRestartPolicyAlways),
					string(containerinstance.ContainerGroupRestartPolicyNever),
					string(containerinstance.ContainerGroupRestartPolicyOnFailure),
				}, false),
			},

			"dns_name_label": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"dns_name_label_reuse_policy": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(containerinstance.DnsNameLabelReusePolicyUnsecure),
				ValidateFunc: validation.StringInSlice([]string{
					string(containerinstance.DnsNameLabelReusePolicyNoreuse),
					string(containerinstance.DnsNameLabelReusePolicyResourceGroupReuse),
					string(containerinstance.DnsNameLabelReusePolicySubscriptionReuse),
					string(containerinstance.DnsNameLabelReusePolicyTenantReuse),
					string(containerinstance.DnsNameLabelReusePolicyUnsecure),
				}, false),
			},

			"exposed_port": {
				Type:       pluginsdk.TypeSet,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Set:        resourceContainerGroupPortsHash,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"port": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validate.PortNumber,
						},

						"protocol": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  string(containerinstance.ContainerGroupNetworkProtocolTCP),
							ValidateFunc: validation.StringInSlice([]string{
								string(containerinstance.ContainerGroupNetworkProtocolTCP),
								string(containerinstance.ContainerGroupNetworkProtocolUDP),
							}, false),
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
							Computed: true,
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

						"cpu_limit": {
							Type:         pluginsdk.TypeFloat,
							Optional:     true,
							ValidateFunc: validation.FloatAtLeast(0.0),
						},

						"memory_limit": {
							Type:         pluginsdk.TypeFloat,
							Optional:     true,
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
										ValidateFunc: validate.PortNumber,
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
							Computed: true,
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
										Type:     pluginsdk.TypeString,
										Optional: true,
										ForceNew: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(containerinstance.LogAnalyticsLogTypeContainerInsights),
											string(containerinstance.LogAnalyticsLogTypeContainerInstanceLogs),
										}, false),
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
				ValidateFunc: keyVaultValidate.NestedItemId,
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
		},
		CustomizeDiff: func(ctx context.Context, d *pluginsdk.ResourceDiff, i interface{}) error {
			if p := d.Get("priority").(string); p == string(containerinstance.ContainerGroupPrioritySpot) {
				if d.Get("ip_address_type").(string) != "None" {
					return fmt.Errorf("`ip_address_type` has to be `None` when `priority` is set to `Spot`")
				}
			}
			return nil
		},
	}

	if !features.FourPointOhBeta() {
		resource.Schema["container"].Elem.(*pluginsdk.Resource).Schema["gpu"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeList,
			Optional:   true,
			MaxItems:   1,
			ForceNew:   true,
			Deprecated: "The `gpu` block has been deprecated since K80 and P100 GPU Skus have been retired and remaining GPU resources are not fully supported and not appropriate for production workloads. This block will be removed in v4.0 of the AzureRM provider.",
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"count": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						ForceNew: true,
						ValidateFunc: validation.IntInSlice([]int{
							1,
							2,
							4,
						}),
					},

					"sku": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							"K80",
							"P100",
							"V100",
						}, false),
					},
				},
			},
		}
		resource.Schema["container"].Elem.(*pluginsdk.Resource).Schema["gpu_limit"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeList,
			Optional:   true,
			MaxItems:   1,
			Deprecated: "The `gpu_limit` block has been deprecated since K80 and P100 GPU Skus have been retired and remaining GPU resources are not fully supported and not appropriate for production workloads. This block will be removed in v4.0 of the AzureRM provider.",
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"count": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"sku": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"K80",
							"P100",
							"V100",
						}, false),
					},
				},
			},
		}
	}

	return resource
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

func resourceContainerGroupCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerInstanceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := containerinstance.NewContainerGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.ContainerGroupsGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_container_group", id.ID())
	}

	location := location.Normalize(d.Get("location").(string))
	OSType := d.Get("os_type").(string)
	IPAddressType := d.Get("ip_address_type").(string)
	restartPolicy := containerinstance.ContainerGroupRestartPolicy(d.Get("restart_policy").(string))
	diagnosticsRaw := d.Get("diagnostics").([]interface{})
	diagnostics := expandContainerGroupDiagnostics(diagnosticsRaw)
	dnsConfig := d.Get("dns_config").([]interface{})
	addedEmptyDirs := map[string]bool{}
	subnets, err := expandContainerGroupSubnets(d.Get("subnet_ids").(*pluginsdk.Set).List())
	if err != nil {
		return err
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*pluginsdk.Set).List())
	initContainers, initContainerVolumes, err := expandContainerGroupInitContainers(d, addedEmptyDirs)
	if err != nil {
		return err
	}

	containers, containerGroupPorts, containerVolumes, err := expandContainerGroupContainers(d, addedEmptyDirs)
	if err != nil {
		return err
	}
	var containerGroupVolumes []containerinstance.Volume
	if initContainerVolumes != nil {
		containerGroupVolumes = initContainerVolumes
	}
	if containerGroupVolumes != nil {
		containerGroupVolumes = append(containerGroupVolumes, containerVolumes...)
	}

	containerGroup := containerinstance.ContainerGroup{
		Name:     pointer.FromString(id.ContainerGroupName),
		Location: &location,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: containerinstance.ContainerGroupPropertiesProperties{
			Sku:                      pointer.To(containerinstance.ContainerGroupSku(d.Get("sku").(string))),
			InitContainers:           initContainers,
			Containers:               containers,
			Diagnostics:              diagnostics,
			RestartPolicy:            &restartPolicy,
			OsType:                   containerinstance.OperatingSystemTypes(OSType),
			Volumes:                  &containerGroupVolumes,
			ImageRegistryCredentials: expandContainerImageRegistryCredentials(d),
			DnsConfig:                expandContainerGroupDnsConfig(dnsConfig),
			SubnetIds:                subnets,
		},
		Zones: &zones,
	}

	expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	containerGroup.Identity = expandedIdentity

	if IPAddressType != "None" {
		containerGroup.Properties.IPAddress = &containerinstance.IPAddress{
			Ports: containerGroupPorts,
			Type:  containerinstance.ContainerGroupIPAddressType(IPAddressType),
		}

		if dnsNameLabel := d.Get("dns_name_label").(string); dnsNameLabel != "" {
			containerGroup.Properties.IPAddress.DnsNameLabel = &dnsNameLabel
		}
		if dnsNameLabelReusePolicy := d.Get("dns_name_label_reuse_policy").(string); dnsNameLabelReusePolicy != "" {
			containerGroup.Properties.IPAddress.AutoGeneratedDomainNameLabelScope = (*containerinstance.DnsNameLabelReusePolicy)(&dnsNameLabelReusePolicy)
		}
	}

	if keyVaultKeyId := d.Get("key_vault_key_id").(string); keyVaultKeyId != "" {
		keyId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(keyVaultKeyId)
		if err != nil {
			return fmt.Errorf("parsing Key Vault Key ID: %+v", err)
		}
		containerGroup.Properties.EncryptionProperties = &containerinstance.EncryptionProperties{
			VaultBaseURL: keyId.KeyVaultBaseUrl,
			KeyName:      keyId.Name,
			KeyVersion:   keyId.Version,
		}

		if keyVaultUAI := d.Get("key_vault_user_assigned_identity_id").(string); keyVaultUAI != "" {
			containerGroup.Properties.EncryptionProperties.Identity = &keyVaultUAI
		}
	}

	if priority := d.Get("priority").(string); priority != "" {
		containerGroup.Properties.Priority = pointer.To(containerinstance.ContainerGroupPriority(priority))
	}

	// Avoid parallel provisioning if "subnet_ids" are given.
	if subnets != nil && len(*subnets) != 0 {
		for _, item := range *subnets {
			subnet, err := commonids.ParseSubnetID(item.Id)
			if err != nil {
				return fmt.Errorf(`parsing subnet id %q: %v`, item.Id, err)
			}

			locks.ByID(subnet.ID())
			defer locks.UnlockByID(subnet.ID())
		}
	}

	if err := client.ContainerGroupsCreateOrUpdateThenPoll(ctx, id, containerGroup); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceContainerGroupRead(d, meta)
}

func resourceContainerGroupUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerInstanceClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := containerinstance.ParseContainerGroupID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.ContainerGroupsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("reading %s: %v", id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("reading %s: `model` was nil", id)
	}

	model := *existing.Model

	if d.HasChange("identity") {
		expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		model.Identity = expandedIdentity

		// As API doesn't return the value of StorageAccountKey, so it has to get the value from tf config and set it to request payload. Otherwise, the Update API call would fail
		addedEmptyDirs := map[string]bool{}
		_, initContainerVolumes, err := expandContainerGroupInitContainers(d, addedEmptyDirs)
		if err != nil {
			return err
		}
		_, _, containerVolumes, err := expandContainerGroupContainers(d, addedEmptyDirs)
		if err != nil {
			return err
		}
		var containerGroupVolumes []containerinstance.Volume
		if initContainerVolumes != nil {
			containerGroupVolumes = initContainerVolumes
		}
		if containerGroupVolumes != nil {
			containerGroupVolumes = append(containerGroupVolumes, containerVolumes...)
		}
		model.Properties.Volumes = pointer.To(containerGroupVolumes)

		// As Update API doesn't support to update identity, so it has to use CreateOrUpdate API to update identity
		if err := client.ContainerGroupsCreateOrUpdateThenPoll(ctx, *id, model); err != nil {
			return fmt.Errorf("updating %s: %+v", *id, err)
		}
	}

	if d.HasChange("tags") {
		updateParameters := containerinstance.Resource{
			Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
		}

		// As CreateOrUpdate API doesn't support to update tags, so it has to use Update API to update tags
		if _, err := client.ContainerGroupsUpdate(ctx, *id, updateParameters); err != nil {
			return fmt.Errorf("updating tags %s: %+v", *id, err)
		}
	}

	return resourceContainerGroupRead(d, meta)
}

func resourceContainerGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerInstanceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := containerinstance.ParseContainerGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ContainerGroupsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", id.ContainerGroupName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		identity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}

		d.Set("zones", zones.FlattenUntyped(model.Zones))

		props := model.Properties

		var sku string
		if v := props.Sku; v != nil {
			sku = string(*v)
		}
		d.Set("sku", sku)

		var priority string
		if v := props.Priority; v != nil {
			priority = string(*v)
		}
		d.Set("priority", priority)

		containerConfigs := flattenContainerGroupContainers(d, &props.Containers, props.Volumes)
		if err := d.Set("container", containerConfigs); err != nil {
			return fmt.Errorf("setting `container`: %+v", err)
		}
		initContainerConfigs := flattenContainerGroupInitContainers(d, props.InitContainers, props.Volumes)
		if err := d.Set("init_container", initContainerConfigs); err != nil {
			return fmt.Errorf("setting `init_container`: %+v", err)
		}

		if err := d.Set("image_registry_credential", flattenContainerImageRegistryCredentials(d, props.ImageRegistryCredentials)); err != nil {
			return fmt.Errorf("setting `image_registry_credential`: %+v", err)
		}

		if address := props.IPAddress; address != nil {
			d.Set("ip_address_type", address.Type)
			d.Set("ip_address", address.IP)
			exposedPorts := make([]interface{}, len(address.Ports))
			for i := range address.Ports {
				exposedPorts[i] = (address.Ports)[i]
			}
			d.Set("exposed_port", flattenPorts(exposedPorts))
			d.Set("dns_name_label", address.DnsNameLabel)
			d.Set("fqdn", address.Fqdn)

			if address.AutoGeneratedDomainNameLabelScope != nil {
				d.Set("dns_name_label_reuse_policy", string(*address.AutoGeneratedDomainNameLabelScope))
			} else {
				d.Set("dns_name_label_reuse_policy", containerinstance.DnsNameLabelReusePolicyUnsecure)
			}
		} else {
			d.Set("dns_name_label_reuse_policy", pointer.FromString(string(containerinstance.DnsNameLabelReusePolicyUnsecure)))
		}

		restartPolicy := ""
		if props.RestartPolicy != nil {
			restartPolicy = string(*props.RestartPolicy)
		}
		d.Set("restart_policy", restartPolicy)

		d.Set("os_type", string(props.OsType))
		d.Set("dns_config", flattenContainerGroupDnsConfig(props.DnsConfig))

		if err := d.Set("diagnostics", flattenContainerGroupDiagnostics(d, props.Diagnostics)); err != nil {
			return fmt.Errorf("setting `diagnostics`: %+v", err)
		}

		subnets, err := flattenContainerGroupSubnets(props.SubnetIds)
		if err != nil {
			return err
		}
		if err := d.Set("subnet_ids", subnets); err != nil {
			return fmt.Errorf("setting `subnet_ids`: %+v", err)
		}

		if kvProps := props.EncryptionProperties; kvProps != nil {
			var keyVaultUri, keyName, keyVersion string
			if kvProps.VaultBaseURL != "" {
				keyVaultUri = kvProps.VaultBaseURL
			} else {
				return fmt.Errorf("empty value returned for Key Vault URI")
			}
			if kvProps.KeyName != "" {
				keyName = kvProps.KeyName
			} else {
				return fmt.Errorf("empty value returned for Key Vault Key Name")
			}
			keyVersion = kvProps.KeyVersion
			keyId, err := keyVaultParse.NewNestedItemID(keyVaultUri, keyVaultParse.NestedItemTypeKey, keyName, keyVersion)
			if err != nil {
				return err
			}
			d.Set("key_vault_key_id", keyId.ID())
			d.Set("key_vault_user_assigned_identity_id", pointer.From(kvProps.Identity))
		}
	}

	return nil
}

func flattenPorts(ports []interface{}) *pluginsdk.Set {
	if len(ports) > 0 {
		flatPorts := make([]interface{}, 0)
		for _, p := range ports {
			port := make(map[string]interface{})
			switch t := p.(type) {
			case containerinstance.Port:
				port["port"] = int(t.Port)
				proto := ""
				if t.Protocol != nil {
					proto = string(*t.Protocol)
				}
				port["protocol"] = proto
			case containerinstance.ContainerPort:
				port["port"] = int(t.Port)
				proto := ""
				if t.Protocol != nil {
					proto = string(*t.Protocol)
				}
				port["protocol"] = proto
			}
			flatPorts = append(flatPorts, port)
		}
		return pluginsdk.NewSet(resourceContainerGroupPortsHash, flatPorts)
	}
	return pluginsdk.NewSet(resourceContainerGroupPortsHash, make([]interface{}, 0))
}

func resourceContainerGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerInstanceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := containerinstance.ParseContainerGroupID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.ContainerGroupsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			// already deleted
			return nil
		}

		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	if model := existing.Model; model != nil {
		props := model.Properties
		if subnetIDs := props.SubnetIds; subnetIDs != nil && len(*subnetIDs) != 0 {
			// Avoid parallel deletion if "subnet_ids" are given.
			for _, item := range *subnetIDs {
				subnet, err := commonids.ParseSubnetID(item.Id)
				if err != nil {
					return fmt.Errorf(`parsing subnet id %q: %v`, item.Id, err)
				}

				locks.ByID(subnet.ID())
				defer locks.UnlockByID(subnet.ID())
			}
		}
	}

	if err := client.ContainerGroupsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	return nil
}

func expandContainerGroupInitContainers(d *pluginsdk.ResourceData, addedEmptyDirs map[string]bool) (*[]containerinstance.InitContainerDefinition, []containerinstance.Volume, error) {
	containersConfig := d.Get("init_container").([]interface{})
	containers := make([]containerinstance.InitContainerDefinition, 0)
	containerGroupVolumes := make([]containerinstance.Volume, 0)
	for _, containerConfig := range containersConfig {
		data := containerConfig.(map[string]interface{})

		name := data["name"].(string)
		image := data["image"].(string)

		container := containerinstance.InitContainerDefinition{
			Name: name,
			Properties: containerinstance.InitContainerPropertiesDefinition{
				Image:           pointer.FromString(image),
				SecurityContext: expandContainerSecurityContext(data["security"].([]interface{})),
			},
		}

		// Set both sensitive and non-secure environment variables
		var envVars *[]containerinstance.EnvironmentVariable
		var secEnvVars *[]containerinstance.EnvironmentVariable

		// Expand environment_variables into slice
		if v, ok := data["environment_variables"]; ok {
			envVars = expandContainerEnvironmentVariables(v, false)
		}

		// Expand secure_environment_variables into slice
		if v, ok := data["secure_environment_variables"]; ok {
			secEnvVars = expandContainerEnvironmentVariables(v, true)
		}

		// Combine environment variable slices
		*envVars = append(*envVars, *secEnvVars...)

		// Set both secure and non secure environment variables
		container.Properties.EnvironmentVariables = envVars

		if v, ok := data["commands"]; ok {
			c := v.([]interface{})
			command := make([]string, 0)
			for _, v := range c {
				command = append(command, v.(string))
			}

			container.Properties.Command = &command
		}

		if v, ok := data["volume"]; ok {
			volumeMounts, _, err := expandSingleContainerVolume(v)
			if err != nil {
				return nil, nil, err
			}
			container.Properties.VolumeMounts = volumeMounts

			expandedContainerGroupVolumes, err := expandContainerVolume(v, addedEmptyDirs, containerGroupVolumes)
			if err != nil {
				return nil, nil, err
			}
			containerGroupVolumes = expandedContainerGroupVolumes
		}

		containers = append(containers, container)
	}

	return &containers, containerGroupVolumes, nil
}

func expandContainerSecurityContext(input []interface{}) *containerinstance.SecurityContextDefinition {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})

	output := &containerinstance.SecurityContextDefinition{
		Privileged: pointer.To(raw["privilege_enabled"].(bool)),
	}

	return output
}

func flattenContainerSecurityContext(input *containerinstance.SecurityContextDefinition) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var privileged bool
	if v := input.Privileged; v != nil {
		privileged = *v
	}

	return []interface{}{
		map[string]interface{}{
			"privilege_enabled": privileged,
		},
	}
}

func expandContainerGroupContainers(d *pluginsdk.ResourceData, addedEmptyDirs map[string]bool) ([]containerinstance.Container, []containerinstance.Port, []containerinstance.Volume, error) {
	containersConfig := d.Get("container").([]interface{})
	containers := make([]containerinstance.Container, 0)
	containerInstancePorts := make([]containerinstance.Port, 0)
	containerGroupPorts := make([]containerinstance.Port, 0)
	containerGroupVolumes := make([]containerinstance.Volume, 0)

	for _, containerConfig := range containersConfig {
		data := containerConfig.(map[string]interface{})

		name := data["name"].(string)
		image := data["image"].(string)
		cpu := data["cpu"].(float64)
		memory := data["memory"].(float64)

		container := containerinstance.Container{
			Name: name,
			Properties: containerinstance.ContainerProperties{
				Image: image,
				Resources: containerinstance.ResourceRequirements{
					Requests: containerinstance.ResourceRequests{
						MemoryInGB: memory,
						Cpu:        cpu,
					},
				},
				SecurityContext: expandContainerSecurityContext(data["security"].([]interface{})),
			},
		}

		cpuLimit := data["cpu_limit"].(float64)
		memLimit := data["memory_limit"].(float64)

		if !(cpuLimit == 0.0 && memLimit == 0.0) {
			limits := &containerinstance.ResourceLimits{}
			if cpuLimit != 0.0 {
				limits.Cpu = &cpuLimit
			}
			if memLimit != 0.0 {
				limits.MemoryInGB = &memLimit
			}

			container.Properties.Resources.Limits = limits
		}

		if !features.FourPointOhBeta() {
			if v, ok := data["gpu"]; ok {
				gpus := v.([]interface{})
				for _, gpuRaw := range gpus {
					if gpuRaw == nil {
						continue
					}
					v := gpuRaw.(map[string]interface{})
					gpuCount := int32(v["count"].(int))
					gpuSku := containerinstance.GpuSku(v["sku"].(string))

					gpus := containerinstance.GpuResource{
						Count: int64(gpuCount),
						Sku:   gpuSku,
					}
					container.Properties.Resources.Requests.Gpu = &gpus
				}
			}

			gpuLimit, ok := data["gpu_limit"].([]interface{})
			if ok && len(gpuLimit) == 1 && gpuLimit[0] != nil {
				if container.Properties.Resources.Limits == nil {
					container.Properties.Resources.Limits = &containerinstance.ResourceLimits{}
				}

				v := gpuLimit[0].(map[string]interface{})
				container.Properties.Resources.Limits.Gpu = &containerinstance.GpuResource{}
				if v := int64(v["count"].(int)); v != 0 {
					container.Properties.Resources.Limits.Gpu.Count = v
				}
				if v := containerinstance.GpuSku(v["sku"].(string)); v != "" {
					container.Properties.Resources.Limits.Gpu.Sku = v
				}

			}
		}

		if v, ok := data["ports"].(*pluginsdk.Set); ok && len(v.List()) > 0 {
			var ports []containerinstance.ContainerPort
			for _, v := range v.List() {
				portObj := v.(map[string]interface{})

				port := int64(portObj["port"].(int))
				proto := portObj["protocol"].(string)

				containerProtocol := containerinstance.ContainerNetworkProtocol(proto)
				ports = append(ports, containerinstance.ContainerPort{
					Port:     port,
					Protocol: &containerProtocol,
				})
				groupProtocol := containerinstance.ContainerGroupNetworkProtocol(proto)
				containerInstancePorts = append(containerInstancePorts, containerinstance.Port{
					Port:     port,
					Protocol: &groupProtocol,
				})
			}
			container.Properties.Ports = &ports
		}

		// Set both sensitive and non-secure environment variables
		var envVars *[]containerinstance.EnvironmentVariable
		var secEnvVars *[]containerinstance.EnvironmentVariable

		// Expand environment_variables into slice
		if v, ok := data["environment_variables"]; ok {
			envVars = expandContainerEnvironmentVariables(v, false)
		}

		// Expand secure_environment_variables into slice
		if v, ok := data["secure_environment_variables"]; ok {
			secEnvVars = expandContainerEnvironmentVariables(v, true)
		}

		// Combine environment variable slices
		*envVars = append(*envVars, *secEnvVars...)

		// Set both secure and non secure environment variables
		container.Properties.EnvironmentVariables = envVars

		if v, ok := data["commands"]; ok {
			c := v.([]interface{})
			command := make([]string, 0)
			for _, v := range c {
				command = append(command, v.(string))
			}

			container.Properties.Command = &command
		}

		if v, ok := data["volume"]; ok {
			volumeMounts, _, err := expandSingleContainerVolume(v)
			if err != nil {
				return nil, nil, nil, err
			}
			container.Properties.VolumeMounts = volumeMounts

			expandedContainerGroupVolumes, err := expandContainerVolume(v, addedEmptyDirs, containerGroupVolumes)
			if err != nil {
				return nil, nil, nil, err
			}
			containerGroupVolumes = expandedContainerGroupVolumes
		}

		if v, ok := data["liveness_probe"]; ok {
			container.Properties.LivenessProbe = expandContainerProbe(v)
		}

		if v, ok := data["readiness_probe"]; ok {
			container.Properties.ReadinessProbe = expandContainerProbe(v)
		}

		containers = append(containers, container)
	}

	// Determine ports to be exposed on the group level, based on exposed_ports
	// and on what ports have been exposed on individual containers.
	if v, ok := d.Get("exposed_port").(*pluginsdk.Set); ok && len(v.List()) > 0 {
		cgpMap := make(map[int64]map[containerinstance.ContainerGroupNetworkProtocol]bool)
		for _, p := range containerInstancePorts {
			if p.Protocol == nil {
				continue
			}
			protocol := *p.Protocol

			if val, ok := cgpMap[p.Port]; ok {
				val[protocol] = true
				cgpMap[p.Port] = val
			} else {
				protoMap := map[containerinstance.ContainerGroupNetworkProtocol]bool{protocol: true}
				cgpMap[p.Port] = protoMap
			}
		}

		for _, p := range v.List() {
			portConfig := p.(map[string]interface{})
			port := int64(portConfig["port"].(int))
			proto := portConfig["protocol"].(string)
			if !cgpMap[port][containerinstance.ContainerGroupNetworkProtocol(proto)] {
				return nil, nil, nil, fmt.Errorf("Port %d/%s is not exposed on any individual container in the container group.\n"+
					"An exposed_ports block contains %d/%s, but no individual container has a ports block with the same port "+
					"and protocol. Any ports exposed on the container group must also be exposed on an individual container.",
					port, proto, port, proto)
			}
			portProtocol := containerinstance.ContainerGroupNetworkProtocol(proto)
			containerGroupPorts = append(containerGroupPorts, containerinstance.Port{
				Port:     port,
				Protocol: &portProtocol,
			})
		}
	} else {
		containerGroupPorts = containerInstancePorts // remove in 3.0 of the provider
	}

	return containers, containerGroupPorts, containerGroupVolumes, nil
}

func expandContainerVolume(v interface{}, addedEmptyDirs map[string]bool, containerGroupVolumes []containerinstance.Volume) ([]containerinstance.Volume, error) {
	_, containerVolumes, err := expandSingleContainerVolume(v)
	if err != nil {
		return nil, err
	}
	if containerVolumes != nil {
		for _, cgVol := range *containerVolumes {
			if cgVol.EmptyDir != nil {
				if addedEmptyDirs[cgVol.Name] {
					// empty_dir-volumes are allowed to overlap across containers, in fact that is their primary purpose,
					// but the containerGroup must not declare same name of such volumes twice.
					continue
				}
				addedEmptyDirs[cgVol.Name] = true
			}
			containerGroupVolumes = append(containerGroupVolumes, cgVol)
		}
	}
	return containerGroupVolumes, nil
}

func expandContainerEnvironmentVariables(input interface{}, secure bool) *[]containerinstance.EnvironmentVariable {
	envVars := input.(map[string]interface{})
	output := make([]containerinstance.EnvironmentVariable, 0, len(envVars))

	if secure {
		for k, v := range envVars {
			ev := containerinstance.EnvironmentVariable{
				Name:        k,
				SecureValue: pointer.FromString(v.(string)),
			}

			output = append(output, ev)
		}
	} else {
		for k, v := range envVars {
			ev := containerinstance.EnvironmentVariable{
				Name:  k,
				Value: pointer.FromString(v.(string)),
			}

			output = append(output, ev)
		}
	}
	return &output
}

func expandContainerImageRegistryCredentials(d *pluginsdk.ResourceData) *[]containerinstance.ImageRegistryCredential {
	credsRaw := d.Get("image_registry_credential").([]interface{})
	if len(credsRaw) == 0 {
		return nil
	}

	output := make([]containerinstance.ImageRegistryCredential, 0, len(credsRaw))

	for _, c := range credsRaw {
		credConfig := c.(map[string]interface{})

		imageRegistryCredential := containerinstance.ImageRegistryCredential{}
		if v := credConfig["server"]; v != nil && v != "" {
			imageRegistryCredential.Server = v.(string)
		}
		if v := credConfig["username"]; v != nil && v != "" {
			imageRegistryCredential.Username = pointer.FromString(v.(string))
		}
		if v := credConfig["password"]; v != nil && v != "" {
			imageRegistryCredential.Password = pointer.FromString(v.(string))
		}
		if v := credConfig["user_assigned_identity_id"]; v != nil && v != "" {
			imageRegistryCredential.Identity = pointer.FromString(v.(string))
		}

		output = append(output, imageRegistryCredential)
	}

	return &output
}

func expandSingleContainerVolume(input interface{}) (*[]containerinstance.VolumeMount, *[]containerinstance.Volume, error) {
	volumesRaw := input.([]interface{})

	if len(volumesRaw) == 0 {
		return nil, nil, nil
	}

	volumeMounts := make([]containerinstance.VolumeMount, 0)
	containerGroupVolumes := make([]containerinstance.Volume, 0)

	for _, volumeRaw := range volumesRaw {
		volumeConfig := volumeRaw.(map[string]interface{})

		name := volumeConfig["name"].(string)
		mountPath := volumeConfig["mount_path"].(string)
		readOnly := volumeConfig["read_only"].(bool)
		emptyDir := volumeConfig["empty_dir"].(bool)
		shareName := volumeConfig["share_name"].(string)
		storageAccountName := volumeConfig["storage_account_name"].(string)
		storageAccountKey := volumeConfig["storage_account_key"].(string)

		vm := containerinstance.VolumeMount{
			Name:      name,
			MountPath: mountPath,
			ReadOnly:  pointer.FromBool(readOnly),
		}

		volumeMounts = append(volumeMounts, vm)

		cv := containerinstance.Volume{
			Name: name,
		}

		secret := expandSecrets(volumeConfig["secret"].(map[string]interface{}))

		gitRepoVolume := expandGitRepoVolume(volumeConfig["git_repo"].([]interface{}))

		switch {
		case emptyDir:
			if shareName != "" || storageAccountName != "" || storageAccountKey != "" || secret != nil || gitRepoVolume != nil {
				return nil, nil, fmt.Errorf("only one of `empty_dir` volume, `git_repo` volume, `secret` volume or storage account volume (`share_name`, `storage_account_name`, and `storage_account_key`) can be specified")
			}
			var m interface{} = map[string]string{}
			cv.EmptyDir = &m
		case gitRepoVolume != nil:
			if shareName != "" || storageAccountName != "" || storageAccountKey != "" || secret != nil {
				return nil, nil, fmt.Errorf("only one of `empty_dir` volume, `git_repo` volume, `secret` volume or storage account volume (`share_name`, `storage_account_name`, and `storage_account_key`) can be specified")
			}
			cv.GitRepo = gitRepoVolume
		case secret != nil:
			if shareName != "" || storageAccountName != "" || storageAccountKey != "" {
				return nil, nil, fmt.Errorf("only one of `empty_dir` volume, `git_repo` volume, `secret` volume or storage account volume (`share_name`, `storage_account_name`, and `storage_account_key`) can be specified")
			}
			cv.Secret = &secret
		default:
			if shareName == "" && storageAccountName == "" && storageAccountKey == "" {
				return nil, nil, fmt.Errorf("exactly one of `empty_dir` volume, `git_repo` volume, `secret` volume or storage account volume (`share_name`, `storage_account_name`, and `storage_account_key`) must be specified")
			} else if shareName == "" || storageAccountName == "" || storageAccountKey == "" {
				return nil, nil, fmt.Errorf("when using a storage account volume, all of `share_name`, `storage_account_name`, `storage_account_key` must be specified")
			}
			cv.AzureFile = &containerinstance.AzureFileVolume{
				ShareName:          shareName,
				ReadOnly:           pointer.FromBool(readOnly),
				StorageAccountName: storageAccountName,
				StorageAccountKey:  pointer.FromString(storageAccountKey),
			}
		}

		containerGroupVolumes = append(containerGroupVolumes, cv)
	}

	return &volumeMounts, &containerGroupVolumes, nil
}

func expandGitRepoVolume(input []interface{}) *containerinstance.GitRepoVolume {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	gitRepoVolume := &containerinstance.GitRepoVolume{
		Repository: v["url"].(string),
	}
	if directory := v["directory"].(string); directory != "" {
		gitRepoVolume.Directory = pointer.FromString(directory)
	}
	if revision := v["revision"].(string); revision != "" {
		gitRepoVolume.Revision = pointer.FromString(revision)
	}
	return gitRepoVolume
}

func expandSecrets(secretsMap map[string]interface{}) map[string]string {
	if len(secretsMap) == 0 {
		return nil
	}
	output := make(map[string]string, len(secretsMap))

	for name, value := range secretsMap {
		output[name] = value.(string)
	}

	return output
}

func expandContainerProbe(input interface{}) *containerinstance.ContainerProbe {
	probe := containerinstance.ContainerProbe{}
	probeRaw := input.([]interface{})

	if len(probeRaw) == 0 {
		return nil
	}

	for _, p := range probeRaw {
		if p == nil {
			continue
		}
		probeConfig := p.(map[string]interface{})

		if v := probeConfig["initial_delay_seconds"].(int); v > 0 {
			probe.InitialDelaySeconds = pointer.FromInt64(int64(v))
		}

		if v := probeConfig["period_seconds"].(int); v > 0 {
			probe.PeriodSeconds = pointer.FromInt64(int64(v))
		}

		if v := probeConfig["failure_threshold"].(int); v > 0 {
			probe.FailureThreshold = pointer.FromInt64(int64(v))
		}

		if v := probeConfig["success_threshold"].(int); v > 0 {
			probe.SuccessThreshold = pointer.FromInt64(int64(v))
		}

		if v := probeConfig["timeout_seconds"].(int); v > 0 {
			probe.TimeoutSeconds = pointer.FromInt64(int64(v))
		}

		commands := probeConfig["exec"].([]interface{})
		if len(commands) > 0 {
			exec := containerinstance.ContainerExec{
				Command: utils.ExpandStringSlice(commands),
			}
			probe.Exec = &exec
		}

		httpRaw := probeConfig["http_get"].([]interface{})
		if len(httpRaw) > 0 {
			for _, httpget := range httpRaw {
				if httpget == nil {
					continue
				}
				x := httpget.(map[string]interface{})

				path := x["path"].(string)
				port := x["port"].(int)
				scheme := x["scheme"].(string)

				httpGetScheme := containerinstance.Scheme(scheme)

				probe.HTTPGet = &containerinstance.ContainerHTTPGet{
					Path:        pointer.FromString(path),
					Port:        int64(port),
					Scheme:      &httpGetScheme,
					HTTPHeaders: expandContainerProbeHttpHeaders(x["http_headers"].(map[string]interface{})),
				}
			}
		}
	}
	return &probe
}

func expandContainerProbeHttpHeaders(input map[string]interface{}) *[]containerinstance.HTTPHeader {
	if len(input) == 0 {
		return nil
	}

	headers := []containerinstance.HTTPHeader{}
	for k, v := range input {
		header := containerinstance.HTTPHeader{
			Name:  pointer.FromString(k),
			Value: pointer.FromString(v.(string)),
		}
		headers = append(headers, header)
	}
	return &headers
}

func flattenContainerProbeHttpHeaders(input *[]containerinstance.HTTPHeader) map[string]interface{} {
	if input == nil {
		return nil
	}

	output := map[string]interface{}{}
	for _, header := range *input {
		name := ""
		if header.Name != nil {
			name = *header.Name
		}
		value := ""
		if header.Value != nil {
			value = *header.Value
		}
		output[name] = value
	}
	return output
}

func flattenContainerImageRegistryCredentials(d *pluginsdk.ResourceData, input *[]containerinstance.ImageRegistryCredential) []interface{} {
	if input == nil {
		return nil
	}
	configsOld := d.Get("image_registry_credential").([]interface{})

	output := make([]interface{}, 0)
	for i, cred := range *input {
		credConfig := make(map[string]interface{})
		credConfig["server"] = cred.Server
		credConfig["username"] = cred.Username
		credConfig["user_assigned_identity_id"] = cred.Identity

		if len(configsOld) > i {
			data := configsOld[i].(map[string]interface{})
			oldServer := data["server"].(string)
			if cred.Server == oldServer {
				if v, ok := d.GetOk(fmt.Sprintf("image_registry_credential.%d.password", i)); ok {
					credConfig["password"] = v.(string)
				}
			}
		}

		output = append(output, credConfig)
	}
	return output
}

func flattenContainerGroupInitContainers(d *pluginsdk.ResourceData, initContainers *[]containerinstance.InitContainerDefinition, containerGroupVolumes *[]containerinstance.Volume) []interface{} {
	if initContainers == nil {
		return nil
	}
	// map old container names to index so we can look up things up
	nameIndexMap := map[string]int{}
	for i, c := range d.Get("init_container").([]interface{}) {
		cfg := c.(map[string]interface{})
		nameIndexMap[cfg["name"].(string)] = i
	}

	containerCfg := make([]interface{}, 0, len(*initContainers))
	for _, container := range *initContainers {
		name := container.Name

		// get index from name
		index := nameIndexMap[name]

		containerConfig := make(map[string]interface{})
		containerConfig["name"] = name

		if v := container.Properties.Image; v != nil {
			containerConfig["image"] = *v
		}

		if container.Properties.EnvironmentVariables != nil {
			if len(*container.Properties.EnvironmentVariables) > 0 {
				containerConfig["environment_variables"] = flattenContainerEnvironmentVariables(container.Properties.EnvironmentVariables)
				containerConfig["secure_environment_variables"] = flattenContainerSecureEnvironmentVariables(container.Properties.EnvironmentVariables, d, index, "init_container")
			}
		}

		commands := make([]string, 0)
		if command := container.Properties.Command; command != nil {
			commands = *command
		}
		containerConfig["commands"] = commands

		if containerGroupVolumes != nil && container.Properties.VolumeMounts != nil {
			containersConfigRaw := d.Get("container").([]interface{})
			flattenContainerVolume(containerConfig, containersConfigRaw, container.Name, container.Properties.VolumeMounts, containerGroupVolumes)
		}

		containerConfig["security"] = flattenContainerSecurityContext(container.Properties.SecurityContext)

		containerCfg = append(containerCfg, containerConfig)
	}

	return containerCfg
}

func flattenContainerGroupContainers(d *pluginsdk.ResourceData, containers *[]containerinstance.Container, containerGroupVolumes *[]containerinstance.Volume) []interface{} {
	// map old container names to index so we can look up things up
	nameIndexMap := map[string]int{}
	for i, c := range d.Get("container").([]interface{}) {
		cfg := c.(map[string]interface{})
		nameIndexMap[cfg["name"].(string)] = i
	}

	containerCfg := make([]interface{}, 0, len(*containers))
	for _, container := range *containers {
		name := container.Name

		// get index from name
		index := nameIndexMap[name]

		containerConfig := make(map[string]interface{})
		containerConfig["name"] = name

		containerConfig["image"] = container.Properties.Image

		resources := container.Properties.Resources
		resourceRequests := resources.Requests
		containerConfig["cpu"] = resourceRequests.Cpu
		containerConfig["memory"] = resourceRequests.MemoryInGB

		if !features.FourPointOhBeta() {
			gpus := make([]interface{}, 0)
			if v := resourceRequests.Gpu; v != nil {
				gpu := make(map[string]interface{})
				gpu["count"] = v.Count
				gpu["sku"] = string(v.Sku)
				gpus = append(gpus, gpu)
			}
			containerConfig["gpu"] = gpus
		}

		if resourceLimits := resources.Limits; resourceLimits != nil {
			if v := resourceLimits.Cpu; v != nil {
				containerConfig["cpu_limit"] = *v
			}
			if v := resourceLimits.MemoryInGB; v != nil {
				containerConfig["memory_limit"] = *v
			}

			if !features.FourPointOhBeta() {
				gpus := make([]interface{}, 0)
				if v := resourceLimits.Gpu; v != nil {
					gpu := make(map[string]interface{})
					gpu["count"] = v.Count
					gpu["sku"] = string(v.Sku)
					gpus = append(gpus, gpu)
				}
				containerConfig["gpu_limit"] = gpus
			}
		}

		containerPorts := make([]interface{}, len(*container.Properties.Ports))
		if container.Properties.Ports != nil {
			for i := range *container.Properties.Ports {
				containerPorts[i] = (*container.Properties.Ports)[i]
			}
		}
		containerConfig["ports"] = flattenPorts(containerPorts)

		if container.Properties.EnvironmentVariables != nil {
			if len(*container.Properties.EnvironmentVariables) > 0 {
				containerConfig["environment_variables"] = flattenContainerEnvironmentVariables(container.Properties.EnvironmentVariables)
				containerConfig["secure_environment_variables"] = flattenContainerSecureEnvironmentVariables(container.Properties.EnvironmentVariables, d, index, "container")
			}
		}

		commands := make([]string, 0)
		if command := container.Properties.Command; command != nil {
			commands = *command
		}
		containerConfig["commands"] = commands

		if containerGroupVolumes != nil && container.Properties.VolumeMounts != nil {
			containersConfigRaw := d.Get("container").([]interface{})
			flattenContainerVolume(containerConfig, containersConfigRaw, container.Name, container.Properties.VolumeMounts, containerGroupVolumes)
		}

		containerConfig["liveness_probe"] = flattenContainerProbes(container.Properties.LivenessProbe)
		containerConfig["readiness_probe"] = flattenContainerProbes(container.Properties.ReadinessProbe)
		containerConfig["security"] = flattenContainerSecurityContext(container.Properties.SecurityContext)

		containerCfg = append(containerCfg, containerConfig)
	}

	return containerCfg
}

func flattenContainerVolume(containerConfig map[string]interface{}, containersConfigRaw []interface{}, containerName string, volumeMounts *[]containerinstance.VolumeMount, containerGroupVolumes *[]containerinstance.Volume) {
	// Also pass in the container volume config from schema
	var containerVolumesConfig *[]interface{}
	for _, containerConfigRaw := range containersConfigRaw {
		data := containerConfigRaw.(map[string]interface{})
		nameRaw := data["name"].(string)
		if nameRaw == containerName {
			// found container config for current container
			// extract volume mounts from config
			if v, ok := data["volume"]; ok {
				containerVolumesRaw := v.([]interface{})
				containerVolumesConfig = &containerVolumesRaw
			}
		}
	}
	volumeConfigs := make([]interface{}, 0)

	if volumeMounts == nil {
		containerConfig["volume"] = nil
		return
	}

	for _, vm := range *volumeMounts {
		volumeConfig := make(map[string]interface{})
		volumeConfig["name"] = vm.Name
		volumeConfig["mount_path"] = vm.MountPath
		if vm.ReadOnly != nil {
			volumeConfig["read_only"] = *vm.ReadOnly
		}

		// find corresponding volume in container group volumes
		// and use the data
		if containerGroupVolumes != nil {
			for _, cgv := range *containerGroupVolumes {
				if cgv.Name == vm.Name {
					if file := cgv.AzureFile; file != nil {
						volumeConfig["share_name"] = file.ShareName
						volumeConfig["storage_account_name"] = file.StorageAccountName
						// skip storage_account_key, is always nil
					}

					if cgv.EmptyDir != nil {
						volumeConfig["empty_dir"] = true
					}

					volumeConfig["git_repo"] = flattenGitRepoVolume(cgv.GitRepo)
				}
			}
		}

		// find corresponding volume in config
		// and use the data
		if containerVolumesConfig != nil {
			for _, cvr := range *containerVolumesConfig {
				cv := cvr.(map[string]interface{})
				rawName := cv["name"].(string)
				if vm.Name == rawName {
					storageAccountKey := cv["storage_account_key"].(string)
					volumeConfig["storage_account_key"] = storageAccountKey
					volumeConfig["secret"] = cv["secret"]
				}
			}
		}

		volumeConfigs = append(volumeConfigs, volumeConfig)
	}

	containerConfig["volume"] = volumeConfigs
}

func flattenContainerSecureEnvironmentVariables(input *[]containerinstance.EnvironmentVariable, d *pluginsdk.ResourceData, oldContainerIndex int, rootPropName string) map[string]interface{} {
	output := make(map[string]interface{})

	if input == nil {
		return output
	}

	for _, envVar := range *input {
		if envVar.Value == nil {
			envVarValue := d.Get(fmt.Sprintf("%s.%d.secure_environment_variables.%s", rootPropName, oldContainerIndex, envVar.Name))
			output[envVar.Name] = envVarValue
		}
	}

	return output
}
func flattenContainerEnvironmentVariables(input *[]containerinstance.EnvironmentVariable) map[string]interface{} {
	output := make(map[string]interface{})

	if input == nil {
		return output
	}

	for _, envVar := range *input {
		if envVar.Value != nil {
			log.Printf("[DEBUG] NOT SECURE: Name: %s - Value: %s", envVar.Name, *envVar.Value)
			output[envVar.Name] = *envVar.Value
		}
	}

	return output
}

func flattenGitRepoVolume(input *containerinstance.GitRepoVolume) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	var revision, directory, repository string
	if input.Directory != nil {
		directory = *input.Directory
	}
	if input.Revision != nil {
		revision = *input.Revision
	}
	repository = input.Repository
	return []interface{}{
		map[string]interface{}{
			"url":       repository,
			"directory": directory,
			"revision":  revision,
		},
	}
}

func flattenContainerProbes(input *containerinstance.ContainerProbe) []interface{} {
	outputs := make([]interface{}, 0)
	if input == nil {
		return outputs
	}

	output := make(map[string]interface{})

	if v := input.Exec; v != nil {
		output["exec"] = *v.Command
	}

	httpGets := make([]interface{}, 0)
	if get := input.HTTPGet; get != nil {
		httpGet := make(map[string]interface{})
		if v := get.Path; v != nil {
			httpGet["path"] = *v
		}
		httpGet["port"] = get.Port
		httpGet["scheme"] = get.Scheme
		httpGet["http_headers"] = flattenContainerProbeHttpHeaders(get.HTTPHeaders)
		httpGets = append(httpGets, httpGet)
	}
	output["http_get"] = httpGets

	if v := input.FailureThreshold; v != nil {
		output["failure_threshold"] = *v
	}

	if v := input.InitialDelaySeconds; v != nil {
		output["initial_delay_seconds"] = *v
	}

	if v := input.PeriodSeconds; v != nil {
		output["period_seconds"] = *v
	}

	if v := input.SuccessThreshold; v != nil {
		output["success_threshold"] = *v
	}

	if v := input.TimeoutSeconds; v != nil {
		output["timeout_seconds"] = *v
	}

	outputs = append(outputs, output)
	return outputs
}

func expandContainerGroupDiagnostics(input []interface{}) *containerinstance.ContainerGroupDiagnostics {
	if len(input) == 0 {
		return nil
	}

	vs := input[0].(map[string]interface{})

	analyticsVs := vs["log_analytics"].([]interface{})
	analyticsV := analyticsVs[0].(map[string]interface{})

	workspaceId := analyticsV["workspace_id"].(string)
	workspaceKey := analyticsV["workspace_key"].(string)

	logAnalytics := containerinstance.LogAnalytics{
		WorkspaceId:  workspaceId,
		WorkspaceKey: workspaceKey,
	}

	if logType := analyticsV["log_type"].(string); logType != "" {
		t := containerinstance.LogAnalyticsLogType(logType)
		logAnalytics.LogType = &t

		metadataMap := analyticsV["metadata"].(map[string]interface{})
		metadata := make(map[string]string)
		for k, v := range metadataMap {
			strValue := v.(string)
			metadata[k] = strValue
		}

		logAnalytics.Metadata = &metadata
	}

	return &containerinstance.ContainerGroupDiagnostics{LogAnalytics: &logAnalytics}
}

func flattenContainerGroupDiagnostics(d *pluginsdk.ResourceData, input *containerinstance.ContainerGroupDiagnostics) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	logAnalytics := make([]interface{}, 0)

	if la := input.LogAnalytics; la != nil {
		output := make(map[string]interface{})

		logType := ""
		if la.LogType != nil {
			logType = string(*la.LogType)
		}
		output["log_type"] = logType

		metadata := make(map[string]interface{})
		if la.Metadata != nil {
			for k, v := range *la.Metadata {
				metadata[k] = v
			}
		}
		output["metadata"] = metadata
		output["workspace_id"] = la.WorkspaceId

		// the existing config may not exist at Import time, protect against it.
		workspaceKey := ""
		if existingDiags := d.Get("diagnostics").([]interface{}); len(existingDiags) > 0 {
			existingDiag := existingDiags[0].(map[string]interface{})
			if existingLA := existingDiag["log_analytics"].([]interface{}); len(existingLA) > 0 {
				vs := existingLA[0].(map[string]interface{})
				if key := vs["workspace_key"]; key != nil && key.(string) != "" {
					workspaceKey = key.(string)
				}
			}
		}
		output["workspace_key"] = workspaceKey

		logAnalytics = append(logAnalytics, output)
	}

	return []interface{}{
		map[string]interface{}{
			"log_analytics": logAnalytics,
		},
	}
}

func resourceContainerGroupPortsHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%d-", m["port"].(int)))
		buf.WriteString(fmt.Sprintf("%s-", m["protocol"].(string)))
	}

	return pluginsdk.HashString(buf.String())
}

func flattenContainerGroupDnsConfig(input *containerinstance.DnsConfiguration) []interface{} {
	output := make(map[string]interface{})

	if input == nil {
		return make([]interface{}, 0)
	}

	// We're converting to TypeSet here from an API response that looks like "a b c" (assumes space delimited)
	var searchDomains []string
	if input.SearchDomains != nil {
		searchDomains = strings.Fields(*input.SearchDomains)
	}
	output["search_domains"] = searchDomains

	// We're converting to TypeSet here from an API response that looks like "a b c" (assumes space delimited)
	var options []string
	if input.Options != nil {
		options = strings.Fields(*input.Options)
	}
	output["options"] = options
	output["nameservers"] = input.NameServers

	return []interface{}{output}
}

func expandContainerGroupDnsConfig(input interface{}) *containerinstance.DnsConfiguration {
	dnsConfigRaw := input.([]interface{})
	if len(dnsConfigRaw) > 0 && dnsConfigRaw[0] != nil {
		config := dnsConfigRaw[0].(map[string]interface{})

		nameservers := []string{}
		for _, v := range config["nameservers"].([]interface{}) {
			nameservers = append(nameservers, v.(string))
		}
		options := []string{}
		for _, v := range config["options"].(*pluginsdk.Set).List() {
			options = append(options, v.(string))
		}
		searchDomains := []string{}
		for _, v := range config["search_domains"].(*pluginsdk.Set).List() {
			searchDomains = append(searchDomains, v.(string))
		}

		return &containerinstance.DnsConfiguration{
			Options:       pointer.FromString(strings.Join(options, " ")),
			SearchDomains: pointer.FromString(strings.Join(searchDomains, " ")),
			NameServers:   nameservers,
		}
	}

	return nil
}

func flattenContainerGroupSubnets(input *[]containerinstance.ContainerGroupSubnetId) ([]interface{}, error) {
	subnetIDs := make([]interface{}, 0)
	if input == nil {
		return subnetIDs, nil
	}

	for _, resourceRef := range *input {
		if resourceRef.Id == "" {
			continue
		}

		id, err := commonids.ParseSubnetIDInsensitively(resourceRef.Id)
		if err != nil {
			return nil, fmt.Errorf(`parsing subnet id %q: %v`, resourceRef.Id, err)
		}

		subnetIDs = append(subnetIDs, id.ID())
	}

	return subnetIDs, nil
}

func expandContainerGroupSubnets(input []interface{}) (*[]containerinstance.ContainerGroupSubnetId, error) {
	if len(input) == 0 {
		return nil, nil
	}

	results := make([]containerinstance.ContainerGroupSubnetId, 0)
	for _, item := range input {
		id, err := commonids.ParseSubnetID(item.(string))
		if err != nil {
			return nil, fmt.Errorf(`parsing subnet id %q: %v`, item, err)
		}

		results = append(results, containerinstance.ContainerGroupSubnetId{
			Id: id.ID(),
		})
	}
	return &results, nil
}
