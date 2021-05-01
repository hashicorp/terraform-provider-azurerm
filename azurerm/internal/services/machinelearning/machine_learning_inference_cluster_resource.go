package machinelearning

import (
	"fmt"
	"math"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2021-02-01/containerservice"
	"github.com/Azure/azure-sdk-for-go/services/machinelearningservices/mgmt/2020-04-01/machinelearningservices"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/machinelearning/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/machinelearning/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"

	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

const min_number_of_nodes_prod int32 = 12

func resourceAksInferenceCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAksInferenceClusterCreateUpdate,
		Read:   resourceAksInferenceClusterRead,
		Update: resourceAksInferenceClusterCreateUpdate,
		Delete: resourceAksInferenceClusterDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.InferenceClusterID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InferenceClusterName,
			},

			"workspace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"kubernetes_cluster_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"kubernetes_cluster_rg": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.KubernetesClusterResourceGroupName,
			},

			"cluster_purpose": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Dev",
				ValidateFunc: validate.ClusterPurpose,
			},

			"node_pool_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NodePoolName,
			},

			"identity": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(machinelearningservices.SystemAssigned),
							}, false),
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"ssl_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"ssl_certificate_custom": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"cname": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sku_name": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceAksInferenceClusterCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	machine_learning_workspaces_client := meta.(*clients.Client).MachineLearning.WorkspacesClient
	machine_learning_compute_client := meta.(*clients.Client).MachineLearning.MachineLearningComputeClient
	kubernetes_clusters_client := meta.(*clients.Client).Containers.KubernetesClustersClient
	node_pools_client := meta.(*clients.Client).Containers.AgentPoolsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	workspace_name := d.Get("workspace_name").(string)
	resource_group_name := d.Get("resource_group_name").(string)

	existing, err := machine_learning_compute_client.Get(ctx, resource_group_name, workspace_name, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("error checking for existing Inference Cluster %q in Workspace %q (Resource Group %q): %s", name, workspace_name, resource_group_name, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_machine_learning_inference_cluster", *existing.ID)
	}

	// Get SKU from Workspace
	aml_ws, _ := machine_learning_workspaces_client.Get(ctx, resource_group_name, workspace_name)

	sku := aml_ws.Sku

	kubernetes_cluster_name := d.Get("kubernetes_cluster_name").(string)
	kubernetes_cluster_rg := d.Get("kubernetes_cluster_rg").(string)

	// Get Existing AKS
	aks_cluster, _ := kubernetes_clusters_client.Get(ctx, kubernetes_cluster_rg, kubernetes_cluster_name)

	pool_name := d.Get("node_pool_name").(string)
	node_pool, _ := node_pools_client.Get(ctx, kubernetes_cluster_rg, kubernetes_cluster_name, pool_name)

	t := d.Get("tags").(map[string]interface{})

	identity := d.Get("identity")

	ssl_config := expandSSLConfig(d)
	cluster_purpose := d.Get("cluster_purpose").(string)
	var map_cluster_purpose = map[string]string{
		"Dev":  "DevTest",
		"Test": "DevTest",
		"Prod": "FastProd",
	}
	aks_cluster_purpose := machinelearningservices.ClusterPurpose(map_cluster_purpose[cluster_purpose])
	aks_properties := expandAksProperties(&aks_cluster, &node_pool, ssl_config, aks_cluster_purpose)

	location := azure.NormalizeLocation(d.Get("location").(string))
	description := d.Get("description").(string)
	aks_compute_properties := expandAksComputeProperties(aks_properties, &aks_cluster, &node_pool, location, description)
	compute_properties, is_aks := (machinelearningservices.BasicCompute).AsAKS(aks_compute_properties)

	if !is_aks {
		return fmt.Errorf("error: No AKS cluster")
	}

	inference_cluster_parameters := machinelearningservices.ComputeResource{
		// Properties - Compute properties
		Properties: compute_properties,
		// ID - READ-ONLY; Specifies the resource ID.
		// Name - READ-ONLY; Specifies the name of the resource.
		// Identity - The identity of the resource.
		Identity: expandMachineLearningComputeClusterIdentity(identity.([]interface{})),
		// Location - Specifies the location of the resource.
		Location: &location,
		// Type - READ-ONLY; Specifies the type of the resource.
		// Tags - Contains resource tags defined as key/value pairs.
		Tags: tags.Expand(t),
		// Sku - The sku of the workspace.
		Sku: sku,
	}

	if v, ok := d.GetOk("description"); ok {
		aks_compute_properties.Description = utils.String(v.(string))
	}

	future, ml_err := machine_learning_compute_client.CreateOrUpdate(ctx, resource_group_name, workspace_name, name, inference_cluster_parameters)
	if ml_err != nil {
		return fmt.Errorf("error creating Inference Cluster %q in workspace %q (Resource Group %q): %+v", name, workspace_name, resource_group_name, ml_err)
	}
	if err := future.WaitForCompletionRef(ctx, machine_learning_compute_client.Client); err != nil {
		return fmt.Errorf("error waiting for creation of Inference Cluster %q in workspace %q (Resource Group %q): %+v", name, workspace_name, resource_group_name, err)
	}
	resp, ml_get_err := machine_learning_compute_client.Get(ctx, resource_group_name, workspace_name, name)
	if ml_get_err != nil {
		return fmt.Errorf("error retrieving Inference Cluster Compute %q in workspace %q (Resource Group %q): %+v", name, workspace_name, resource_group_name, ml_get_err)
	}

	if resp.ID == nil {
		return fmt.Errorf("cannot read Inference Cluster ID %q in workspace %q (Resource Group %q) ID", name, workspace_name, resource_group_name)
	}

	d.SetId(*resp.ID)

	return resourceAksInferenceClusterRead(d, meta)
}

