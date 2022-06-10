package containers

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2021-03-01/containerinstance"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-08-01/network"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	networkParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceContainerGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceContainerGroupCreate,
		Read:   resourceContainerGroupRead,
		Delete: resourceContainerGroupDelete,
		Update: resourceContainerGroupUpdate,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ContainerGroupID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

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

			"network_profile_id": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  networkValidate.NetworkProfileID,
				ConflictsWith: []string{"dns_name_label"},
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

						"username": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"password": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							Sensitive:    true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"tags": tags.Schema(),

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

			"exposed_port": {
				Type:       pluginsdk.TypeSet,
				Optional:   true, // change to 'Required' in 3.0 of the provider
				ForceNew:   true,
				Computed:   true,                           // remove in 3.0 of the provider
				ConfigMode: pluginsdk.SchemaConfigModeAttr, // remove in 3.0 of the provider
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

						//lintignore:XS003
						"gpu": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							ForceNew: true,
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

func resourceContainerGroupCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.GroupsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewContainerGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_container_group", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	OSType := d.Get("os_type").(string)
	IPAddressType := d.Get("ip_address_type").(string)
	t := d.Get("tags").(map[string]interface{})
	restartPolicy := d.Get("restart_policy").(string)
	diagnosticsRaw := d.Get("diagnostics").([]interface{})
	diagnostics := expandContainerGroupDiagnostics(diagnosticsRaw)
	dnsConfig := d.Get("dns_config").([]interface{})
	addedEmptyDirs := map[string]bool{}
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
		Name:     utils.String(id.Name),
		Location: &location,
		Tags:     tags.Expand(t),
		ContainerGroupProperties: &containerinstance.ContainerGroupProperties{
			InitContainers:           initContainers,
			Containers:               containers,
			Diagnostics:              diagnostics,
			RestartPolicy:            containerinstance.ContainerGroupRestartPolicy(restartPolicy),
			OsType:                   containerinstance.OperatingSystemTypes(OSType),
			Volumes:                  &containerGroupVolumes,
			ImageRegistryCredentials: expandContainerImageRegistryCredentials(d),
			DNSConfig:                expandContainerGroupDnsConfig(dnsConfig),
		},
	}

	// Container Groups with OS Type Windows do not support managed identities but the API also does not accept Identity Type: None
	// https://github.com/Azure/azure-rest-api-specs/issues/18122
	if OSType != string(containerinstance.OperatingSystemTypesWindows) {
		expandedIdentity, err := expandContainerGroupIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		containerGroup.Identity = expandedIdentity
	}

	if IPAddressType != "None" {
		containerGroup.ContainerGroupProperties.IPAddress = &containerinstance.IPAddress{
			Ports: containerGroupPorts,
			Type:  containerinstance.ContainerGroupIPAddressType(IPAddressType),
		}

		if dnsNameLabel := d.Get("dns_name_label").(string); dnsNameLabel != "" {
			containerGroup.ContainerGroupProperties.IPAddress.DNSNameLabel = &dnsNameLabel
		}
	}

	// https://docs.microsoft.com/en-us/azure/container-instances/container-instances-vnet#virtual-network-deployment-limitations
	// https://docs.microsoft.com/en-us/azure/container-instances/container-instances-vnet#preview-limitations
	if networkProfileID := d.Get("network_profile_id").(string); networkProfileID != "" {
		id, _ := networkParse.NetworkProfileID(networkProfileID)
		networkProfileIDNorm := id.ID()
		// Avoid parallel provisioning if "network_profile_id" is given.
		// See: https://github.com/hashicorp/terraform-provider-azurerm/issues/15025
		locks.ByID(networkProfileIDNorm)
		defer locks.UnlockByID(networkProfileIDNorm)

		if strings.ToLower(OSType) != "linux" {
			return fmt.Errorf("Currently only Linux containers can be deployed to virtual networks")
		}
		containerGroup.ContainerGroupProperties.NetworkProfile = &containerinstance.ContainerGroupNetworkProfile{
			ID: &networkProfileIDNorm,
		}
	}

	if keyVaultKeyId := d.Get("key_vault_key_id").(string); keyVaultKeyId != "" {
		keyId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(keyVaultKeyId)
		if err != nil {
			return fmt.Errorf("parsing Key Vault Key ID: %+v", err)
		}
		containerGroup.ContainerGroupProperties.EncryptionProperties = &containerinstance.EncryptionProperties{
			VaultBaseURL: utils.String(keyId.KeyVaultBaseUrl),
			KeyName:      utils.String(keyId.Name),
			KeyVersion:   utils.String(keyId.Version),
		}
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, containerGroup)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceContainerGroupRead(d, meta)
}

func resourceContainerGroupUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.GroupsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ContainerGroupID(d.Id())
	if err != nil {
		return err
	}

	t := d.Get("tags").(map[string]interface{})

	parameters := containerinstance.Resource{
		Tags: tags.Expand(t),
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.Name, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceContainerGroupRead(d, meta)
}

func resourceContainerGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.GroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ContainerGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	identity, err := flattenContainerGroupIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if props := resp.ContainerGroupProperties; props != nil {
		containerConfigs := flattenContainerGroupContainers(d, resp.Containers, props.Volumes)
		if err := d.Set("container", containerConfigs); err != nil {
			return fmt.Errorf("setting `container`: %+v", err)
		}
		initContainerConfigs := flattenContainerGroupInitContainers(d, resp.InitContainers, props.Volumes)
		if err := d.Set("init_container", initContainerConfigs); err != nil {
			return fmt.Errorf("setting `init_container`: %+v", err)
		}

		if err := d.Set("image_registry_credential", flattenContainerImageRegistryCredentials(d, props.ImageRegistryCredentials)); err != nil {
			return fmt.Errorf("setting `image_registry_credential`: %+v", err)
		}

		if address := props.IPAddress; address != nil {
			d.Set("ip_address_type", address.Type)
			d.Set("ip_address", address.IP)
			exposedPorts := make([]interface{}, len(*resp.ContainerGroupProperties.IPAddress.Ports))
			for i := range *resp.ContainerGroupProperties.IPAddress.Ports {
				exposedPorts[i] = (*resp.ContainerGroupProperties.IPAddress.Ports)[i]
			}
			d.Set("exposed_port", flattenPorts(exposedPorts))
			d.Set("dns_name_label", address.DNSNameLabel)
			d.Set("fqdn", address.Fqdn)
		}

		d.Set("restart_policy", string(props.RestartPolicy))
		d.Set("os_type", string(props.OsType))
		d.Set("dns_config", flattenContainerGroupDnsConfig(resp.DNSConfig))

		if err := d.Set("diagnostics", flattenContainerGroupDiagnostics(d, props.Diagnostics)); err != nil {
			return fmt.Errorf("setting `diagnostics`: %+v", err)
		}

		if kvProps := props.EncryptionProperties; kvProps != nil {
			var keyVaultUri, keyName, keyVersion string
			if kvProps.VaultBaseURL != nil && *kvProps.VaultBaseURL != "" {
				keyVaultUri = *kvProps.VaultBaseURL
			} else {
				return fmt.Errorf("empty value returned for Key Vault URI")
			}
			if kvProps.KeyName != nil && *kvProps.KeyName != "" {
				keyName = *kvProps.KeyName
			} else {
				return fmt.Errorf("empty value returned for Key Vault Key Name")
			}
			if kvProps.KeyVersion != nil {
				keyVersion = *kvProps.KeyVersion
			}
			keyId, err := keyVaultParse.NewNestedItemID(keyVaultUri, "keys", keyName, keyVersion)
			if err != nil {
				return err
			}
			d.Set("key_vault_key_id", keyId.ID())
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenPorts(ports []interface{}) *pluginsdk.Set {
	if len(ports) > 0 {
		flatPorts := make([]interface{}, 0)
		for _, p := range ports {
			port := make(map[string]interface{})
			switch t := p.(type) {
			case containerinstance.Port:
				if v := t.Port; v != nil {
					port["port"] = int(*v)
				}
				port["protocol"] = string(t.Protocol)
			case containerinstance.ContainerPort:
				if v := t.Port; v != nil {
					port["port"] = int(*v)
				}
				port["protocol"] = string(t.Protocol)
			}
			flatPorts = append(flatPorts, port)
		}
		return pluginsdk.NewSet(resourceContainerGroupPortsHash, flatPorts)
	}
	return pluginsdk.NewSet(resourceContainerGroupPortsHash, make([]interface{}, 0))
}

func resourceContainerGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.GroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ContainerGroupID(d.Id())
	if err != nil {
		return err
	}

	networkProfileId := ""
	existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			// already deleted
			return nil
		}

		return fmt.Errorf("retrieving Container Group %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if props := existing.ContainerGroupProperties; props != nil {
		if profile := props.NetworkProfile; profile != nil {
			if profile.ID != nil {
				id, err := networkParse.NetworkProfileID(*profile.ID)
				if err != nil {
					return err
				}
				networkProfileId = id.ID()
			}
			// Avoid parallel deletion if "network_profile_id" is given. (not sure whether this is necessary)
			// See: https://github.com/hashicorp/terraform-provider-azurerm/issues/15025
			locks.ByID(networkProfileId)
			defer locks.UnlockByID(networkProfileId)
		}
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Container Group %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Container Group %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if networkProfileId != "" {
		networkProfileClient := meta.(*clients.Client).Network.ProfileClient
		networkProfileId, err := networkParse.NetworkProfileID(networkProfileId)
		if err != nil {
			return err
		}

		// TODO: remove when https://github.com/Azure/azure-sdk-for-go/issues/5082 has been fixed
		log.Printf("[DEBUG] Waiting for Container Group %q (Resource Group %q) to be finish deleting", id.Name, id.ResourceGroup)
		stateConf := &pluginsdk.StateChangeConf{
			Pending:                   []string{"Attached"},
			Target:                    []string{"Detached"},
			Refresh:                   containerGroupEnsureDetachedFromNetworkProfileRefreshFunc(ctx, networkProfileClient, networkProfileId.ResourceGroup, networkProfileId.Name, id.ResourceGroup, id.Name),
			MinTimeout:                15 * time.Second,
			ContinuousTargetOccurence: 5,
			Timeout:                   d.Timeout(pluginsdk.TimeoutDelete),
			NotFoundChecks:            40,
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for Container Group %q (Resource Group %q) to finish deleting: %s", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func containerGroupEnsureDetachedFromNetworkProfileRefreshFunc(ctx context.Context,
	client *network.ProfilesClient,
	networkProfileResourceGroup, networkProfileName,
	containerResourceGroupName, containerName string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		profile, err := client.Get(ctx, networkProfileResourceGroup, networkProfileName, "")
		if err != nil {
			return nil, "Error", fmt.Errorf("retrieving Network Profile %q (Resource Group %q): %s", networkProfileName, networkProfileResourceGroup, err)
		}

		exists := false
		if props := profile.ProfilePropertiesFormat; props != nil {
			if nics := props.ContainerNetworkInterfaces; nics != nil {
				for _, nic := range *nics {
					nicProps := nic.ContainerNetworkInterfacePropertiesFormat
					if nicProps == nil || nicProps.Container == nil || nicProps.Container.ID == nil {
						continue
					}

					parsedId, err := parse.ContainerGroupID(*nicProps.Container.ID)
					if err != nil {
						return nil, "", err
					}

					if parsedId.ResourceGroup != containerResourceGroupName {
						continue
					}

					if parsedId.Name == "" || parsedId.Name != containerName {
						continue
					}

					exists = true
					break
				}
			}
		}

		if exists {
			return nil, "Attached", nil
		}

		return profile, "Detached", nil
	}
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
			Name: utils.String(name),
			InitContainerPropertiesDefinition: &containerinstance.InitContainerPropertiesDefinition{
				Image: utils.String(image),
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
		container.EnvironmentVariables = envVars

		if v, ok := data["commands"]; ok {
			c := v.([]interface{})
			command := make([]string, 0)
			for _, v := range c {
				command = append(command, v.(string))
			}

			container.Command = &command
		}

		if v, ok := data["volume"]; ok {
			volumeMounts, _, err := expandSingleContainerVolume(v)
			if err != nil {
				return nil, nil, err
			}
			container.VolumeMounts = volumeMounts

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

func expandContainerGroupContainers(d *pluginsdk.ResourceData, addedEmptyDirs map[string]bool) (*[]containerinstance.Container, *[]containerinstance.Port, []containerinstance.Volume, error) {
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
			Name: utils.String(name),
			ContainerProperties: &containerinstance.ContainerProperties{
				Image: utils.String(image),
				Resources: &containerinstance.ResourceRequirements{
					Requests: &containerinstance.ResourceRequests{
						MemoryInGB: utils.Float(memory),
						CPU:        utils.Float(cpu),
					},
				},
			},
		}

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
					Count: &gpuCount,
					Sku:   gpuSku,
				}
				container.Resources.Requests.Gpu = &gpus
			}
		}

		if v, ok := data["ports"].(*pluginsdk.Set); ok && len(v.List()) > 0 {
			var ports []containerinstance.ContainerPort
			for _, v := range v.List() {
				portObj := v.(map[string]interface{})

				port := int32(portObj["port"].(int))
				proto := portObj["protocol"].(string)

				ports = append(ports, containerinstance.ContainerPort{
					Port:     &port,
					Protocol: containerinstance.ContainerNetworkProtocol(proto),
				})
				containerInstancePorts = append(containerInstancePorts, containerinstance.Port{
					Port:     &port,
					Protocol: containerinstance.ContainerGroupNetworkProtocol(proto),
				})
			}
			container.Ports = &ports
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
		container.EnvironmentVariables = envVars

		if v, ok := data["commands"]; ok {
			c := v.([]interface{})
			command := make([]string, 0)
			for _, v := range c {
				command = append(command, v.(string))
			}

			container.Command = &command
		}

		if v, ok := data["volume"]; ok {
			volumeMounts, _, err := expandSingleContainerVolume(v)
			if err != nil {
				return nil, nil, nil, err
			}
			container.VolumeMounts = volumeMounts

			expandedContainerGroupVolumes, err := expandContainerVolume(v, addedEmptyDirs, containerGroupVolumes)
			if err != nil {
				return nil, nil, nil, err
			}
			containerGroupVolumes = expandedContainerGroupVolumes
		}

		if v, ok := data["liveness_probe"]; ok {
			container.ContainerProperties.LivenessProbe = expandContainerProbe(v)
		}

		if v, ok := data["readiness_probe"]; ok {
			container.ContainerProperties.ReadinessProbe = expandContainerProbe(v)
		}

		containers = append(containers, container)
	}

	// Determine ports to be exposed on the group level, based on exposed_ports
	// and on what ports have been exposed on individual containers.
	if v, ok := d.Get("exposed_port").(*pluginsdk.Set); ok && len(v.List()) > 0 {
		cgpMap := make(map[int32]map[containerinstance.ContainerGroupNetworkProtocol]bool)
		for _, p := range containerInstancePorts {
			if val, ok := cgpMap[*p.Port]; ok {
				val[p.Protocol] = true
				cgpMap[*p.Port] = val
			} else {
				protoMap := map[containerinstance.ContainerGroupNetworkProtocol]bool{p.Protocol: true}
				cgpMap[*p.Port] = protoMap
			}
		}

		for _, p := range v.List() {
			portConfig := p.(map[string]interface{})
			port := int32(portConfig["port"].(int))
			proto := portConfig["protocol"].(string)
			if !cgpMap[port][containerinstance.ContainerGroupNetworkProtocol(proto)] {
				return nil, nil, nil, fmt.Errorf("Port %d/%s is not exposed on any individual container in the container group.\n"+
					"An exposed_ports block contains %d/%s, but no individual container has a ports block with the same port "+
					"and protocol. Any ports exposed on the container group must also be exposed on an individual container.",
					port, proto, port, proto)
			}
			containerGroupPorts = append(containerGroupPorts, containerinstance.Port{
				Port:     &port,
				Protocol: containerinstance.ContainerGroupNetworkProtocol(proto),
			})
		}
	} else {
		containerGroupPorts = containerInstancePorts // remove in 3.0 of the provider
	}

	return &containers, &containerGroupPorts, containerGroupVolumes, nil
}

func expandContainerVolume(v interface{}, addedEmptyDirs map[string]bool, containerGroupVolumes []containerinstance.Volume) ([]containerinstance.Volume, error) {
	_, containerVolumes, err := expandSingleContainerVolume(v)
	if err != nil {
		return nil, err
	}
	if containerVolumes != nil {
		for _, cgVol := range *containerVolumes {
			if cgVol.EmptyDir != nil {
				if addedEmptyDirs[*cgVol.Name] {
					// empty_dir-volumes are allowed to overlap across containers, in fact that is their primary purpose,
					// but the containerGroup must not declare same name of such volumes twice.
					continue
				}
				addedEmptyDirs[*cgVol.Name] = true
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
				Name:        utils.String(k),
				SecureValue: utils.String(v.(string)),
			}

			output = append(output, ev)
		}
	} else {
		for k, v := range envVars {
			ev := containerinstance.EnvironmentVariable{
				Name:  utils.String(k),
				Value: utils.String(v.(string)),
			}

			output = append(output, ev)
		}
	}
	return &output
}

func expandContainerGroupIdentity(input []interface{}) (*containerinstance.ContainerGroupIdentity, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := containerinstance.ContainerGroupIdentity{
		Type: containerinstance.ResourceIdentityType(string(expanded.Type)),
	}
	if expanded.Type == identity.TypeUserAssigned || expanded.Type == identity.TypeSystemAssignedUserAssigned {
		out.UserAssignedIdentities = make(map[string]*containerinstance.ContainerGroupIdentityUserAssignedIdentitiesValue)
		for k := range expanded.IdentityIds {
			out.UserAssignedIdentities[k] = &containerinstance.ContainerGroupIdentityUserAssignedIdentitiesValue{
				// intentionally empty
			}
		}
	}

	return &out, nil
}

func expandContainerImageRegistryCredentials(d *pluginsdk.ResourceData) *[]containerinstance.ImageRegistryCredential {
	credsRaw := d.Get("image_registry_credential").([]interface{})
	if len(credsRaw) == 0 {
		return nil
	}

	output := make([]containerinstance.ImageRegistryCredential, 0, len(credsRaw))

	for _, c := range credsRaw {
		credConfig := c.(map[string]interface{})

		output = append(output, containerinstance.ImageRegistryCredential{
			Server:   utils.String(credConfig["server"].(string)),
			Password: utils.String(credConfig["password"].(string)),
			Username: utils.String(credConfig["username"].(string)),
		})
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
			Name:      utils.String(name),
			MountPath: utils.String(mountPath),
			ReadOnly:  utils.Bool(readOnly),
		}

		volumeMounts = append(volumeMounts, vm)

		cv := containerinstance.Volume{
			Name: utils.String(name),
		}

		secret := expandSecrets(volumeConfig["secret"].(map[string]interface{}))

		gitRepoVolume := expandGitRepoVolume(volumeConfig["git_repo"].([]interface{}))

		switch {
		case emptyDir:
			if shareName != "" || storageAccountName != "" || storageAccountKey != "" || secret != nil || gitRepoVolume != nil {
				return nil, nil, fmt.Errorf("only one of `empty_dir` volume, `git_repo` volume, `secret` volume or storage account volume (`share_name`, `storage_account_name`, and `storage_account_key`) can be specified")
			}
			cv.EmptyDir = map[string]string{}
		case gitRepoVolume != nil:
			if shareName != "" || storageAccountName != "" || storageAccountKey != "" || secret != nil {
				return nil, nil, fmt.Errorf("only one of `empty_dir` volume, `git_repo` volume, `secret` volume or storage account volume (`share_name`, `storage_account_name`, and `storage_account_key`) can be specified")
			}
			cv.GitRepo = gitRepoVolume
		case secret != nil:
			if shareName != "" || storageAccountName != "" || storageAccountKey != "" {
				return nil, nil, fmt.Errorf("only one of `empty_dir` volume, `git_repo` volume, `secret` volume or storage account volume (`share_name`, `storage_account_name`, and `storage_account_key`) can be specified")
			}
			cv.Secret = secret
		default:
			if shareName == "" && storageAccountName == "" && storageAccountKey == "" {
				return nil, nil, fmt.Errorf("only one of `empty_dir` volume, `git_repo` volume, `secret` volume or storage account volume (`share_name`, `storage_account_name`, and `storage_account_key`) can be specified")
			} else if shareName == "" || storageAccountName == "" || storageAccountKey == "" {
				return nil, nil, fmt.Errorf("when using a storage account volume, all of `share_name`, `storage_account_name`, `storage_account_key` must be specified")
			}
			cv.AzureFile = &containerinstance.AzureFileVolume{
				ShareName:          utils.String(shareName),
				ReadOnly:           utils.Bool(readOnly),
				StorageAccountName: utils.String(storageAccountName),
				StorageAccountKey:  utils.String(storageAccountKey),
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
		Repository: utils.String(v["url"].(string)),
	}
	if directory := v["directory"].(string); directory != "" {
		gitRepoVolume.Directory = utils.String(directory)
	}
	if revision := v["revision"].(string); revision != "" {
		gitRepoVolume.Revision = utils.String(revision)
	}
	return gitRepoVolume
}

func expandSecrets(secretsMap map[string]interface{}) map[string]*string {
	if len(secretsMap) == 0 {
		return nil
	}
	output := make(map[string]*string, len(secretsMap))

	for name, value := range secretsMap {
		output[name] = utils.String(value.(string))
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
			probe.InitialDelaySeconds = utils.Int32(int32(v))
		}

		if v := probeConfig["period_seconds"].(int); v > 0 {
			probe.PeriodSeconds = utils.Int32(int32(v))
		}

		if v := probeConfig["failure_threshold"].(int); v > 0 {
			probe.FailureThreshold = utils.Int32(int32(v))
		}

		if v := probeConfig["success_threshold"].(int); v > 0 {
			probe.SuccessThreshold = utils.Int32(int32(v))
		}

		if v := probeConfig["timeout_seconds"].(int); v > 0 {
			probe.TimeoutSeconds = utils.Int32(int32(v))
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

				probe.HTTPGet = &containerinstance.ContainerHTTPGet{
					Path:   utils.String(path),
					Port:   utils.Int32(int32(port)),
					Scheme: containerinstance.Scheme(scheme),
				}
			}
		}
	}
	return &probe
}

