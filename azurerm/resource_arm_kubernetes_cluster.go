package azurerm

import (
	"bytes"
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2017-09-30/containerservice"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/kubernetes"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKubernetesClusterCreate,
		Read:   resourceArmKubernetesClusterRead,
		Update: resourceArmKubernetesClusterCreate,
		Delete: resourceArmKubernetesClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"dns_prefix": {
				Type:     schema.TypeString,
				Required: true,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"kubernetes_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

			"linux_profile": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admin_username": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"ssh_key": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key_data": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},

			"agent_pool_profile": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validateKubernetesClusterAgentPoolName(),
						},

						"count": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(1, 50),
						},

						"dns_prefix": {
							Type:     schema.TypeString,
							Computed: true,
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
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"os_disk_size_gb": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},

						"vnet_subnet_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"os_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  containerservice.Linux,
							ValidateFunc: validation.StringInSlice([]string{
								string(containerservice.Linux),
								string(containerservice.Windows),
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},
					},
				},
			},

			"service_principal": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Type:     schema.TypeString,
							Required: true,
						},

						"client_secret": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
				Set: resourceAzureRMKubernetesClusterServicePrincipalProfileHash,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmKubernetesClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	kubernetesClustersClient := client.kubernetesClustersClient

	log.Printf("[INFO] preparing arguments for Azure ARM AKS managed cluster creation.")

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	dnsPrefix := d.Get("dns_prefix").(string)
	kubernetesVersion := d.Get("kubernetes_version").(string)

	linuxProfile := expandAzureRmKubernetesClusterLinuxProfile(d)
	agentProfiles := expandAzureRmKubernetesClusterAgentProfiles(d)
	servicePrincipalProfile := expandAzureRmKubernetesClusterServicePrincipal(d)

	tags := d.Get("tags").(map[string]interface{})

	parameters := containerservice.ManagedCluster{
		Name:     &name,
		Location: &location,
		ManagedClusterProperties: &containerservice.ManagedClusterProperties{
			AgentPoolProfiles:       &agentProfiles,
			DNSPrefix:               &dnsPrefix,
			KubernetesVersion:       &kubernetesVersion,
			LinuxProfile:            &linuxProfile,
			ServicePrincipalProfile: servicePrincipalProfile,
		},
		Tags: expandTags(tags),
	}

	ctx := client.StopContext
	future, err := kubernetesClustersClient.CreateOrUpdate(ctx, resGroup, name, parameters)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, kubernetesClustersClient.Client)
	if err != nil {
		return err
	}

	read, err := kubernetesClustersClient.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read AKS Managed Cluster %q (Resource Group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmKubernetesClusterRead(d, meta)
}

func resourceArmKubernetesClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	kubernetesClustersClient := meta.(*ArmClient).kubernetesClustersClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["managedClusters"]

	ctx := client.StopContext
	resp, err := kubernetesClustersClient.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on AKS Managed Cluster %q (resource group %q): %+v", name, resGroup, err)
	}

	profile, err := kubernetesClustersClient.GetAccessProfiles(ctx, resGroup, name, "clusterUser")
	if err != nil {
		return fmt.Errorf("Error getting access profile while making Read request on AKS Managed Cluster %q (resource group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.ManagedClusterProperties; props != nil {
		d.Set("dns_prefix", props.DNSPrefix)
		d.Set("fqdn", props.Fqdn)
		d.Set("kubernetes_version", props.KubernetesVersion)

		linuxProfile := flattenAzureRmKubernetesClusterLinuxProfile(props.LinuxProfile)
		if err := d.Set("linux_profile", linuxProfile); err != nil {
			return fmt.Errorf("Error setting `linux_profile`: %+v", err)
		}

		agentPoolProfiles := flattenAzureRmKubernetesClusterAgentPoolProfiles(props.AgentPoolProfiles, resp.Fqdn)
		if err := d.Set("agent_pool_profile", agentPoolProfiles); err != nil {
			return fmt.Errorf("Error setting `agent_pool_profile`: %+v", err)
		}

		servicePrincipal := flattenAzureRmKubernetesClusterServicePrincipalProfile(resp.ManagedClusterProperties.ServicePrincipalProfile)
		if err := d.Set("service_principal", servicePrincipal); err != nil {
			return fmt.Errorf("Error setting `service_principal`: %+v", err)
		}
	}

	kubeConfigRaw, kubeConfig := flattenAzureRmKubernetesClusterAccessProfile(&profile)
	d.Set("kube_config_raw", kubeConfigRaw)

	if err := d.Set("kube_config", kubeConfig); err != nil {
		return fmt.Errorf("Error setting `kube_config`: %+v", err)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmKubernetesClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	kubernetesClustersClient := client.kubernetesClustersClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["managedClusters"]

	ctx := client.StopContext
	future, err := kubernetesClustersClient.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error issuing AzureRM delete request of AKS Managed Cluster %q (resource Group %q): %+v", name, resGroup, err)
	}

	return future.WaitForCompletion(ctx, kubernetesClustersClient.Client)
}