func resourceAksInferenceClusterRead(d *schema.ResourceData, meta interface{}) error {
	machine_learning_compute_client := meta.(*clients.Client).MachineLearning.MachineLearningComputeClient
	kubernetes_clusters_client := meta.(*clients.Client).Containers.KubernetesClustersClient
	node_pools_client := meta.(*clients.Client).Containers.AgentPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.InferenceClusterID(d.Id())
	if err != nil {
		return fmt.Errorf("error parsing Inference Cluster ID `%q`: %+v", d.Id(), err)
	}

	resp, err := machine_learning_compute_client.Get(ctx, id.ResourceGroup, id.Name, id.InferenceClusterName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error making Read request on Inference Cluster %q in Workspace %q (Resource Group %q): %+v", id.InferenceClusterName, id.Name, id.ResourceGroup, err)
	}

	d.Set("workspace_name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("name", id.InferenceClusterName)

	aks_resp, _ := kubernetes_clusters_client.ListByResourceGroup(ctx, id.ResourceGroup)
	kubernetes_cluster_name := aks_resp.Values()[0].Name
	resource_id := aks_resp.Values()[0].ID

	// get SSL configuration from ak1 below by getting an aks to pass to AsAKS by using the methods in create function... =
	node_pool_list, _ := node_pools_client.List(ctx, id.ResourceGroup, *kubernetes_cluster_name)
	pool_name := node_pool_list.Values()[0].Name

	d.Set("kubernetes_cluster_name", kubernetes_cluster_name)
	d.Set("kubernetes_cluster_rg", id.ResourceGroup)
	d.Set("resource_id", resource_id)
	d.Set("node_pool_name", pool_name)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	if err := d.Set("identity", flattenAksInferenceClusterIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("error flattening identity on Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceAksInferenceClusterDelete(d *schema.ResourceData, meta interface{}) error {
	machine_learning_compute_client := meta.(*clients.Client).MachineLearning.MachineLearningComputeClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := parse.InferenceClusterID(d.Id())
	if err != nil {
		return fmt.Errorf("error parsing Inference Cluster ID `%q`: %+v", d.Id(), err)
	}
	underlying_resource_action := machinelearningservices.Detach
	future, err := machine_learning_compute_client.Delete(ctx, id.ResourceGroup, id.Name, id.InferenceClusterName, underlying_resource_action)
	if err != nil {
		return fmt.Errorf("error deleting Inference Cluster %q in workspace %q (Resource Group %q): %+v", id.InferenceClusterName, id.Name, id.ResourceGroup, err)
	}
	if err := future.WaitForCompletionRef(ctx, machine_learning_compute_client.Client); err != nil {
		return fmt.Errorf("error waiting for deletion of Inference Cluster %q in workspace %q (Resource Group %q): %+v", id.InferenceClusterName, id.Name, id.ResourceGroup, err)
	}
	return nil
}

func expandSSLConfig(d *schema.ResourceData) *machinelearningservices.SslConfiguration {
	// Documentation: https://docs.microsoft.com/en-us/azure/machine-learning/how-to-secure-web-service
	var map_ssl_enabled = map[bool]string{
		true:  "Enabled",
		false: "Disabled",
	}

	// --- SSL Configurations ---
	ssl_enabled := map_ssl_enabled[d.Get("ssl_enabled").(bool)]

	// SSL Certificate default values
	cert_file := ""
	key_file := ""
	cname := ""
	leaf_domain_label := ""
	overwrite_existing_domain := false

	// SSL Custom Certificate settings
	ssl_certificate_custom_refs := d.Get("ssl_certificate_custom").(*schema.Set).List()

	if len(ssl_certificate_custom_refs) > 0 && ssl_certificate_custom_refs[0] != nil {
		ssl_certificate_custom_ref := ssl_certificate_custom_refs[0].(map[string]interface{})
		cert_file = ssl_certificate_custom_ref["cert"].(string)
		key_file = ssl_certificate_custom_ref["key"].(string)
		cname = ssl_certificate_custom_ref["cname"].(string)
	}

	return &machinelearningservices.SslConfiguration{
		Status:                  machinelearningservices.Status1(ssl_enabled),
		Cert:                    utils.String(cert_file),
		Key:                     utils.String(key_file),
		Cname:                   utils.String(cname),
		LeafDomainLabel:         utils.String(leaf_domain_label),
		OverwriteExistingDomain: utils.Bool(overwrite_existing_domain)}
}

func flattenCustomSSLConfig(ssl_config *machinelearningservices.SslConfiguration) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"cert":  *(ssl_config.Cert),
			"key":   *(ssl_config.Key),
			"cname": *(ssl_config.Cname),
		},
	}
}

