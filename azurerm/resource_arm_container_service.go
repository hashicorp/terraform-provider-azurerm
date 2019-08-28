package azurerm

import (
	"fmt"
	"log"

	"time"

	"bytes"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-06-01/containerservice"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmContainerService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmContainerServiceCreateUpdate,
		Read:   resourceArmContainerServiceRead,
		Update: resourceArmContainerServiceCreateUpdate,
		Delete: resourceArmContainerServiceDelete,

		DeprecationMessage: `Azure Container Service (ACS) has been deprecated in favour of Azure (Managed) Kubernetes Service (AKS).

Azure will remove support for ACS Clusters on January 31, 2020. In preparation for this, the AzureRM Provider will remove support for the 'azurerm_container_service' resource in the next major version of the AzureRM Provider, which is targeted for Early 2019.

If you're using ACS with Kubernetes, we'd recommend migrating to AKS / the 'azurerm_kubernetes_cluster' resource.

More information can be found here: https://azure.microsoft.com/en-us/updates/azure-container-service-will-retire-on-january-31-2020/
`,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"orchestration_platform": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArmContainerServiceOrchestrationPlatform,
			},

			"master_profile": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validateArmContainerServiceMasterProfileCount,
						},

						"dns_prefix": {
							Type:     schema.TypeString,
							Required: true,
						},

						"fqdn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: resourceAzureRMContainerServiceMasterProfileHash,
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
				Set: resourceAzureRMContainerServiceLinuxProfilesHash,
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
							ValidateFunc: validateArmContainerServiceAgentPoolProfileCount,
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
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},
				},
				Set: resourceAzureRMContainerServiceAgentPoolProfilesHash,
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
				Set: resourceAzureRMContainerServiceServicePrincipalProfileHash,
			},

			"diagnostics_profile": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"storage_uri": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: resourceAzureRMContainerServiceDiagnosticProfilesHash,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmContainerServiceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	ctx := meta.(*ArmClient).StopContext
	containerServiceClient := client.containers.ServicesClient

	log.Printf("[INFO] preparing arguments for Azure ARM Container Service creation.")

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := containerServiceClient.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Container Service %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_container_service", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	orchestrationPlatform := d.Get("orchestration_platform").(string)

	masterProfile := expandAzureRmContainerServiceMasterProfile(d)
	linuxProfile := expandAzureRmContainerServiceLinuxProfile(d)
	agentProfiles := expandAzureRmContainerServiceAgentProfiles(d)
	diagnosticsProfile := expandAzureRmContainerServiceDiagnostics(d)

	t := d.Get("tags").(map[string]interface{})

	parameters := containerservice.ContainerService{
		Name:     &name,
		Location: &location,
		Properties: &containerservice.Properties{
			MasterProfile: &masterProfile,
			LinuxProfile:  &linuxProfile,
			OrchestratorProfile: &containerservice.OrchestratorProfileType{
				OrchestratorType: containerservice.OrchestratorTypes(orchestrationPlatform),
			},
			AgentPoolProfiles:  &agentProfiles,
			DiagnosticsProfile: &diagnosticsProfile,
		},
		Tags: tags.Expand(t),
	}

	servicePrincipalProfile := expandAzureRmContainerServiceServicePrincipal(d)
	if servicePrincipalProfile != nil {
		parameters.ServicePrincipalProfile = servicePrincipalProfile
	}

	if _, err := containerServiceClient.CreateOrUpdate(ctx, resGroup, name, parameters); err != nil {
		return err
	}

	read, err := containerServiceClient.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Container Service %s (resource group %s) ID", name, resGroup)
	}

	log.Printf("[DEBUG] Waiting for Container Service (%s) to become available", d.Get("name"))
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Updating", "Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    containerServiceStateRefreshFunc(client, resGroup, name),
		Timeout:    30 * time.Minute,
		MinTimeout: 15 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Container Service (%s) to become available: %s", d.Get("name"), err)
	}

	d.SetId(*read.ID)

	return resourceArmContainerServiceRead(d, meta)
}

