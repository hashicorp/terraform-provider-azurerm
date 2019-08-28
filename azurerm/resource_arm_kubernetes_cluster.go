package azurerm

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-06-01/containerservice"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/kubernetes"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKubernetesClusterCreateUpdate,
		Read:   resourceArmKubernetesClusterRead,
		Update: resourceArmKubernetesClusterCreateUpdate,
		Delete: resourceArmKubernetesClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			if v, exists := diff.GetOk("network_profile"); exists {
				rawProfiles := v.([]interface{})
				if len(rawProfiles) == 0 {
					return nil
				}

				// then ensure the conditionally-required fields are set
				profile := rawProfiles[0].(map[string]interface{})
				networkPlugin := profile["network_plugin"].(string)

				if networkPlugin != "kubenet" && networkPlugin != "azure" {
					return nil
				}

				dockerBridgeCidr := profile["docker_bridge_cidr"].(string)
				dnsServiceIP := profile["dns_service_ip"].(string)
				serviceCidr := profile["service_cidr"].(string)

				// All empty values.
				if dockerBridgeCidr == "" && dnsServiceIP == "" && serviceCidr == "" {
					return nil
				}

				// All set values.
				if dockerBridgeCidr != "" && dnsServiceIP != "" && serviceCidr != "" {
					return nil
				}

				return fmt.Errorf("`docker_bridge_cidr`, `dns_service_ip` and `service_cidr` should all be empty or all should be set.")
			}

			return nil
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"dns_prefix": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.KubernetesDNSPrefix,
			},

			"kubernetes_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"agent_pool_profile": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.KubernetesAgentPoolName,
						},

						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  string(containerservice.AvailabilitySet),
							ValidateFunc: validation.StringInSlice([]string{
								string(containerservice.AvailabilitySet),
								string(containerservice.VirtualMachineScaleSets),
							}, false),
						},

						"count": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(1, 100),
						},

						"max_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 100),
						},

						"min_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 100),
						},

						"enable_auto_scaling": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"availability_zones": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						// TODO: remove this field in the next major version
						"dns_prefix": {
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "This field has been removed by Azure",
						},

						"fqdn": {
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "This field has been deprecated. Use the parent `fqdn` instead",
						},

						"vm_size": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validate.NoEmptyStrings,
						},

						"os_disk_size_gb": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							Computed:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},

						"vnet_subnet_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"os_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  string(containerservice.Linux),
							ValidateFunc: validation.StringInSlice([]string{
								string(containerservice.Linux),
								string(containerservice.Windows),
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"max_pods": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},

						"node_taints": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			// TODO: 2.0 - we should be able to make this a List to be able to detect changes in the Client Secret
			"service_principal": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"client_secret": {
							Type:         schema.TypeString,
							ForceNew:     true,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
				Set: resourceKubernetesClusterServicePrincipalProfileHash,
			},

			// Optional
			"addon_profile": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http_application_routing": {
							Type:     schema.TypeList,
							MaxItems: 1,
							ForceNew: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										ForceNew: true,
										Required: true,
									},
									"http_application_routing_zone_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},

						"oms_agent": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"log_analytics_workspace_id": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: azure.ValidateResourceID,
									},
								},
							},
						},

						"aci_connector_linux": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"subnet_name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.NoEmptyStrings,
									},
								},
							},
						},
					},
				},
			},

			"linux_profile": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admin_username": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.KubernetesAdminUserName,
						},
						"ssh_key": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 1,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key_data": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validate.NoEmptyStrings,
									},
								},
							},
						},
					},
				},
			},

			"windows_profile": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admin_username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"admin_password": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"network_profile": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_plugin": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(containerservice.Azure),
								string(containerservice.Kubenet),
							}, false),
						},

						"network_policy": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(containerservice.NetworkPolicyCalico),
								string(containerservice.NetworkPolicyAzure),
							}, false),
						},

						"dns_service_ip": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validate.IPv4Address,
						},

						"docker_bridge_cidr": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validate.CIDR,
						},

						"pod_cidr": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validate.CIDR,
						},

						"service_cidr": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validate.CIDR,
						},

						"load_balancer_sku": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(containerservice.Basic),
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(containerservice.Basic),
								string(containerservice.Standard),
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},
				},
			},

			"role_based_access_control": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},
						"azure_active_directory": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"client_app_id": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validate.UUID,
									},

									"server_app_id": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validate.UUID,
									},

									"server_app_secret": {
										Type:         schema.TypeString,
										ForceNew:     true,
										Required:     true,
										Sensitive:    true,
										ValidateFunc: validate.NoEmptyStrings,
									},

									"tenant_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
										// OrEmpty since this can be sourced from the client config if it's not specified
										ValidateFunc: validate.UUIDOrEmpty,
									},
								},
							},
						},
					},
				},
			},

			"tags": tags.Schema(),

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// Computed
			"kube_admin_config": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"client_certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_key": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"cluster_ca_certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"kube_admin_config_raw": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"kube_config": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"client_certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_key": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"cluster_ca_certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"kube_config_raw": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"node_resource_group": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"api_server_authorized_ip_ranges": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.CIDR,
				},
			},

			"enable_pod_security_policy": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceArmKubernetesClusterCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).containers.KubernetesClustersClient
	ctx := meta.(*ArmClient).StopContext
	tenantId := meta.(*ArmClient).tenantId

	log.Printf("[INFO] preparing arguments for Managed Kubernetes Cluster create/update.")

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Kubernetes Cluster %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_kubernetes_cluster", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	dnsPrefix := d.Get("dns_prefix").(string)
	kubernetesVersion := d.Get("kubernetes_version").(string)

	linuxProfile := expandKubernetesClusterLinuxProfile(d)
	agentProfiles, err := expandKubernetesClusterAgentPoolProfiles(d)
	if err != nil {
		return err
	}
	windowsProfile := expandKubernetesClusterWindowsProfile(d)
	servicePrincipalProfile := expandAzureRmKubernetesClusterServicePrincipal(d)
	networkProfile := expandKubernetesClusterNetworkProfile(d)
	addonProfiles := expandKubernetesClusterAddonProfiles(d)

	t := d.Get("tags").(map[string]interface{})

	rbacRaw := d.Get("role_based_access_control").([]interface{})
	rbacEnabled, azureADProfile := expandKubernetesClusterRoleBasedAccessControl(rbacRaw, tenantId)

	apiServerAuthorizedIPRangesRaw := d.Get("api_server_authorized_ip_ranges").(*schema.Set).List()
	apiServerAuthorizedIPRanges := utils.ExpandStringSlice(apiServerAuthorizedIPRangesRaw)

	nodeResourceGroup := d.Get("node_resource_group").(string)

	enablePodSecurityPolicy := d.Get("enable_pod_security_policy").(bool)

	parameters := containerservice.ManagedCluster{
		Name:     &name,
		Location: &location,
		ManagedClusterProperties: &containerservice.ManagedClusterProperties{
			APIServerAuthorizedIPRanges: apiServerAuthorizedIPRanges,
			AadProfile:                  azureADProfile,
			AddonProfiles:               addonProfiles,
			AgentPoolProfiles:           &agentProfiles,
			DNSPrefix:                   utils.String(dnsPrefix),
			EnableRBAC:                  utils.Bool(rbacEnabled),
			KubernetesVersion:           utils.String(kubernetesVersion),
			LinuxProfile:                linuxProfile,
			WindowsProfile:              windowsProfile,
			NetworkProfile:              networkProfile,
			ServicePrincipalProfile:     servicePrincipalProfile,
			NodeResourceGroup:           utils.String(nodeResourceGroup),
			EnablePodSecurityPolicy:     utils.Bool(enablePodSecurityPolicy),
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating/updating Managed Kubernetes Cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Managed Kubernetes Cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Managed Kubernetes Cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ID for Managed Kubernetes Cluster %q (Resource Group %q)", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmKubernetesClusterRead(d, meta)
}

func resourceArmKubernetesClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).containers.KubernetesClustersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["managedClusters"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Managed Kubernetes Cluster %q was not found in Resource Group %q - removing from state!", name, resGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Managed Kubernetes Cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	profile, err := client.GetAccessProfile(ctx, resGroup, name, "clusterUser")
	if err != nil {
		return fmt.Errorf("Error retrieving Access Profile for Managed Kubernetes Cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.ManagedClusterProperties; props != nil {
		d.Set("dns_prefix", props.DNSPrefix)
		d.Set("fqdn", props.Fqdn)
		d.Set("kubernetes_version", props.KubernetesVersion)
		d.Set("node_resource_group", props.NodeResourceGroup)
		d.Set("enable_pod_security_policy", props.EnablePodSecurityPolicy)

		apiServerAuthorizedIPRanges := utils.FlattenStringSlice(props.APIServerAuthorizedIPRanges)
		if err := d.Set("api_server_authorized_ip_ranges", apiServerAuthorizedIPRanges); err != nil {
			return fmt.Errorf("Error setting `api_server_authorized_ip_ranges`: %+v", err)
		}

		addonProfiles := flattenKubernetesClusterAddonProfiles(props.AddonProfiles)
		if err := d.Set("addon_profile", addonProfiles); err != nil {
			return fmt.Errorf("Error setting `addon_profile`: %+v", err)
		}

		agentPoolProfiles := flattenKubernetesClusterAgentPoolProfiles(props.AgentPoolProfiles, resp.Fqdn)
		if err := d.Set("agent_pool_profile", agentPoolProfiles); err != nil {
			return fmt.Errorf("Error setting `agent_pool_profile`: %+v", err)
		}

		linuxProfile := flattenKubernetesClusterLinuxProfile(props.LinuxProfile)
		if err := d.Set("linux_profile", linuxProfile); err != nil {
			return fmt.Errorf("Error setting `linux_profile`: %+v", err)
		}

		windowsProfile := flattenKubernetesClusterWindowsProfile(props.WindowsProfile, d)
		if err := d.Set("windows_profile", windowsProfile); err != nil {
			return fmt.Errorf("Error setting `windows_profile`: %+v", err)
		}

		networkProfile := flattenKubernetesClusterNetworkProfile(props.NetworkProfile)
		if err := d.Set("network_profile", networkProfile); err != nil {
			return fmt.Errorf("Error setting `network_profile`: %+v", err)
		}

		roleBasedAccessControl := flattenKubernetesClusterRoleBasedAccessControl(props, d)
		if err := d.Set("role_based_access_control", roleBasedAccessControl); err != nil {
			return fmt.Errorf("Error setting `role_based_access_control`: %+v", err)
		}

		servicePrincipal := flattenAzureRmKubernetesClusterServicePrincipalProfile(props.ServicePrincipalProfile)
		if err := d.Set("service_principal", servicePrincipal); err != nil {
			return fmt.Errorf("Error setting `service_principal`: %+v", err)
		}

		// adminProfile is only available for RBAC enabled clusters with AAD
		if props.AadProfile != nil {
			adminProfile, err := client.GetAccessProfile(ctx, resGroup, name, "clusterAdmin")
			if err != nil {
				return fmt.Errorf("Error retrieving Admin Access Profile for Managed Kubernetes Cluster %q (Resource Group %q): %+v", name, resGroup, err)
			}

			adminKubeConfigRaw, adminKubeConfig := flattenKubernetesClusterAccessProfile(adminProfile)
			d.Set("kube_admin_config_raw", adminKubeConfigRaw)
			if err := d.Set("kube_admin_config", adminKubeConfig); err != nil {
				return fmt.Errorf("Error setting `kube_admin_config`: %+v", err)
			}
		} else {
			d.Set("kube_admin_config_raw", "")
			d.Set("kube_admin_config", []interface{}{})
		}
	}

	kubeConfigRaw, kubeConfig := flattenKubernetesClusterAccessProfile(profile)
	d.Set("kube_config_raw", kubeConfigRaw)
	if err := d.Set("kube_config", kubeConfig); err != nil {
		return fmt.Errorf("Error setting `kube_config`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmKubernetesClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).containers.KubernetesClustersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["managedClusters"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Managed Kubernetes Cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deletion of Managed Kubernetes Cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func flattenKubernetesClusterAccessProfile(profile containerservice.ManagedClusterAccessProfile) (*string, []interface{}) {
	if accessProfile := profile.AccessProfile; accessProfile != nil {
		if kubeConfigRaw := accessProfile.KubeConfig; kubeConfigRaw != nil {
			rawConfig := string(*kubeConfigRaw)
			var flattenedKubeConfig []interface{}

			if strings.Contains(rawConfig, "apiserver-id:") {
				kubeConfigAAD, err := kubernetes.ParseKubeConfigAAD(rawConfig)
				if err != nil {
					return utils.String(rawConfig), []interface{}{}
				}

				flattenedKubeConfig = flattenKubernetesClusterKubeConfigAAD(*kubeConfigAAD)
			} else {
				kubeConfig, err := kubernetes.ParseKubeConfig(rawConfig)
				if err != nil {
					return utils.String(rawConfig), []interface{}{}
				}

				flattenedKubeConfig = flattenKubernetesClusterKubeConfig(*kubeConfig)
			}

			return utils.String(rawConfig), flattenedKubeConfig
		}
	}
	return nil, []interface{}{}
}

func expandKubernetesClusterAddonProfiles(d *schema.ResourceData) map[string]*containerservice.ManagedClusterAddonProfile {
	profiles := d.Get("addon_profile").([]interface{})
	if len(profiles) == 0 {
		return nil
	}

	profile := profiles[0].(map[string]interface{})
	addonProfiles := map[string]*containerservice.ManagedClusterAddonProfile{}

	httpApplicationRouting := profile["http_application_routing"].([]interface{})
	if len(httpApplicationRouting) > 0 {
		value := httpApplicationRouting[0].(map[string]interface{})
		enabled := value["enabled"].(bool)
		addonProfiles["httpApplicationRouting"] = &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(enabled),
		}
	}

	omsAgent := profile["oms_agent"].([]interface{})
	if len(omsAgent) > 0 {
		value := omsAgent[0].(map[string]interface{})
		config := make(map[string]*string)
		enabled := value["enabled"].(bool)

		if workspaceId, ok := value["log_analytics_workspace_id"]; ok {
			config["logAnalyticsWorkspaceResourceID"] = utils.String(workspaceId.(string))
		}

		addonProfiles["omsagent"] = &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(enabled),
			Config:  config,
		}
	}

	aciConnector := profile["aci_connector_linux"].([]interface{})
	if len(aciConnector) > 0 {
		value := aciConnector[0].(map[string]interface{})
		config := make(map[string]*string)
		enabled := value["enabled"].(bool)

		if subnetName, ok := value["subnet_name"]; ok {
			config["SubnetName"] = utils.String(subnetName.(string))
		}

		addonProfiles["aciConnectorLinux"] = &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(enabled),
			Config:  config,
		}
	}

	return addonProfiles
}