func expandAksNetworkingConfiguration(aks *containerservice.ManagedCluster, node_pool *containerservice.AgentPool) *machinelearningservices.AksNetworkingConfiguration {
	subnet_id := *(node_pool.ManagedClusterAgentPoolProfileProperties).VnetSubnetID //d.Get("subnet_id").(string)
	service_cidr := *(aks.NetworkProfile.ServiceCidr)
	dns_service_ip := *(aks.NetworkProfile.DNSServiceIP)
	docker_bridge_cidr := *(aks.NetworkProfile.DockerBridgeCidr)

	return &machinelearningservices.AksNetworkingConfiguration{
		SubnetID:         utils.String(subnet_id),
		ServiceCidr:      utils.String(service_cidr),
		DNSServiceIP:     utils.String(dns_service_ip),
		DockerBridgeCidr: utils.String(docker_bridge_cidr)}
}

func expandAksProperties(aks_cluster *containerservice.ManagedCluster, node_pool *containerservice.AgentPool,
	ssl_config *machinelearningservices.SslConfiguration, cluster_purpose machinelearningservices.ClusterPurpose) *machinelearningservices.AKSProperties {
	// https://github.com/Azure/azure-sdk-for-go/blob/v53.1.0/services/containerservice/mgmt/2020-12-01/containerservice/models.go#L1865
	fqdn := *(aks_cluster.ManagedClusterProperties.Fqdn)
	agent_count := *(node_pool.ManagedClusterAgentPoolProfileProperties).Count
	agent_vmsize := string(node_pool.ManagedClusterAgentPoolProfileProperties.VMSize)

	if agent_count < min_number_of_nodes_prod && cluster_purpose == "FastProd" {
		min_number_of_cores_needed := int(math.Ceil(float64(min_number_of_nodes_prod) / float64(agent_count)))
		err := fmt.Errorf("error: you should pick a VM with at least %d cores", min_number_of_cores_needed)
		fmt.Println(err.Error())
	}

	return &machinelearningservices.AKSProperties{
		// ClusterFqdn - Cluster fully qualified domain name
		ClusterFqdn: utils.String(fqdn),
		// SystemServices - READ-ONLY; System services
		// AgentCount - Number of agents
		AgentCount: utils.Int32(agent_count),
		// AgentVMSize - Agent virtual machine size
		AgentVMSize: utils.String(agent_vmsize),
		// SslConfiguration - SSL configuration
		SslConfiguration: ssl_config,
		// AksNetworkingConfiguration - AKS networking configuration for vnet
		AksNetworkingConfiguration: expandAksNetworkingConfiguration(aks_cluster, node_pool),
		// ClusterPurpose - Possible values include: 'FastProd', 'DenseProd', 'DevTest'
		ClusterPurpose: cluster_purpose,
	}
}