func resourceArmContainerServiceRead(d *schema.ResourceData, meta interface{}) error {
	containerServiceClient := meta.(*ArmClient).containers.ServicesClient

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["containerServices"]

	ctx := meta.(*ArmClient).StopContext
	resp, err := containerServiceClient.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Container Service %s: %s", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("orchestration_platform", string(resp.Properties.OrchestratorProfile.OrchestratorType))

	masterProfiles := flattenAzureRmContainerServiceMasterProfile(*resp.Properties.MasterProfile)
	d.Set("master_profile", &masterProfiles)

	linuxProfile := flattenAzureRmContainerServiceLinuxProfile(*resp.Properties.LinuxProfile)
	d.Set("linux_profile", &linuxProfile)

	agentPoolProfiles := flattenAzureRmContainerServiceAgentPoolProfiles(resp.Properties.AgentPoolProfiles)
	d.Set("agent_pool_profile", &agentPoolProfiles)

	servicePrincipal := flattenAzureRmContainerServiceServicePrincipalProfile(resp.Properties.ServicePrincipalProfile)
	if servicePrincipal != nil {
		d.Set("service_principal", servicePrincipal)
	}

	diagnosticProfile := flattenAzureRmContainerServiceDiagnosticsProfile(resp.Properties.DiagnosticsProfile)
	if diagnosticProfile != nil {
		d.Set("diagnostics_profile", diagnosticProfile)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmContainerServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	containerServiceClient := client.containers.ServicesClient

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["containerServices"]

	ctx := meta.(*ArmClient).StopContext
	future, err := containerServiceClient.Delete(ctx, resGroup, name)

	if err != nil {
		return fmt.Errorf("Error issuing Azure ARM delete request of Container Service '%s': %s", name, err)
	}

	return future.WaitForCompletionRef(ctx, containerServiceClient.Client)
}

func flattenAzureRmContainerServiceMasterProfile(profile containerservice.MasterProfile) *schema.Set {
	masterProfiles := &schema.Set{
		F: resourceAzureRMContainerServiceMasterProfileHash,
	}

	masterProfile := make(map[string]interface{}, 3)

	masterProfile["count"] = int(*profile.Count)
	masterProfile["dns_prefix"] = *profile.DNSPrefix
	masterProfile["fqdn"] = *profile.Fqdn

	masterProfiles.Add(masterProfile)

	return masterProfiles
}

func flattenAzureRmContainerServiceLinuxProfile(profile containerservice.LinuxProfile) *schema.Set {
	profiles := &schema.Set{
		F: resourceAzureRMContainerServiceLinuxProfilesHash,
	}

	values := map[string]interface{}{}

	sshKeys := &schema.Set{
		F: resourceAzureRMContainerServiceLinuxProfilesSSHKeysHash,
	}
	for _, ssh := range *profile.SSH.PublicKeys {
		keys := map[string]interface{}{}
		keys["key_data"] = *ssh.KeyData
		sshKeys.Add(keys)
	}

	values["admin_username"] = *profile.AdminUsername
	values["ssh_key"] = sshKeys
	profiles.Add(values)

	return profiles
}

func flattenAzureRmContainerServiceAgentPoolProfiles(profiles *[]containerservice.AgentPoolProfile) *schema.Set {
	agentPoolProfiles := &schema.Set{
		F: resourceAzureRMContainerServiceAgentPoolProfilesHash,
	}

	for _, profile := range *profiles {
		agentPoolProfile := map[string]interface{}{}
		agentPoolProfile["count"] = int(*profile.Count)
		agentPoolProfile["dns_prefix"] = *profile.DNSPrefix
		agentPoolProfile["fqdn"] = *profile.Fqdn
		agentPoolProfile["name"] = *profile.Name
		agentPoolProfile["vm_size"] = string(profile.VMSize)
		agentPoolProfiles.Add(agentPoolProfile)
	}

	return agentPoolProfiles
}

func flattenAzureRmContainerServiceServicePrincipalProfile(profile *containerservice.ServicePrincipalProfile) *schema.Set {

	if profile == nil {
		return nil
	}

	servicePrincipalProfiles := &schema.Set{
		F: resourceAzureRMContainerServiceServicePrincipalProfileHash,
	}

	values := map[string]interface{}{}

	values["client_id"] = *profile.ClientID
	if profile.Secret != nil {
		values["client_secret"] = *profile.Secret
	}

	servicePrincipalProfiles.Add(values)

	return servicePrincipalProfiles
}

func flattenAzureRmContainerServiceDiagnosticsProfile(profile *containerservice.DiagnosticsProfile) *schema.Set {
	diagnosticProfiles := &schema.Set{
		F: resourceAzureRMContainerServiceDiagnosticProfilesHash,
	}

	values := map[string]interface{}{}

	values["enabled"] = *profile.VMDiagnostics.Enabled
	if profile.VMDiagnostics.StorageURI != nil {
		values["storage_uri"] = *profile.VMDiagnostics.StorageURI
	}
	diagnosticProfiles.Add(values)

	return diagnosticProfiles
}

func expandAzureRmContainerServiceDiagnostics(d *schema.ResourceData) containerservice.DiagnosticsProfile {
	configs := d.Get("diagnostics_profile").(*schema.Set).List()

	data := configs[0].(map[string]interface{})

	enabled := data["enabled"].(bool)

	return containerservice.DiagnosticsProfile{
		VMDiagnostics: &containerservice.VMDiagnostics{
			Enabled: &enabled,
		},
	}
}

func expandAzureRmContainerServiceLinuxProfile(d *schema.ResourceData) containerservice.LinuxProfile {
	profiles := d.Get("linux_profile").(*schema.Set).List()
	config := profiles[0].(map[string]interface{})

	adminUsername := config["admin_username"].(string)

	linuxKeys := config["ssh_key"].(*schema.Set).List()
	sshPublicKeys := make([]containerservice.SSHPublicKey, 0)

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

func expandAzureRmContainerServiceMasterProfile(d *schema.ResourceData) containerservice.MasterProfile {
	configs := d.Get("master_profile").(*schema.Set).List()
	config := configs[0].(map[string]interface{})

	count := int32(config["count"].(int))
	dnsPrefix := config["dns_prefix"].(string)

	profile := containerservice.MasterProfile{
		Count:     &count,
		DNSPrefix: &dnsPrefix,
	}

	return profile
}

func expandAzureRmContainerServiceServicePrincipal(d *schema.ResourceData) *containerservice.ServicePrincipalProfile {

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

func expandAzureRmContainerServiceAgentProfiles(d *schema.ResourceData) []containerservice.AgentPoolProfile {
	configs := d.Get("agent_pool_profile").(*schema.Set).List()
	config := configs[0].(map[string]interface{})
	profiles := make([]containerservice.AgentPoolProfile, 0, len(configs))

	name := config["name"].(string)
	count := int32(config["count"].(int))
	dnsPrefix := config["dns_prefix"].(string)
	vmSize := config["vm_size"].(string)

	profile := containerservice.AgentPoolProfile{
		Name:      &name,
		Count:     &count,
		VMSize:    containerservice.VMSizeTypes(vmSize),
		DNSPrefix: &dnsPrefix,
	}

	profiles = append(profiles, profile)

	return profiles
}

func containerServiceStateRefreshFunc(client *ArmClient, resourceGroupName string, containerServiceName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := client.StopContext
		res, err := client.containers.ServicesClient.Get(ctx, resourceGroupName, containerServiceName)
		if err != nil {
			return nil, "", fmt.Errorf("Error issuing read request in containerServiceStateRefreshFunc to Azure ARM for Container Service '%s' (RG: '%s'): %s", containerServiceName, resourceGroupName, err)
		}

		return res, *res.Properties.ProvisioningState, nil
	}
}

func resourceAzureRMContainerServiceMasterProfileHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%d-", m["count"].(int)))
		buf.WriteString(fmt.Sprintf("%s-", m["dns_prefix"].(string)))
	}

	return hashcode.String(buf.String())
}