func flattenKubernetesClusterAddonProfiles(profile map[string]*containerservice.ManagedClusterAddonProfile) []interface{} {
	values := make(map[string]interface{})

	routes := make([]interface{}, 0)
	if httpApplicationRouting := profile["httpApplicationRouting"]; httpApplicationRouting != nil {
		enabled := false
		if enabledVal := httpApplicationRouting.Enabled; enabledVal != nil {
			enabled = *enabledVal
		}

		zoneName := ""
		if v := httpApplicationRouting.Config["HTTPApplicationRoutingZoneName"]; v != nil {
			zoneName = *v
		}

		output := map[string]interface{}{
			"enabled":                            enabled,
			"http_application_routing_zone_name": zoneName,
		}
		routes = append(routes, output)
	}
	values["http_application_routing"] = routes

	agents := make([]interface{}, 0)
	if omsAgent := profile["omsagent"]; omsAgent != nil {
		enabled := false
		if enabledVal := omsAgent.Enabled; enabledVal != nil {
			enabled = *enabledVal
		}

		workspaceId := ""
		if workspaceResourceID := omsAgent.Config["logAnalyticsWorkspaceResourceID"]; workspaceResourceID != nil {
			workspaceId = *workspaceResourceID
		}

		output := map[string]interface{}{
			"enabled":                    enabled,
			"log_analytics_workspace_id": workspaceId,
		}
		agents = append(agents, output)
	}
	values["oms_agent"] = agents

	aciConnectors := make([]interface{}, 0)
	if aciConnector := profile["aciConnectorLinux"]; aciConnector != nil {
		enabled := false
		if enabledVal := aciConnector.Enabled; enabledVal != nil {
			enabled = *enabledVal
		}

		subnetName := ""
		if v := aciConnector.Config["SubnetName"]; v != nil {
			subnetName = *v
		}

		output := map[string]interface{}{
			"enabled":     enabled,
			"subnet_name": subnetName,
		}
		aciConnectors = append(aciConnectors, output)
	}
	values["aci_connector_linux"] = aciConnectors

	return []interface{}{values}
}