func expandAksComputeProperties(aks_properties *machinelearningservices.AKSProperties, aks_cluster *containerservice.ManagedCluster,
	node_pool *containerservice.AgentPool, location string, description string) machinelearningservices.AKS {

	return machinelearningservices.AKS{
		// Properties - AKS properties
		Properties: aks_properties,
		// ComputeLocation - Location for the underlying compute
		ComputeLocation: &location,
		// ProvisioningState - READ-ONLY; The provision state of the cluster. Valid values are Unknown, Updating, Provisioning, Succeeded, and Failed. Possible values include: 'ProvisioningStateUnknown', 'ProvisioningStateUpdating', 'ProvisioningStateCreating', 'ProvisioningStateDeleting', 'ProvisioningStateSucceeded', 'ProvisioningStateFailed', 'ProvisioningStateCanceled'
		// Description - The description of the Machine Learning compute.
		Description: &description,
		// CreatedOn - READ-ONLY; The date and time when the compute was created.
		// ModifiedOn - READ-ONLY; The date and time when the compute was last modified.
		// ResourceID - ARM resource id of the underlying compute
		ResourceID: aks_cluster.ID,
		// ProvisioningErrors - READ-ONLY; Errors during provisioning
		// IsAttachedCompute - READ-ONLY; Indicating whether the compute was provisioned by user and brought from outside if true, or machine learning service provisioned it if false.
		// ComputeType - Possible values include: 'ComputeTypeCompute', 'ComputeTypeAKS1', 'ComputeTypeAmlCompute1', 'ComputeTypeVirtualMachine1', 'ComputeTypeHDInsight1', 'ComputeTypeDataFactory1', 'ComputeTypeDatabricks1', 'ComputeTypeDataLakeAnalytics1'
		ComputeType: "ComputeTypeAKS1",
	}
}

func expandMachineLearningComputeClusterIdentity(input []interface{}) *machinelearningservices.Identity {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	identityType := machinelearningservices.ResourceIdentityType(v["type"].(string))

	identity := machinelearningservices.Identity{
		Type: identityType,
	}

	return &identity
}

func flattenAksInferenceClusterIdentity(identity *machinelearningservices.Identity) []interface{} {
	if identity == nil {
		return []interface{}{}
	}

	t := string(identity.Type)

	principalID := ""
	if identity.PrincipalID != nil {
		principalID = *identity.PrincipalID
	}

	tenantID := ""
	if identity.TenantID != nil {
		tenantID = *identity.TenantID
	}

	return []interface{}{
		map[string]interface{}{
			"type":         t,
			"principal_id": principalID,
			"tenant_id":    tenantID,
		},
	}
}