func flattenAzureRmKubernetesClusterLinuxProfile(input *containerservice.LinuxProfile) []interface{} {
	values := make(map[string]interface{})
	sshKeys := make([]interface{}, 0)

	if profile := input; profile != nil {
		if username := profile.AdminUsername; username != nil {
			values["admin_username"] = *username
		}

		if ssh := profile.SSH; ssh != nil {
			if keys := ssh.PublicKeys; keys != nil {
				for _, sshKey := range *keys {
					outputs := make(map[string]interface{}, 0)
					if keyData := sshKey.KeyData; keyData != nil {
						outputs["key_data"] = *keyData
					}
					sshKeys = append(sshKeys, outputs)
				}
			}
		}
	}
	values["ssh_key"] = sshKeys

	return []interface{}{values}
}

func flattenAzureRmKubernetesClusterAgentPoolProfiles(profiles *[]containerservice.AgentPoolProfile, fqdn *string) []interface{} {
	agentPoolProfiles := make([]interface{}, 0)

	for _, profile := range *profiles {
		agentPoolProfile := make(map[string]interface{})

		if profile.Count != nil {
			agentPoolProfile["count"] = int(*profile.Count)
		}

		if profile.DNSPrefix != nil {
			agentPoolProfile["dns_prefix"] = *profile.DNSPrefix
		}

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

		agentPoolProfiles = append(agentPoolProfiles, agentPoolProfile)
	}

	return agentPoolProfiles
}

func flattenAzureRmKubernetesClusterServicePrincipalProfile(profile *containerservice.ServicePrincipalProfile) *schema.Set {
	if profile == nil {
		return nil
	}

	servicePrincipalProfiles := &schema.Set{
		F: resourceAzureRMKubernetesClusterServicePrincipalProfileHash,
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

func flattenAzureRmKubernetesClusterAccessProfile(profile *containerservice.ManagedClusterAccessProfile) (*string, []interface{}) {
	if profile != nil {
		if accessProfile := profile.AccessProfile; accessProfile != nil {
			if kubeConfigRaw := accessProfile.KubeConfig; kubeConfigRaw != nil {
				rawConfig := string(*kubeConfigRaw)

				kubeConfig, err := kubernetes.ParseKubeConfig(rawConfig)
				if err != nil {
					return utils.String(rawConfig), []interface{}{}
				}

				flattenedKubeConfig := flattenKubernetesClusterKubeConfig(*kubeConfig)
				return utils.String(rawConfig), flattenedKubeConfig
			}
		}
	}
	return nil, []interface{}{}
}

func flattenKubernetesClusterKubeConfig(config kubernetes.KubeConfig) []interface{} {
	values := make(map[string]interface{})

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

func expandAzureRmKubernetesClusterLinuxProfile(d *schema.ResourceData) containerservice.LinuxProfile {
	profiles := d.Get("linux_profile").([]interface{})
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

	return profile
}

func expandAzureRmKubernetesClusterServicePrincipal(d *schema.ResourceData) *containerservice.ServicePrincipalProfile {
	value, exists := d.GetOk("service_principal")
	if !exists {
		return nil
	}

	configs := value.(*schema.Set).List()

	config := configs[0].(map[string]interface{})

	clientId := config["client_id"].(string)
	clientSecret := config["client_secret"].(string)

	principal := containerservice.ServicePrincipalProfile{
		ClientID: &clientId,
		Secret:   &clientSecret,
	}

	return &principal
}

func expandAzureRmKubernetesClusterAgentProfiles(d *schema.ResourceData) []containerservice.AgentPoolProfile {
	configs := d.Get("agent_pool_profile").([]interface{})
	config := configs[0].(map[string]interface{})
	profiles := make([]containerservice.AgentPoolProfile, 0, len(configs))

	name := config["name"].(string)
	count := int32(config["count"].(int))
	dnsPrefix := config["dns_prefix"].(string)
	vmSize := config["vm_size"].(string)
	osDiskSizeGB := int32(config["os_disk_size_gb"].(int))
	osType := config["os_type"].(string)

	profile := containerservice.AgentPoolProfile{
		Name:           utils.String(name),
		Count:          utils.Int32(count),
		VMSize:         containerservice.VMSizeTypes(vmSize),
		DNSPrefix:      utils.String(dnsPrefix),
		OsDiskSizeGB:   utils.Int32(osDiskSizeGB),
		StorageProfile: containerservice.ManagedDisks,
		OsType:         containerservice.OSType(osType),
	}

	vnetSubnetID := config["vnet_subnet_id"].(string)
	if vnetSubnetID != "" {
		profile.VnetSubnetID = utils.String(vnetSubnetID)
	}

	profiles = append(profiles, profile)

	return profiles
}

func resourceAzureRMKubernetesClusterServicePrincipalProfileHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["client_id"].(string)))
	}

	return hashcode.String(buf.String())
}

func validateKubernetesClusterAgentPoolName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-z]{1}[a-z0-9]{0,11}$"),
		"Agent Pool names must start with a lowercase letter, have max length of 12, and only have characters a-z0-9.",
	)
}
