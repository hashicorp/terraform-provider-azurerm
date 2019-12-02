package azurerm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-06-01/containerservice"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/kubernetes"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKubernetesClusterCreate,
		Read:   resourceArmKubernetesClusterRead,
		Update: resourceArmKubernetesClusterUpdate,
		Delete: resourceArmKubernetesClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
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
				podCidr := profile["pod_cidr"].(string)

				// Azure network plugin is not compatible with pod_cidr
				if podCidr != "" && networkPlugin == "azure" {
					return fmt.Errorf("`pod_cidr` and `azure` cannot be set together.")
				}

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

			"default_node_pool": containers.SchemaDefaultNodePool(),

			// TODO: remove in 2.0
			"agent_pool_profile": {
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				Deprecated: "This has been replaced by `default_node_pool` and will be removed in version 2.0 of the AzureRM Provider",
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

						"enable_node_public_ip": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"service_principal": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"client_secret": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			// Optional
			"addon_profile": containers.SchemaKubernetesAddOnProfiles(),

			"api_server_authorized_ip_ranges": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.CIDR,
				},
			},

			// TODO: remove Computed in 2.0
			"enable_pod_security_policy": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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

			"node_resource_group": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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

			"windows_profile": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admin_username": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
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

			// Computed
			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

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
		},
	}
}

func resourceArmKubernetesClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Containers.KubernetesClustersClient
	ctx, cancel := timeouts.ForCreate(meta.(*ArmClient).StopContext, d)
	defer cancel()
	tenantId := meta.(*ArmClient).Account.TenantId

	log.Printf("[INFO] preparing arguments for Managed Kubernetes Cluster create.")

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	if features.ShouldResourcesBeImported() {
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

	linuxProfileRaw := d.Get("linux_profile").([]interface{})
	linuxProfile := expandKubernetesClusterLinuxProfile(linuxProfileRaw)

	agentProfiles, err := containers.ExpandDefaultNodePool(d)
	if err != nil {
		return fmt.Errorf("Error expanding `default_node_pool`: %+v", err)
	}

	// TODO: remove me in 2.0
	if agentProfiles == nil {
		agentProfilesRaw := d.Get("agent_pool_profile").([]interface{})
		agentProfilesLegacy, err := expandKubernetesClusterAgentPoolProfiles(agentProfilesRaw, true)
		if err != nil {
			return err
		}

		agentProfiles = &agentProfilesLegacy
	}

	addOnProfilesRaw := d.Get("addon_profile").([]interface{})
	addonProfiles := containers.ExpandKubernetesAddOnProfiles(addOnProfilesRaw)

	networkProfileRaw := d.Get("network_profile").([]interface{})
	networkProfile := expandKubernetesClusterNetworkProfile(networkProfileRaw)

	rbacRaw := d.Get("role_based_access_control").([]interface{})
	rbacEnabled, azureADProfile := expandKubernetesClusterRoleBasedAccessControl(rbacRaw, tenantId)

	// since the Create and Update use separate methods, there's no point extracting this out
	servicePrincipalProfileRaw := d.Get("service_principal").([]interface{})
	servicePrincipalProfileVal := servicePrincipalProfileRaw[0].(map[string]interface{})
	servicePrincipalProfile := &containerservice.ManagedClusterServicePrincipalProfile{
		ClientID: utils.String(servicePrincipalProfileVal["client_id"].(string)),
		Secret:   utils.String(servicePrincipalProfileVal["client_secret"].(string)),
	}

	t := d.Get("tags").(map[string]interface{})

	windowsProfileRaw := d.Get("windows_profile").([]interface{})
	windowsProfile := expandKubernetesClusterWindowsProfile(windowsProfileRaw)

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
			AgentPoolProfiles:           agentProfiles,
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
		return fmt.Errorf("Error creating Managed Kubernetes Cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Managed Kubernetes Cluster %q (Resource Group %q): %+v", name, resGroup, err)
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

func resourceArmKubernetesClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	nodePoolsClient := meta.(*ArmClient).Containers.AgentPoolsClient
	clusterClient := meta.(*ArmClient).Containers.KubernetesClustersClient
	ctx, cancel := timeouts.ForUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()
	tenantId := meta.(*ArmClient).Account.TenantId

	log.Printf("[INFO] preparing arguments for Managed Kubernetes Cluster update.")

	id, err := containers.ParseKubernetesClusterID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Name

	d.Partial(true)

	if d.HasChange("service_principal") {
		log.Printf("[DEBUG] Updating the Service Principal for Kubernetes Cluster %q (Resource Group %q)..", name, resourceGroup)
		servicePrincipals := d.Get("service_principal").([]interface{})
		servicePrincipalRaw := servicePrincipals[0].(map[string]interface{})

		clientId := servicePrincipalRaw["client_id"].(string)
		clientSecret := servicePrincipalRaw["client_secret"].(string)

		params := containerservice.ManagedClusterServicePrincipalProfile{
			ClientID: utils.String(clientId),
			Secret:   utils.String(clientSecret),
		}
		future, err := clusterClient.ResetServicePrincipalProfile(ctx, resourceGroup, name, params)
		if err != nil {
			return fmt.Errorf("Error updating Service Principal for Kubernetes Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, clusterClient.Client); err != nil {
			return fmt.Errorf("Error waiting for update of Service Principal for Kubernetes Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
		log.Printf("[DEBUG] Updated the Service Principal for Kubernetes Cluster %q (Resource Group %q).", name, resourceGroup)
	}

	// we need to conditionally update the cluster
	existing, err := clusterClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving existing Kubernetes Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if existing.ManagedClusterProperties == nil {
		return fmt.Errorf("Error retrieving existing Kubernetes Cluster %q (Resource Group %q): `properties` was nil", name, resourceGroup)
	}

	// since there's multiple reasons why we could be called into Update, we use this to only update if something's changed that's not SP/Version
	updateCluster := false

	if d.HasChange("addon_profile") {
		updateCluster = true
		addOnProfilesRaw := d.Get("addon_profile").([]interface{})
		addonProfiles := containers.ExpandKubernetesAddOnProfiles(addOnProfilesRaw)
		existing.ManagedClusterProperties.AddonProfiles = addonProfiles
	}

	if d.HasChange("api_server_authorized_ip_ranges") {
		updateCluster = true
		apiServerAuthorizedIPRangesRaw := d.Get("api_server_authorized_ip_ranges").(*schema.Set).List()
		existing.APIServerAuthorizedIPRanges = utils.ExpandStringSlice(apiServerAuthorizedIPRangesRaw)
	}

	if d.HasChange("enable_pod_security_policy") {
		updateCluster = true
		enablePodSecurityPolicy := d.Get("enable_pod_security_policy").(bool)
		existing.ManagedClusterProperties.EnablePodSecurityPolicy = utils.Bool(enablePodSecurityPolicy)
	}

	if d.HasChange("linux_profile") {
		updateCluster = true
		linuxProfileRaw := d.Get("linux_profile").([]interface{})
		linuxProfile := expandKubernetesClusterLinuxProfile(linuxProfileRaw)
		existing.ManagedClusterProperties.LinuxProfile = linuxProfile
	}

	if d.HasChange("network_profile") {
		updateCluster = true
		networkProfileRaw := d.Get("network_profile").([]interface{})
		networkProfile := expandKubernetesClusterNetworkProfile(networkProfileRaw)
		existing.ManagedClusterProperties.NetworkProfile = networkProfile
	}

	if d.HasChange("role_based_access_control") {
		updateCluster = true
		rbacRaw := d.Get("role_based_access_control").([]interface{})
		rbacEnabled, azureADProfile := expandKubernetesClusterRoleBasedAccessControl(rbacRaw, tenantId)
		existing.ManagedClusterProperties.AadProfile = azureADProfile
		existing.ManagedClusterProperties.EnableRBAC = utils.Bool(rbacEnabled)
	}

	if d.HasChange("tags") {
		updateCluster = true
		t := d.Get("tags").(map[string]interface{})
		existing.Tags = tags.Expand(t)
	}

	if d.HasChange("windows_profile") {
		updateCluster = true
		windowsProfileRaw := d.Get("windows_profile").([]interface{})
		windowsProfile := expandKubernetesClusterWindowsProfile(windowsProfileRaw)
		existing.ManagedClusterProperties.WindowsProfile = windowsProfile
	}

	if updateCluster {
		log.Printf("[DEBUG] Updating the Kubernetes Cluster %q (Resource Group %q)..", name, resourceGroup)
		future, err := clusterClient.CreateOrUpdate(ctx, resourceGroup, name, existing)
		if err != nil {
			return fmt.Errorf("Error updating Managed Kubernetes Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, clusterClient.Client); err != nil {
			return fmt.Errorf("Error waiting for update of Managed Kubernetes Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
		log.Printf("[DEBUG] Updated the Kubernetes Cluster %q (Resource Group %q)..", name, resourceGroup)
	}

	// update the node pool using the separate API
	if d.HasChange("default_node_pool") || d.HasChange("agent_pool_profile") {
		log.Printf("[DEBUG] Updating of Default Node Pool..")

		agentProfiles, err := containers.ExpandDefaultNodePool(d)
		if err != nil {
			return fmt.Errorf("Error expanding `default_node_pool`: %+v", err)
		}

		// TODO: remove me in 2.0
		if agentProfiles == nil {
			agentProfilesRaw := d.Get("agent_pool_profile").([]interface{})
			agentProfilesLegacy, err := expandKubernetesClusterAgentPoolProfiles(agentProfilesRaw, false)
			if err != nil {
				return err
			}

			agentProfiles = &agentProfilesLegacy
		}

		agentProfile := containers.ConvertDefaultNodePoolToAgentPool(agentProfiles)
		agentPool, err := nodePoolsClient.CreateOrUpdate(ctx, resourceGroup, name, *agentProfile.Name, agentProfile)
		if err != nil {
			return fmt.Errorf("Error updating Default Node Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		if err := agentPool.WaitForCompletionRef(ctx, nodePoolsClient.Client); err != nil {
			return fmt.Errorf("Error waiting for update of Default Node Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
		log.Printf("[DEBUG] Updated Default Node Pool.")
	}

	// then roll the version of Kubernetes if necessary
	if d.HasChange("kubernetes_version") {
		existing, err = clusterClient.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Error retrieving existing Kubernetes Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
		if existing.ManagedClusterProperties == nil {
			return fmt.Errorf("Error retrieving existing Kubernetes Cluster %q (Resource Group %q): `properties` was nil", name, resourceGroup)
		}

		kubernetesVersion := d.Get("kubernetes_version").(string)
		log.Printf("[DEBUG] Upgrading the version of Kubernetes to %q..", kubernetesVersion)
		existing.ManagedClusterProperties.KubernetesVersion = utils.String(kubernetesVersion)

		future, err := clusterClient.CreateOrUpdate(ctx, resourceGroup, name, existing)
		if err != nil {
			return fmt.Errorf("Error updating Managed Kubernetes Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, clusterClient.Client); err != nil {
			return fmt.Errorf("Error waiting for update of Managed Kubernetes Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		log.Printf("[DEBUG] Upgraded the version of Kubernetes to %q..", kubernetesVersion)
	}

	d.Partial(false)

	return resourceArmKubernetesClusterRead(d, meta)
}

func resourceArmKubernetesClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Containers.KubernetesClustersClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := containers.ParseKubernetesClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Managed Kubernetes Cluster %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Managed Kubernetes Cluster %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	profile, err := client.GetAccessProfile(ctx, id.ResourceGroup, id.Name, "clusterUser")
	if err != nil {
		return fmt.Errorf("Error retrieving Access Profile for Managed Kubernetes Cluster %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
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

		addonProfiles := containers.FlattenKubernetesAddOnProfiles(props.AddonProfiles)
		if err := d.Set("addon_profile", addonProfiles); err != nil {
			return fmt.Errorf("Error setting `addon_profile`: %+v", err)
		}

		// TODO: remove me in 2.0
		agentPoolProfiles := flattenKubernetesClusterAgentPoolProfiles(props.AgentPoolProfiles, resp.Fqdn)
		if err := d.Set("agent_pool_profile", agentPoolProfiles); err != nil {
			return fmt.Errorf("Error setting `agent_pool_profile`: %+v", err)
		}

		flattenedDefaultNodePool, err := containers.FlattenDefaultNodePool(props.AgentPoolProfiles, d)
		if err != nil {
			return fmt.Errorf("Error flattening `default_node_pool`: %+v", err)
		}
		if err := d.Set("default_node_pool", flattenedDefaultNodePool); err != nil {
			return fmt.Errorf("Error setting `default_node_pool`: %+v", err)
		}

		linuxProfile := flattenKubernetesClusterLinuxProfile(props.LinuxProfile)
		if err := d.Set("linux_profile", linuxProfile); err != nil {
			return fmt.Errorf("Error setting `linux_profile`: %+v", err)
		}

		networkProfile := flattenKubernetesClusterNetworkProfile(props.NetworkProfile)
		if err := d.Set("network_profile", networkProfile); err != nil {
			return fmt.Errorf("Error setting `network_profile`: %+v", err)
		}

		roleBasedAccessControl := flattenKubernetesClusterRoleBasedAccessControl(props, d)
		if err := d.Set("role_based_access_control", roleBasedAccessControl); err != nil {
			return fmt.Errorf("Error setting `role_based_access_control`: %+v", err)
		}

		servicePrincipal := flattenAzureRmKubernetesClusterServicePrincipalProfile(props.ServicePrincipalProfile, d)
		if err := d.Set("service_principal", servicePrincipal); err != nil {
			return fmt.Errorf("Error setting `service_principal`: %+v", err)
		}

		windowsProfile := flattenKubernetesClusterWindowsProfile(props.WindowsProfile, d)
		if err := d.Set("windows_profile", windowsProfile); err != nil {
			return fmt.Errorf("Error setting `windows_profile`: %+v", err)
		}

		// adminProfile is only available for RBAC enabled clusters with AAD
		if props.AadProfile != nil {
			adminProfile, err := client.GetAccessProfile(ctx, id.ResourceGroup, id.Name, "clusterAdmin")
			if err != nil {
				return fmt.Errorf("Error retrieving Admin Access Profile for Managed Kubernetes Cluster %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
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
	client := meta.(*ArmClient).Containers.KubernetesClustersClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := containers.ParseKubernetesClusterID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting Managed Kubernetes Cluster %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deletion of Managed Kubernetes Cluster %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
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

func expandKubernetesClusterAgentPoolProfiles(input []interface{}, isNewResource bool) ([]containerservice.ManagedClusterAgentPoolProfile, error) {
	profiles := make([]containerservice.ManagedClusterAgentPoolProfile, 0)

	for _, v := range input {
		config := v.(map[string]interface{})

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
			if *profile.EnableAutoScaling && !isNewResource {
				profile.Count = nil
			}
		}

		if availabilityZones := utils.ExpandStringSlice(config["availability_zones"].([]interface{})); len(*availabilityZones) > 0 {
			profile.AvailabilityZones = availabilityZones
		}

		if *profile.EnableAutoScaling && (profile.MinCount == nil || profile.MaxCount == nil) {
			return nil, fmt.Errorf("Can't create an AKS cluster with autoscaling enabled but not setting min_count or max_count")
		}

		if nodeTaints := utils.ExpandStringSlice(config["node_taints"].([]interface{})); len(*nodeTaints) > 0 {
			profile.NodeTaints = nodeTaints
		}

		if enableNodePublicIP := config["enable_node_public_ip"]; enableNodePublicIP != nil {
			profile.EnableNodePublicIP = utils.Bool(enableNodePublicIP.(bool))
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
		count := 0
		if profile.Count != nil {
			count = int(*profile.Count)
		}

		enableAutoScaling := false
		if profile.EnableAutoScaling != nil {
			enableAutoScaling = *profile.EnableAutoScaling
		}

		fqdnVal := ""
		if fqdn != nil {
			// temporarily persist the parent FQDN here until `fqdn` is removed from the `agent_pool_profile`
			fqdnVal = *fqdn
		}

		maxCount := 0
		if profile.MaxCount != nil {
			maxCount = int(*profile.MaxCount)
		}

		maxPods := 0
		if profile.MaxPods != nil {
			maxPods = int(*profile.MaxPods)
		}

		minCount := 0
		if profile.MinCount != nil {
			minCount = int(*profile.MinCount)
		}

		name := ""
		if profile.Name != nil {
			name = *profile.Name
		}

		osDiskSizeGB := 0
		if profile.OsDiskSizeGB != nil {
			osDiskSizeGB = int(*profile.OsDiskSizeGB)
		}

		subnetId := ""
		if profile.VnetSubnetID != nil {
			subnetId = *profile.VnetSubnetID
		}

		enableNodePublicIP := false
		if profile.EnableNodePublicIP != nil {
			enableNodePublicIP = *profile.EnableNodePublicIP
		}

		agentPoolProfile := map[string]interface{}{
			"availability_zones":    utils.FlattenStringSlice(profile.AvailabilityZones),
			"count":                 count,
			"enable_auto_scaling":   enableAutoScaling,
			"enable_node_public_ip": enableNodePublicIP,
			"max_count":             maxCount,
			"max_pods":              maxPods,
			"min_count":             minCount,
			"name":                  name,
			"node_taints":           utils.FlattenStringSlice(profile.NodeTaints),
			"os_disk_size_gb":       osDiskSizeGB,
			"os_type":               string(profile.OsType),
			"type":                  string(profile.Type),
			"vm_size":               string(profile.VMSize),
			"vnet_subnet_id":        subnetId,

			// TODO: remove in 2.0
			"fqdn": fqdnVal,
		}

		agentPoolProfiles = append(agentPoolProfiles, agentPoolProfile)
	}

	return agentPoolProfiles
}

func expandKubernetesClusterLinuxProfile(input []interface{}) *containerservice.LinuxProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]interface{})

	adminUsername := config["admin_username"].(string)
	linuxKeys := config["ssh_key"].([]interface{})

	keyData := ""
	if key, ok := linuxKeys[0].(map[string]interface{}); ok {
		keyData = key["key_data"].(string)
	}

	return &containerservice.LinuxProfile{
		AdminUsername: &adminUsername,
		SSH: &containerservice.SSHConfiguration{
			PublicKeys: &[]containerservice.SSHPublicKey{
				{
					KeyData: &keyData,
				},
			},
		},
	}
}

func flattenKubernetesClusterLinuxProfile(profile *containerservice.LinuxProfile) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	adminUsername := ""
	if username := profile.AdminUsername; username != nil {
		adminUsername = *username
	}

	sshKeys := make([]interface{}, 0)
	if ssh := profile.SSH; ssh != nil {
		if keys := ssh.PublicKeys; keys != nil {
			for _, sshKey := range *keys {
				keyData := ""
				if kd := sshKey.KeyData; kd != nil {
					keyData = *kd
				}
				sshKeys = append(sshKeys, map[string]interface{}{
					"key_data": keyData,
				})
			}
		}
	}

	return []interface{}{
		map[string]interface{}{
			"admin_username": adminUsername,
			"ssh_key":        sshKeys,
		},
	}
}

func expandKubernetesClusterWindowsProfile(input []interface{}) *containerservice.ManagedClusterWindowsProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]interface{})

	adminUsername := config["admin_username"].(string)
	adminPassword := config["admin_password"].(string)

	profile := containerservice.ManagedClusterWindowsProfile{
		AdminUsername: &adminUsername,
		AdminPassword: &adminPassword,
	}

	return &profile
}

func flattenKubernetesClusterWindowsProfile(profile *containerservice.ManagedClusterWindowsProfile, d *schema.ResourceData) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	adminUsername := ""
	if username := profile.AdminUsername; username != nil {
		adminUsername = *username
	}

	// admin password isn't returned, so let's look it up
	adminPassword := ""
	if v, ok := d.GetOk("windows_profile.0.admin_password"); ok {
		adminPassword = v.(string)
	}

	return []interface{}{
		map[string]interface{}{
			"admin_password": adminPassword,
			"admin_username": adminUsername,
		},
	}
}

func expandKubernetesClusterNetworkProfile(input []interface{}) *containerservice.NetworkProfileType {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]interface{})

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

	dnsServiceIP := ""
	if profile.DNSServiceIP != nil {
		dnsServiceIP = *profile.DNSServiceIP
	}

	dockerBridgeCidr := ""
	if profile.DockerBridgeCidr != nil {
		dockerBridgeCidr = *profile.DockerBridgeCidr
	}

	serviceCidr := ""
	if profile.ServiceCidr != nil {
		serviceCidr = *profile.ServiceCidr
	}

	podCidr := ""
	if profile.PodCidr != nil {
		podCidr = *profile.PodCidr
	}

	return []interface{}{
		map[string]interface{}{
			"dns_service_ip":     dnsServiceIP,
			"docker_bridge_cidr": dockerBridgeCidr,
			"load_balancer_sku":  string(profile.LoadBalancerSku),
			"network_plugin":     string(profile.NetworkPlugin),
			"network_policy":     string(profile.NetworkPolicy),
			"pod_cidr":           podCidr,
			"service_cidr":       serviceCidr,
		},
	}
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
		clientAppId := ""
		if profile.ClientAppID != nil {
			clientAppId = *profile.ClientAppID
		}

		serverAppId := ""
		if profile.ServerAppID != nil {
			serverAppId = *profile.ServerAppID
		}

		serverAppSecret := ""
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
						serverAppSecret = v.(string)
					}
				}
			}
		}

		tenantId := ""
		if profile.TenantID != nil {
			tenantId = *profile.TenantID
		}

		results = append(results, map[string]interface{}{
			"client_app_id":     clientAppId,
			"server_app_id":     serverAppId,
			"server_app_secret": serverAppSecret,
			"tenant_id":         tenantId,
		})
	}

	return []interface{}{
		map[string]interface{}{
			"enabled":                rbacEnabled,
			"azure_active_directory": results,
		},
	}
}

func flattenAzureRmKubernetesClusterServicePrincipalProfile(profile *containerservice.ManagedClusterServicePrincipalProfile, d *schema.ResourceData) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	clientId := ""
	if v := profile.ClientID; v != nil {
		clientId = *v
	}

	// client secret isn't returned by the API so pass the existing value along
	clientSecret := ""
	if sp, ok := d.GetOk("service_principal"); ok {
		var val []interface{}

		// prior to 1.34 this was a *schema.Set, now it's a List - try both
		if v, ok := sp.([]interface{}); ok {
			val = v
		} else if v, ok := sp.(*schema.Set); ok {
			val = v.List()
		}

		if len(val) > 0 && val[0] != nil {
			raw := val[0].(map[string]interface{})
			clientSecret = raw["client_secret"].(string)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"client_id":     clientId,
			"client_secret": clientSecret,
		},
	}
}

func flattenKubernetesClusterKubeConfig(config kubernetes.KubeConfig) []interface{} {
	// we don't size-check these since they're validated in the Parse method
	cluster := config.Clusters[0].Cluster
	user := config.Users[0].User
	name := config.Users[0].Name

	return []interface{}{
		map[string]interface{}{
			"client_certificate":     user.ClientCertificteData,
			"client_key":             user.ClientKeyData,
			"cluster_ca_certificate": cluster.ClusterAuthorityData,
			"host":                   cluster.Server,
			"password":               user.Token,
			"username":               name,
		},
	}
}

func flattenKubernetesClusterKubeConfigAAD(config kubernetes.KubeConfigAAD) []interface{} {
	// we don't size-check these since they're validated in the Parse method
	cluster := config.Clusters[0].Cluster
	name := config.Users[0].Name

	return []interface{}{
		map[string]interface{}{
			"client_certificate":     "",
			"client_key":             "",
			"cluster_ca_certificate": cluster.ClusterAuthorityData,
			"host":                   cluster.Server,
			"password":               "",
			"username":               name,
		},
	}
}