func resourceAzureRMContainerServiceLinuxProfilesHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["admin_username"].(string)))
	}

	return hashcode.String(buf.String())
}

func resourceAzureRMContainerServiceLinuxProfilesSSHKeysHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["key_data"].(string)))
	}

	return hashcode.String(buf.String())
}

func resourceAzureRMContainerServiceAgentPoolProfilesHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%d-", m["count"].(int)))
		buf.WriteString(fmt.Sprintf("%s-", m["dns_prefix"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["vm_size"].(string)))
	}

	return hashcode.String(buf.String())
}

func resourceAzureRMContainerServiceServicePrincipalProfileHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["client_id"].(string)))
	}

	return hashcode.String(buf.String())
}

func resourceAzureRMContainerServiceDiagnosticProfilesHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%t", m["enabled"].(bool)))
	}

	return hashcode.String(buf.String())
}

func validateArmContainerServiceOrchestrationPlatform(v interface{}, _ string) (warnings []string, errors []error) {
	value := v.(string)
	capacities := map[string]bool{
		"DCOS":       true,
		"Kubernetes": true,
		"Swarm":      true,
	}

	if !capacities[value] {
		errors = append(errors, fmt.Errorf("Container Service: Orchestration Platgorm can only be DCOS / Kubernetes / Swarm"))
	}
	return warnings, errors
}

func validateArmContainerServiceMasterProfileCount(v interface{}, _ string) (warnings []string, errors []error) {
	value := v.(int)
	capacities := map[int]bool{
		1: true,
		3: true,
		5: true,
	}

	if !capacities[value] {
		errors = append(errors, fmt.Errorf("The number of master nodes must be 1, 3 or 5."))
	}
	return warnings, errors
}

func validateArmContainerServiceAgentPoolProfileCount(v interface{}, _ string) (warnings []string, errors []error) {
	value := v.(int)
	if value > 100 || 0 >= value {
		errors = append(errors, fmt.Errorf("The Count for an Agent Pool Profile can only be between 1 and 100."))
	}
	return warnings, errors
}