func expandKubernetesClusterAgentPoolProfiles(d *schema.ResourceData) ([]containerservice.ManagedClusterAgentPoolProfile, error) {

	configs := d.Get("agent_pool_profile").([]interface{})

	profiles := make([]containerservice.ManagedClusterAgentPoolProfile, 0)
	for config_id := range configs {
		config := configs[config_id].(map[string]interface{})

		name := config["name"].(string)
		poolType := config["type"].(string)
		count := int32(config["count"].(int))
		vmSize := config["vm_size"].(string)
		osDiskSizeGB := int32(config["os_disk_size_gb"].(int))
		osType := config["os_type"].(string)

		profile := containerservice.ManagedClusterAgentPoolProfile{
			Name:         utils.String(name),
			Type:         containerservice.AgentPoolType(poolType),
			Count:        utils.Int32(count),
			VMSize:       containerservice.VMSizeTypes(vmSize),
			OsDiskSizeGB: utils.Int32(osDiskSizeGB),
			OsType:       containerservice.OSType(osType),
		}

		if maxPods := int32(config["max_pods"].(int)); maxPods > 0 {
			profile.MaxPods = utils.Int32(maxPods)
		}

		vnetSubnetID := config["vnet_subnet_id"].(string)
		if vnetSubnetID != "" {
			profile.VnetSubnetID = utils.String(vnetSubnetID)
		}

		if maxCount := int32(config["max_count"].(int)); maxCount > 0 {
			profile.MaxCount = utils.Int32(maxCount)
		}

		if minCount := int32(config["min_count"].(int)); minCount > 0 {
			profile.MinCount = utils.Int32(minCount)
		}

		if enableAutoScalingItf := config["enable_auto_scaling"]; enableAutoScalingItf != nil {
			profile.EnableAutoScaling = utils.Bool(enableAutoScalingItf.(bool))

			// Auto scaling will change the number of nodes, but the original count number should not be sent again.
			// This avoid the cluster being resized after creation.
			if *profile.EnableAutoScaling && !d.IsNewResource() {
				profile.Count = nil
			}
		}

		if availavilityZones := utils.ExpandStringSlice(config["availability_zones"].([]interface{})); len(*availavilityZones) > 0 {
			profile.AvailabilityZones = availavilityZones
		}

		if *profile.EnableAutoScaling && (profile.MinCount == nil || profile.MaxCount == nil) {
			return nil, fmt.Errorf("Can't create an AKS cluster with autoscaling enabled but not setting min_count or max_count")
		}

		if nodeTaints := utils.ExpandStringSlice(config["node_taints"].([]interface{})); len(*nodeTaints) > 0 {
			profile.NodeTaints = nodeTaints
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func flattenKubernetesClusterAgentPoolProfiles(profiles *[]containerservice.ManagedClusterAgentPoolProfile, fqdn *string) []interface{} {
	if profiles == nil {
		return []interface{}{}
	}

	agentPoolProfiles := make([]interface{}, 0)

	for _, profile := range *profiles {
		agentPoolProfile := make(map[string]interface{})

		if profile.Type != "" {
			agentPoolProfile["type"] = string(profile.Type)
		}

		if profile.Count != nil {
			agentPoolProfile["count"] = int(*profile.Count)
		}

		if profile.MinCount != nil {
			agentPoolProfile["min_count"] = int(*profile.MinCount)
		}

		if profile.MaxCount != nil {
			agentPoolProfile["max_count"] = int(*profile.MaxCount)
		}

		if profile.EnableAutoScaling != nil {
			agentPoolProfile["enable_auto_scaling"] = *profile.EnableAutoScaling
		}

		agentPoolProfile["availability_zones"] = utils.FlattenStringSlice(profile.AvailabilityZones)

		if fqdn != nil {
			// temporarily persist the parent FQDN here until `fqdn` is removed from the `agent_pool_profile`
			agentPoolProfile["fqdn"] = *fqdn
		}

		if profile.Name != nil {
			agentPoolProfile["name"] = *profile.Name
		}

		if profile.VMSize != "" {
			agentPoolProfile["vm_size"] = string(profile.VMSize)
		}

		if profile.OsDiskSizeGB != nil {
			agentPoolProfile["os_disk_size_gb"] = int(*profile.OsDiskSizeGB)
		}

		if profile.VnetSubnetID != nil {
			agentPoolProfile["vnet_subnet_id"] = *profile.VnetSubnetID
		}

		if profile.OsType != "" {
			agentPoolProfile["os_type"] = string(profile.OsType)
		}

		if profile.MaxPods != nil {
			agentPoolProfile["max_pods"] = int(*profile.MaxPods)
		}

		if profile.NodeTaints != nil {
			agentPoolProfile["node_taints"] = *profile.NodeTaints
		}

		agentPoolProfiles = append(agentPoolProfiles, agentPoolProfile)
	}

	return agentPoolProfiles
}

func expandKubernetesClusterLinuxProfile(d *schema.ResourceData) *containerservice.LinuxProfile {
	profiles := d.Get("linux_profile").([]interface{})

	if len(profiles) == 0 {
		return nil
	}

	config := profiles[0].(map[string]interface{})

	adminUsername := config["admin_username"].(string)
	linuxKeys := config["ssh_key"].([]interface{})

	keyData := ""
	if key, ok := linuxKeys[0].(map[string]interface{}); ok {
		keyData = key["key_data"].(string)
	}
	sshPublicKey := containerservice.SSHPublicKey{
		KeyData: &keyData,
	}

	sshPublicKeys := []containerservice.SSHPublicKey{sshPublicKey}

	profile := containerservice.LinuxProfile{
		AdminUsername: &adminUsername,
		SSH: &containerservice.SSHConfiguration{
			PublicKeys: &sshPublicKeys,
		},
	}

	return &profile
}

func expandKubernetesClusterWindowsProfile(d *schema.ResourceData) *containerservice.ManagedClusterWindowsProfile {
	profiles := d.Get("windows_profile").([]interface{})

	if len(profiles) == 0 {
		return nil
	}

	config := profiles[0].(map[string]interface{})

	adminUsername := config["admin_username"].(string)
	adminPassword := config["admin_password"].(string)

	profile := containerservice.ManagedClusterWindowsProfile{
		AdminUsername: &adminUsername,
		AdminPassword: &adminPassword,
	}

	return &profile
}

func flattenKubernetesClusterLinuxProfile(profile *containerservice.LinuxProfile) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	values := make(map[string]interface{})

	if username := profile.AdminUsername; username != nil {
		values["admin_username"] = *username
	}

	sshKeys := make([]interface{}, 0)
	if ssh := profile.SSH; ssh != nil {
		if keys := ssh.PublicKeys; keys != nil {
			for _, sshKey := range *keys {
				outputs := make(map[string]interface{})
				if keyData := sshKey.KeyData; keyData != nil {
					outputs["key_data"] = *keyData
				}
				sshKeys = append(sshKeys, outputs)
			}
		}
	}

	values["ssh_key"] = sshKeys

	return []interface{}{values}
}

func flattenKubernetesClusterWindowsProfile(profile *containerservice.ManagedClusterWindowsProfile, d *schema.ResourceData) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	values := make(map[string]interface{})

	if username := profile.AdminUsername; username != nil {
		values["admin_username"] = *username
	}

	// admin password isn't returned, so let's look it up
	if v, ok := d.GetOk("windows_profile.0.admin_password"); ok {
		values["admin_password"] = v.(string)
	}

	return []interface{}{values}
}

func expandKubernetesClusterNetworkProfile(d *schema.ResourceData) *containerservice.NetworkProfileType {
	configs := d.Get("network_profile").([]interface{})
	if len(configs) == 0 {
		return nil
	}

	config := configs[0].(map[string]interface{})

	networkPlugin := config["network_plugin"].(string)

	networkPolicy := config["network_policy"].(string)

	loadBalancerSku := config["load_balancer_sku"].(string)

	networkProfile := containerservice.NetworkProfileType{
		NetworkPlugin:   containerservice.NetworkPlugin(networkPlugin),
		NetworkPolicy:   containerservice.NetworkPolicy(networkPolicy),
		LoadBalancerSku: containerservice.LoadBalancerSku(loadBalancerSku),
	}

	if v, ok := config["dns_service_ip"]; ok && v.(string) != "" {
		dnsServiceIP := v.(string)
		networkProfile.DNSServiceIP = utils.String(dnsServiceIP)
	}

	if v, ok := config["pod_cidr"]; ok && v.(string) != "" {
		podCidr := v.(string)
		networkProfile.PodCidr = utils.String(podCidr)
	}

	if v, ok := config["docker_bridge_cidr"]; ok && v.(string) != "" {
		dockerBridgeCidr := v.(string)
		networkProfile.DockerBridgeCidr = utils.String(dockerBridgeCidr)
	}

	if v, ok := config["service_cidr"]; ok && v.(string) != "" {
		serviceCidr := v.(string)
		networkProfile.ServiceCidr = utils.String(serviceCidr)
	}

	return &networkProfile
}

func flattenKubernetesClusterNetworkProfile(profile *containerservice.NetworkProfileType) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	values := make(map[string]interface{})

	values["network_plugin"] = profile.NetworkPlugin

	if profile.NetworkPolicy != "" {
		values["network_policy"] = string(profile.NetworkPolicy)
	}

	if profile.ServiceCidr != nil {
		values["service_cidr"] = *profile.ServiceCidr
	}

	if profile.DNSServiceIP != nil {
		values["dns_service_ip"] = *profile.DNSServiceIP
	}

	if profile.DockerBridgeCidr != nil {
		values["docker_bridge_cidr"] = *profile.DockerBridgeCidr
	}

	if profile.PodCidr != nil {
		values["pod_cidr"] = *profile.PodCidr
	}

	if profile.LoadBalancerSku != "" {
		values["load_balancer_sku"] = string(profile.LoadBalancerSku)
	}

	return []interface{}{values}
}