func flattenContainerGroupIdentity(input *containerinstance.ContainerGroupIdentity) (*[]interface{}, error) {
	var transform *identity.SystemAndUserAssignedMap

	if input != nil {
		transform = &identity.SystemAndUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		if input.PrincipalID != nil {
			transform.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transform.TenantId = *input.TenantID
		}
		for k, v := range input.UserAssignedIdentities {
			transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				ClientId:    v.ClientID,
				PrincipalId: v.PrincipalID,
			}
		}
	}

	return identity.FlattenSystemAndUserAssignedMap(transform)
}

func flattenContainerImageRegistryCredentials(d *pluginsdk.ResourceData, input *[]containerinstance.ImageRegistryCredential) []interface{} {
	if input == nil {
		return nil
	}
	configsOld := d.Get("image_registry_credential").([]interface{})

	output := make([]interface{}, 0)
	for i, cred := range *input {
		credConfig := make(map[string]interface{})
		if cred.Server != nil {
			credConfig["server"] = *cred.Server
		}
		if cred.Username != nil {
			credConfig["username"] = *cred.Username
		}

		if len(configsOld) > i {
			data := configsOld[i].(map[string]interface{})
			oldServer := data["server"].(string)
			if cred.Server != nil && *cred.Server == oldServer {
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
		// TODO fix this crash point
		name := *container.Name

		// get index from name
		index := nameIndexMap[name]

		containerConfig := make(map[string]interface{})
		containerConfig["name"] = name

		if v := container.Image; v != nil {
			containerConfig["image"] = *v
		}

		if container.EnvironmentVariables != nil {
			if len(*container.EnvironmentVariables) > 0 {
				containerConfig["environment_variables"] = flattenContainerEnvironmentVariables(container.EnvironmentVariables, false, d, index)
				containerConfig["secure_environment_variables"] = flattenContainerEnvironmentVariables(container.EnvironmentVariables, true, d, index)
			}
		}

		commands := make([]string, 0)
		if command := container.Command; command != nil {
			commands = *command
		}
		containerConfig["commands"] = commands

		if containerGroupVolumes != nil && container.VolumeMounts != nil {
			containersConfigRaw := d.Get("container").([]interface{})
			flattenContainerVolume(containerConfig, containersConfigRaw, *container.Name, container.VolumeMounts, containerGroupVolumes)
		}
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
		// TODO fix this crash point
		name := *container.Name

		// get index from name
		index := nameIndexMap[name]

		containerConfig := make(map[string]interface{})
		containerConfig["name"] = name

		if v := container.Image; v != nil {
			containerConfig["image"] = *v
		}

		if resources := container.Resources; resources != nil {
			if resourceRequests := resources.Requests; resourceRequests != nil {
				if v := resourceRequests.CPU; v != nil {
					containerConfig["cpu"] = *v
				}
				if v := resourceRequests.MemoryInGB; v != nil {
					containerConfig["memory"] = *v
				}

				gpus := make([]interface{}, 0)
				if v := resourceRequests.Gpu; v != nil {
					gpu := make(map[string]interface{})
					if v.Count != nil {
						gpu["count"] = *v.Count
					}
					gpu["sku"] = string(v.Sku)
					gpus = append(gpus, gpu)
				}
				containerConfig["gpu"] = gpus
			}
		}

		containerPorts := make([]interface{}, len(*container.Ports))
		for i := range *container.Ports {
			containerPorts[i] = (*container.Ports)[i]
		}
		containerConfig["ports"] = flattenPorts(containerPorts)

		if container.EnvironmentVariables != nil {
			if len(*container.EnvironmentVariables) > 0 {
				containerConfig["environment_variables"] = flattenContainerEnvironmentVariables(container.EnvironmentVariables, false, d, index)
				containerConfig["secure_environment_variables"] = flattenContainerEnvironmentVariables(container.EnvironmentVariables, true, d, index)
			}
		}

		commands := make([]string, 0)
		if command := container.Command; command != nil {
			commands = *command
		}
		containerConfig["commands"] = commands

		if containerGroupVolumes != nil && container.VolumeMounts != nil {
			containersConfigRaw := d.Get("container").([]interface{})
			flattenContainerVolume(containerConfig, containersConfigRaw, *container.Name, container.VolumeMounts, containerGroupVolumes)
		}

		containerConfig["liveness_probe"] = flattenContainerProbes(container.LivenessProbe)
		containerConfig["readiness_probe"] = flattenContainerProbes(container.ReadinessProbe)

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
		if vm.Name != nil {
			volumeConfig["name"] = *vm.Name
		}
		if vm.MountPath != nil {
			volumeConfig["mount_path"] = *vm.MountPath
		}
		if vm.ReadOnly != nil {
			volumeConfig["read_only"] = *vm.ReadOnly
		}

		// find corresponding volume in container group volumes
		// and use the data
		if containerGroupVolumes != nil {
			for _, cgv := range *containerGroupVolumes {
				if cgv.Name == nil || vm.Name == nil {
					continue
				}

				if *cgv.Name == *vm.Name {
					if file := cgv.AzureFile; file != nil {
						if file.ShareName != nil {
							volumeConfig["share_name"] = *file.ShareName
						}
						if file.StorageAccountName != nil {
							volumeConfig["storage_account_name"] = *file.StorageAccountName
						}
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
				if vm.Name != nil && *vm.Name == rawName {
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

func flattenContainerEnvironmentVariables(input *[]containerinstance.EnvironmentVariable, isSecure bool, d *pluginsdk.ResourceData, oldContainerIndex int) map[string]interface{} {
	output := make(map[string]interface{})

	if input == nil {
		return output
	}

	if isSecure {
		for _, envVar := range *input {
			if envVar.Name != nil && envVar.Value == nil {
				envVarValue := d.Get(fmt.Sprintf("container.%d.secure_environment_variables.%s", oldContainerIndex, *envVar.Name))
				output[*envVar.Name] = envVarValue
			}
		}
	} else {
		for _, envVar := range *input {
			if envVar.Name != nil && envVar.Value != nil {
				log.Printf("[DEBUG] NOT SECURE: Name: %s - Value: %s", *envVar.Name, *envVar.Value)
				output[*envVar.Name] = *envVar.Value
			}
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
	if input.Repository != nil {
		repository = *input.Repository
	}
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

		if v := get.Port; v != nil {
			httpGet["port"] = *v
		}

		if get.Scheme != "" {
			httpGet["scheme"] = get.Scheme
		}

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
		WorkspaceID:  utils.String(workspaceId),
		WorkspaceKey: utils.String(workspaceKey),
	}

	if logType := analyticsV["log_type"].(string); logType != "" {
		logAnalytics.LogType = containerinstance.LogAnalyticsLogType(logType)

		metadataMap := analyticsV["metadata"].(map[string]interface{})
		metadata := make(map[string]*string)
		for k, v := range metadataMap {
			strValue := v.(string)
			metadata[k] = &strValue
		}

		logAnalytics.Metadata = metadata
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

		output["log_type"] = string(la.LogType)

		metadata := make(map[string]interface{})
		for k, v := range la.Metadata {
			metadata[k] = *v
		}
		output["metadata"] = metadata

		if la.WorkspaceID != nil {
			output["workspace_id"] = *la.WorkspaceID
		}

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

func flattenContainerGroupDnsConfig(input *containerinstance.DNSConfiguration) []interface{} {
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

	// Nameservers is already an array from the API
	var nameservers []string
	if input.NameServers != nil {
		nameservers = *input.NameServers
	}
	output["nameservers"] = nameservers

	return []interface{}{output}
}

func expandContainerGroupDnsConfig(input interface{}) *containerinstance.DNSConfiguration {
	dnsConfigRaw := input.([]interface{})
	if len(dnsConfigRaw) > 0 {
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

		return &containerinstance.DNSConfiguration{
			Options:       utils.String(strings.Join(options, " ")),
			SearchDomains: utils.String(strings.Join(searchDomains, " ")),
			NameServers:   &nameservers,
		}
	}

	return nil
}
