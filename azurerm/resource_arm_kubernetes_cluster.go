package azurerm

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2017-09-30/containerservice"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKubernetesClusterCreate,
		Read:   resourceArmKubernetesClusterRead,
		Update: resourceArmKubernetesClusterCreate,
		Delete: resourceArmKubernetesClusterDelete,

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

			"kubernetes_version": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"linux_profile": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admin_username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ssh_key": {
							Type:     schema.TypeSet,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key_data": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
				Set: resourceAzureRMKubernetesClusterLinuxProfilesHash,
			},

			"agent_pool_profile": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"count": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(1, 100),
						},

						"dns_prefix": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"fqdn": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"vm_size": {
							Type:     schema.TypeString,
							Required: true,
						},

						"os_disk_size_gb": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},

						"storage_profile": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(containerservice.StorageAccount),
								string(containerservice.ManagedDisks),
							}, true),
						},

						"vnet_subnet_id": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"os_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  containerservice.Linux,
						},
					},
				},
				Set: resourceAzureRMKubernetesClusterAgentPoolProfilesHash,
			},

			"service_principal": {
				Type:     schema.TypeSet,
				Optional: true,
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
	location := d.Get("location").(string)
	dnsPrefix := d.Get("dns_prefix").(string)
	kubernetesVersion := d.Get("kubernetes_version").(string)

	linuxProfile := expandAzureRmKubernetesClusterLinuxProfile(d)
	agentProfiles := expandAzureRmKubernetesClusterAgentProfiles(d)

	tags := d.Get("tags").(map[string]interface{})

	parameters := containerservice.ManagedCluster{
		Name:     &name,
		Location: &location,
		ManagedClusterProperties: &containerservice.ManagedClusterProperties{
			DNSPrefix:         &dnsPrefix,
			KubernetesVersion: &kubernetesVersion,
			LinuxProfile:      &linuxProfile,
			AgentPoolProfiles: &agentProfiles,
		},
		Tags: expandTags(tags),
	}

	servicePrincipalProfile := expandAzureRmKubernetesClusterServicePrincipal(d)
	if servicePrincipalProfile != nil {
		parameters.ServicePrincipalProfile = servicePrincipalProfile
	}

	ctx := client.StopContext
	_, error := kubernetesClustersClient.CreateOrUpdate(ctx, resGroup, name, parameters)
	if error != nil {
		return error
	}

	read, err := kubernetesClustersClient.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read AKS managed cluster %s (resource group %s) ID", name, resGroup)
	}

	log.Printf("[DEBUG] Waiting for AKS managed cluster (%s) to become available", d.Get("name"))
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Updating", "Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    kubernetesClusterStateRefreshFunc(client, resGroup, name),
		Timeout:    30 * time.Minute,
		MinTimeout: 15 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for AKS managed cluster (%s) to become available: %s", d.Get("name"), err)
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

		return fmt.Errorf("Error making Read request on Azure Container Service %s: %s", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))
	d.Set("resource_group_name", resGroup)
	d.Set("dns_prefix", resp.DNSPrefix)
	d.Set("kubernetes_version", resp.KubernetesVersion)

	linuxProfile := flattenAzureRmKubernetesClusterLinuxProfile(*resp.ManagedClusterProperties.LinuxProfile)
	d.Set("linux_profile", &linuxProfile)

	agentPoolProfiles := flattenAzureRmKubernetesClusterAgentPoolProfiles(resp.ManagedClusterProperties.AgentPoolProfiles)
	d.Set("agent_pool_profile", &agentPoolProfiles)

	servicePrincipal := flattenAzureRmKubernetesClusterServicePrincipalProfile(resp.ManagedClusterProperties.ServicePrincipalProfile)
	if servicePrincipal != nil {
		d.Set("service_principal", servicePrincipal)
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
		return fmt.Errorf("Error issuing Azure ARM delete request of Container Service '%s': %s", name, err)
	}

	err = future.WaitForCompletion(ctx, kubernetesClustersClient.Client)
	if err != nil {
		return err
	}

	return nil
}

func flattenAzureRmKubernetesClusterLinuxProfile(profile containerservice.LinuxProfile) *schema.Set {
	profiles := &schema.Set{
		F: resourceAzureRMKubernetesClusterLinuxProfilesHash,
	}

	values := make(map[string]interface{})

	sshKeys := &schema.Set{
		F: resourceAzureRMKubernetesClusterLinuxProfilesSSHKeysHash,
	}
	for _, ssh := range *profile.SSH.PublicKeys {
		keys := make(map[string]interface{})
		keys["key_data"] = *ssh.KeyData
		sshKeys.Add(keys)
	}

	values["admin_username"] = *profile.AdminUsername
	values["ssh_key"] = sshKeys
	profiles.Add(values)

	return profiles
}