func expandKubernetesClusterRoleBasedAccessControl(input []interface{}, providerTenantId string) (bool, *containerservice.ManagedClusterAADProfile) {
	if len(input) == 0 {
		return false, nil
	}

	val := input[0].(map[string]interface{})

	rbacEnabled := val["enabled"].(bool)
	azureADsRaw := val["azure_active_directory"].([]interface{})

	var aad *containerservice.ManagedClusterAADProfile
	if len(azureADsRaw) > 0 {
		azureAdRaw := azureADsRaw[0].(map[string]interface{})

		clientAppId := azureAdRaw["client_app_id"].(string)
		serverAppId := azureAdRaw["server_app_id"].(string)
		serverAppSecret := azureAdRaw["server_app_secret"].(string)
		tenantId := azureAdRaw["tenant_id"].(string)

		if tenantId == "" {
			tenantId = providerTenantId
		}

		aad = &containerservice.ManagedClusterAADProfile{
			ClientAppID:     utils.String(clientAppId),
			ServerAppID:     utils.String(serverAppId),
			ServerAppSecret: utils.String(serverAppSecret),
			TenantID:        utils.String(tenantId),
		}
	}

	return rbacEnabled, aad
}

func flattenKubernetesClusterRoleBasedAccessControl(input *containerservice.ManagedClusterProperties, d *schema.ResourceData) []interface{} {
	rbacEnabled := false
	if input.EnableRBAC != nil {
		rbacEnabled = *input.EnableRBAC
	}

	results := make([]interface{}, 0)
	if profile := input.AadProfile; profile != nil {
		output := make(map[string]interface{})

		if profile.ClientAppID != nil {
			output["client_app_id"] = *profile.ClientAppID
		}

		if profile.ServerAppID != nil {
			output["server_app_id"] = *profile.ServerAppID
		}

		// since input.ServerAppSecret isn't returned we're pulling this out of the existing state (which won't work for Imports)
		// role_based_access_control.0.azure_active_directory.0.server_app_secret
		if existing, ok := d.GetOk("role_based_access_control"); ok {
			rbacRawVals := existing.([]interface{})
			if len(rbacRawVals) > 0 {
				rbacRawVal := rbacRawVals[0].(map[string]interface{})
				if azureADVals, ok := rbacRawVal["azure_active_directory"].([]interface{}); ok && len(azureADVals) > 0 {
					azureADVal := azureADVals[0].(map[string]interface{})
					v := azureADVal["server_app_secret"]
					if v != nil {
						output["server_app_secret"] = v.(string)
					}
				}
			}
		}

		if profile.TenantID != nil {
			output["tenant_id"] = *profile.TenantID
		}

		results = append(results, output)
	}

	return []interface{}{
		map[string]interface{}{
			"enabled":                rbacEnabled,
			"azure_active_directory": results,
		},
	}
}

