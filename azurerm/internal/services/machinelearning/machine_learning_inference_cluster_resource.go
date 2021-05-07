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
			},

			"machine_learning_workspace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"kubernetes_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validate.KubernetesClusterID,
			},

			"cluster_purpose": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Dev",
			},

			"node_pool_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
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

			"ssl": {
				Type:     schema.TypeList,
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

			"sku_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceAksInferenceClusterCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	mlWorkspacesClient := meta.(*clients.Client).MachineLearning.WorkspacesClient
	mlComputeClient := meta.(*clients.Client).MachineLearning.MachineLearningComputeClient
	aksClient := meta.(*clients.Client).Containers.KubernetesClustersClient
	poolsClient := meta.(*clients.Client).Containers.AgentPoolsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// Define Inference Cluster Name
	name := d.Get("name").(string)

	// Get Machine Learning Workspace Name and Resource Group from ID
	ml_workspace_id := d.Get("machine_learning_workspace_id").(string)
	ws_id, err := parse.WorkspaceID(ml_workspace_id)
	if err != nil {
		return err
	}

	workspace_name := ws_id.Name
	resource_group_name := ws_id.ResourceGroup

	// Check if Inference Cluster already exists
	existing, err := mlComputeClient.Get(ctx, resource_group_name, workspace_name, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("error checking for existing Inference Cluster %q in Workspace %q (Resource Group %q): %s", name, workspace_name, resource_group_name, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_machine_learning_inference_cluster", *existing.ID)
	}

	// Get SKU from Workspace
	aml_ws, err := mlWorkspacesClient.Get(ctx, resource_group_name, workspace_name)
	if err != nil {
		return err
	}
	sku := aml_ws.Sku

	// Get Kubernetes Cluster Name and Resource Group from ID
	aks_id := d.Get("kubernetes_cluster_id").(string)
	aks_id_details, err := parse.KubernetesClusterID(aks_id)
	if err != nil {
		return err
	}
	kubernetes_cluster_name := aks_id_details.ManagedClusterName
	kubernetes_cluster_rg := aks_id_details.ResourceGroup

	// Get Existing AKS
	aks_cluster, err := aksClient.Get(ctx, kubernetes_cluster_rg, kubernetes_cluster_name)
	if err != nil {
		return err
	}

	pool_name := d.Get("node_pool_name").(string)
	node_pool, err := poolsClient.Get(ctx, kubernetes_cluster_rg, kubernetes_cluster_name, pool_name)
	if err != nil {
		return err
	}

	t := d.Get("tags").(map[string]interface{})

	identity := d.Get("identity").([]interface{})
	ssl_interface := d.Get("ssl").([]interface{})
	ssl := expandSSLConfig(ssl_interface)

	cluster_purpose := d.Get("cluster_purpose").(string)
	var map_cluster_purpose = map[string]string{
		"Dev":  "DevTest",
		"Test": "DevTest",
		"Prod": "FastProd",
	}
	aks_cluster_purpose := machinelearningservices.ClusterPurpose(map_cluster_purpose[cluster_purpose])

	aks_properties := expandAksProperties(&aks_cluster, &node_pool, ssl, aks_cluster_purpose)
	location := azure.NormalizeLocation(d.Get("location").(string))
	description := d.Get("description").(string)
	aks_compute_properties := expandAksComputeProperties(aks_properties, &aks_cluster, location, description)
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
		Identity: expandAksInferenceClusterIdentity(identity),
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

	future, err := mlComputeClient.CreateOrUpdate(ctx, resource_group_name, workspace_name, name, inference_cluster_parameters)
	if err != nil {
		return fmt.Errorf("error creating Inference Cluster %q in workspace %q (Resource Group %q): %+v", name, workspace_name, resource_group_name, err)
	}
	if err := future.WaitForCompletionRef(ctx, mlComputeClient.Client); err != nil {
		return fmt.Errorf("error waiting for creation of Inference Cluster %q in workspace %q (Resource Group %q): %+v", name, workspace_name, resource_group_name, err)
	}
	resp, err := mlComputeClient.Get(ctx, resource_group_name, workspace_name, name)
	if err != nil {
		return fmt.Errorf("error retrieving Inference Cluster Compute %q in workspace %q (Resource Group %q): %+v", name, workspace_name, resource_group_name, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("cannot read Inference Cluster ID %q in workspace %q (Resource Group %q) ID", name, workspace_name, resource_group_name)
	}

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := parse.NewInferenceClusterID(subscriptionId, resource_group_name, workspace_name, name)
	d.SetId(id.ID())

	return resourceAksInferenceClusterRead(d, meta)
}

func resourceAksInferenceClusterRead(d *schema.ResourceData, meta interface{}) error {
	mlWorkspacesClient := meta.(*clients.Client).MachineLearning.WorkspacesClient
	mlComputeClient := meta.(*clients.Client).MachineLearning.MachineLearningComputeClient
	aksClient := meta.(*clients.Client).Containers.KubernetesClustersClient
	poolsClient := meta.(*clients.Client).Containers.AgentPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.InferenceClusterID(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", id.ComputeName)

	// Check that Inference Cluster Response can be read
	resp, err := mlComputeClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.ComputeName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error making Read request on Inference Cluster %q in Workspace %q (Resource Group %q): %+v",
			id.ComputeName, id.WorkspaceName, id.ResourceGroup, err)
	}

	// Retrieve Machine Learning Workspace ID
	ws_resp, err := mlWorkspacesClient.Get(ctx, id.ResourceGroup, id.WorkspaceName)
	if err != nil {
		return err
	}
	d.Set("machine_learning_workspace_id", ws_resp.ID)

	// Retrieve AKS Cluster ID
	aks_resp, err := aksClient.ListByResourceGroup(ctx, id.ResourceGroup)
	if err != nil {
		return err
	}

	aks_id := *(aks_resp.Values()[0].ID)

	// Retrieve AKS Cluster name and Node pool name from ID
	aks_id_details, err := parse.KubernetesClusterID(aks_id)
	if err != nil {
		return err
	}

	kubernetes_cluster_name := aks_id_details.ManagedClusterName
	node_pool_list, _ := poolsClient.List(ctx, id.ResourceGroup, kubernetes_cluster_name)
	pool_name := node_pool_list.Values()[0].Name

	d.Set("kubernetes_cluster_id", aks_id)
	d.Set("node_pool_name", pool_name)

	// Retrieve location
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	// Retrieve Sku
	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	// Retrieve Identity
	if err := d.Set("identity", flattenAksInferenceClusterIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("error flattening identity on Workspace %q (Resource Group %q): %+v",
			id.WorkspaceName, id.ResourceGroup, err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceAksInferenceClusterDelete(d *schema.ResourceData, meta interface{}) error {
	mlComputeClient := meta.(*clients.Client).MachineLearning.MachineLearningComputeClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := parse.InferenceClusterID(d.Id())
	if err != nil {
		return err
	}
	underlying_resource_action := machinelearningservices.Detach
	future, err := mlComputeClient.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.ComputeName, underlying_resource_action)
	if err != nil {
		return fmt.Errorf("error deleting Inference Cluster %q in workspace %q (Resource Group %q): %+v",
			id.ComputeName, id.WorkspaceName, id.ResourceGroup, err)
	}
	if err := future.WaitForCompletionRef(ctx, mlComputeClient.Client); err != nil {
		return fmt.Errorf("error waiting for deletion of Inference Cluster %q in workspace %q (Resource Group %q): %+v",
			id.ComputeName, id.WorkspaceName, id.ResourceGroup, err)
	}
	return nil
}

func expandSSLConfig(input []interface{}) *machinelearningservices.SslConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	// SSL Certificate default values
	ssl_status := "Disabled"

	if !(v["cert"] == "" && v["key"] == "" && v["cname"] == "") {
		ssl_status = "Enabled"
	}

	return &machinelearningservices.SslConfiguration{
		Status:                  machinelearningservices.Status1(ssl_status),
		Cert:                    utils.String(v["cert"].(string)),
		Key:                     utils.String(v["key"].(string)),
		Cname:                   utils.String(v["cname"].(string)),
		LeafDomainLabel:         utils.String(v["leaf_domain_label"].(string)),
		OverwriteExistingDomain: utils.Bool(v["overwrite_existing_domain"].(bool))}
}

func expandAksNetworkingConfiguration(aks *containerservice.ManagedCluster, node_pool *containerservice.AgentPool) *machinelearningservices.AksNetworkingConfiguration {
	subnet_id := *(node_pool.ManagedClusterAgentPoolProfileProperties).VnetSubnetID
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
	ssl *machinelearningservices.SslConfiguration, cluster_purpose machinelearningservices.ClusterPurpose) *machinelearningservices.AKSProperties {
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
		SslConfiguration: ssl,
		// AksNetworkingConfiguration - AKS networking configuration for vnet
		AksNetworkingConfiguration: expandAksNetworkingConfiguration(aks_cluster, node_pool),
		// ClusterPurpose - Possible values include: 'FastProd', 'DenseProd', 'DevTest'
		ClusterPurpose: cluster_purpose,
	}
}

func expandAksComputeProperties(aks_properties *machinelearningservices.AKSProperties, aks_cluster *containerservice.ManagedCluster, location string, description string) machinelearningservices.AKS {
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

func expandAksInferenceClusterIdentity(input []interface{}) *machinelearningservices.Identity {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &machinelearningservices.Identity{
		Type: machinelearningservices.ResourceIdentityType(v["type"].(string)),
	}
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