func flattenAzureRmKubernetesClusterAgentPoolProfiles(profiles *[]containerservice.AgentPoolProfile) *schema.Set {
	agentPoolProfiles := &schema.Set{
		F: resourceAzureRMKubernetesClusterAgentPoolProfilesHash,
	}

	for _, profile := range *profiles {
		agentPoolProfile := make(map[string]interface{})

		if profile.Count != nil {
			agentPoolProfile["count"] = int(*profile.Count)
		}

		if profile.DNSPrefix != nil {
			agentPoolProfile["dns_prefix"] = *profile.DNSPrefix
		}

		if profile.Fqdn != nil {
			agentPoolProfile["fqdn"] = *profile.Fqdn
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

		if profile.StorageProfile != "" {
			agentPoolProfile["storage_profile"] = string(profile.StorageProfile)
		}

		if profile.VnetSubnetID != nil {
			agentPoolProfile["vnet_subnet_id"] = *profile.VnetSubnetID
		}

		if profile.OsType != "" {
			agentPoolProfile["os_type"] = string(profile.OsType)
		}

		agentPoolProfiles.Add(agentPoolProfile)
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

	values["client_id"] = *profile.ClientID
	if profile.Secret != nil {
		values["client_secret"] = *profile.Secret
	}

	servicePrincipalProfiles.Add(values)

	return servicePrincipalProfiles
}

func expandAzureRmKubernetesClusterLinuxProfile(d *schema.ResourceData) containerservice.LinuxProfile {
	profiles := d.Get("linux_profile").(*schema.Set).List()
	config := profiles[0].(map[string]interface{})

	adminUsername := config["admin_username"].(string)

	linuxKeys := config["ssh_key"].(*schema.Set).List()
	sshPublicKeys := []containerservice.SSHPublicKey{}

	key := linuxKeys[0].(map[string]interface{})
	keyData := key["key_data"].(string)

	sshPublicKey := containerservice.SSHPublicKey{
		KeyData: &keyData,
	}

	sshPublicKeys = append(sshPublicKeys, sshPublicKey)

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
	configs := d.Get("agent_pool_profile").(*schema.Set).List()
	config := configs[0].(map[string]interface{})
	profiles := make([]containerservice.AgentPoolProfile, 0, len(configs))

	name := config["name"].(string)
	count := int32(config["count"].(int))
	dnsPrefix := config["dns_prefix"].(string)
	vmSize := config["vm_size"].(string)
	osDiskSizeGB := int32(config["os_disk_size_gb"].(int))
	storageProfile := config["storage_profile"].(string)
	vnetSubnetID := config["vnet_subnet_id"].(string)
	osType := config["os_type"].(string)

	profile := containerservice.AgentPoolProfile{
		Name:           &name,
		Count:          &count,
		VMSize:         containerservice.VMSizeTypes(vmSize),
		DNSPrefix:      &dnsPrefix,
		OsDiskSizeGB:   &osDiskSizeGB,
		StorageProfile: containerservice.StorageProfileTypes(storageProfile),
		VnetSubnetID:   &vnetSubnetID,
		OsType:         containerservice.OSType(osType),
	}

	profiles = append(profiles, profile)

	return profiles
}

func kubernetesClusterStateRefreshFunc(client *ArmClient, resourceGroupName string, kubernetesClusterName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := client.StopContext
		res, err := client.kubernetesClustersClient.Get(ctx, resourceGroupName, kubernetesClusterName)
		if err != nil {
			return nil, "", fmt.Errorf("Error issuing read request in AKSStateRefreshFunc to Azure ARM for AKS managed cluster '%s' (RG: '%s'): %s", kubernetesClusterName, resourceGroupName, err)
		}

		return res, *res.ManagedClusterProperties.ProvisioningState, nil
	}
}

func resourceAzureRMKubernetesClusterLinuxProfilesHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	adminUsername := m["admin_username"].(string)

	buf.WriteString(fmt.Sprintf("%s-", adminUsername))

	return hashcode.String(buf.String())
}

func resourceAzureRMKubernetesClusterLinuxProfilesSSHKeysHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	keyData := m["key_data"].(string)

	buf.WriteString(fmt.Sprintf("%s-", keyData))

	return hashcode.String(buf.String())
}

func resourceAzureRMKubernetesClusterAgentPoolProfilesHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if m["count"] != nil {
		buf.WriteString(fmt.Sprintf("%d-", m["count"].(int)))
	}

	if m["dns_prefix"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["dns_prefix"].(string)))
	}

	if m["name"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	}

	if m["vm_size"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["vm_size"].(string)))
	}

	if m["os_disk_size_gb"] != nil {
		buf.WriteString(fmt.Sprintf("%d-", m["os_disk_size_gb"].(int)))
	}

	if m["storage_profile"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["storage_profile"].(string)))
	}

	if m["vnet_subnet_id"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["vnet_subnet_id"].(string)))
	}

	if m["os_type"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["os_type"].(string)))
	}

	return hashcode.String(buf.String())
}

func resourceAzureRMKubernetesClusterServicePrincipalProfileHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	clientId := m["client_id"].(string)
	buf.WriteString(fmt.Sprintf("%s-", clientId))

	return hashcode.String(buf.String())
}