func expandAzureRmKubernetesClusterServicePrincipal(d *schema.ResourceData) *containerservice.ManagedClusterServicePrincipalProfile {
	value, exists := d.GetOk("service_principal")
	if !exists {
		return nil
	}

	configs := value.(*schema.Set).List()

	config := configs[0].(map[string]interface{})

	clientId := config["client_id"].(string)
	clientSecret := config["client_secret"].(string)

	principal := containerservice.ManagedClusterServicePrincipalProfile{
		ClientID: &clientId,
		Secret:   &clientSecret,
	}

	return &principal
}

func flattenAzureRmKubernetesClusterServicePrincipalProfile(profile *containerservice.ManagedClusterServicePrincipalProfile) *schema.Set {
	if profile == nil {
		return nil
	}

	servicePrincipalProfiles := &schema.Set{
		F: resourceKubernetesClusterServicePrincipalProfileHash,
	}

	values := make(map[string]interface{})

	if clientId := profile.ClientID; clientId != nil {
		values["client_id"] = *clientId
	}
	if secret := profile.Secret; secret != nil {
		values["client_secret"] = *secret
	}

	servicePrincipalProfiles.Add(values)

	return servicePrincipalProfiles
}

func resourceKubernetesClusterServicePrincipalProfileHash(v interface{}) int {
	// TODO: this method should be able to be removed in time
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["client_id"].(string)))
	}

	return hashcode.String(buf.String())
}

func flattenKubernetesClusterKubeConfig(config kubernetes.KubeConfig) []interface{} {
	values := make(map[string]interface{})

	// we don't size-check these since they're validated in the Parse method
	cluster := config.Clusters[0].Cluster
	user := config.Users[0].User
	name := config.Users[0].Name

	values["host"] = cluster.Server
	values["username"] = name
	values["password"] = user.Token
	values["client_certificate"] = user.ClientCertificteData
	values["client_key"] = user.ClientKeyData
	values["cluster_ca_certificate"] = cluster.ClusterAuthorityData

	return []interface{}{values}
}

func flattenKubernetesClusterKubeConfigAAD(config kubernetes.KubeConfigAAD) []interface{} {
	values := make(map[string]interface{})

	// we don't size-check these since they're validated in the Parse method
	cluster := config.Clusters[0].Cluster
	name := config.Users[0].Name

	values["host"] = cluster.Server
	values["username"] = name

	values["password"] = ""
	values["client_certificate"] = ""
	values["client_key"] = ""

	values["cluster_ca_certificate"] = cluster.ClusterAuthorityData

	return []interface{}{values}
}
