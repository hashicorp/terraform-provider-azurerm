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

const minNumberOfNodesProd int32 = 12

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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"machine_learning_workspace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"kubernetes_cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.KubernetesClusterID,
			},

			"cluster_purpose": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Dev",
			},

			"node_pool_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NodePoolID,
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
	unparsedWorkspaceID := d.Get("machine_learning_workspace_id").(string)

	workspaceID, err := parse.WorkspaceID(unparsedWorkspaceID)
	if err != nil {
		return err
	}

	// Check if Inference Cluster already exists
	existing, err := mlComputeClient.Get(ctx, workspaceID.ResourceGroup, workspaceID.Name, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("error checking for existing Inference Cluster %q in Workspace %q (Resource Group %q): %s", name, workspaceID.Name, workspaceID.ResourceGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_machine_learning_inference_cluster", *existing.ID)
	}

	// Get SKU from Workspace
	workspace, err := mlWorkspacesClient.Get(ctx, workspaceID.ResourceGroup, workspaceID.Name)
	if err != nil {
		return err
	}
	sku := workspace.Sku

	// Get Kubernetes Cluster Name and Resource Group from ID
	unparsedAksID := d.Get("kubernetes_cluster_id").(string)
	aksID, err := parse.KubernetesClusterID(unparsedAksID)
	if err != nil {
		return err
	}

	// Get Existing AKS
	aks, err := aksClient.Get(ctx, aksID.ResourceGroup, aksID.ManagedClusterName)
	if err != nil {
		return err
	}

	unparsedAgentPoolID := d.Get("node_pool_id").(string)
	agentPoolID, err := parse.NodePoolID(unparsedAgentPoolID )
	if err != nil {
		return err
	}
	nodePool, err := poolsClient.Get(ctx, aksID.ResourceGroup, aksID.ManagedClusterName, agentPoolID.AgentPoolName)
	if err != nil {
		return err
	}

	t := d.Get("tags").(map[string]interface{})

	identity := d.Get("identity").([]interface{})
	sslInterface := d.Get("ssl").([]interface{})
	ssl := expandSSLConfig(sslInterface)

	unmappedclusterPurpose := d.Get("cluster_purpose").(string)
	var mapClusterPurpose = map[string]string{
		"Dev":  "DevTest",
		"Test": "DevTest",
		"Prod": "FastProd",
	}
	clusterPurpose := machinelearningservices.ClusterPurpose(mapClusterPurpose[unmappedclusterPurpose])

	aksProperties := expandAksProperties(&aks, &nodePool, ssl, clusterPurpose)
	location := azure.NormalizeLocation(d.Get("location").(string))
	description := d.Get("description").(string)
	aksComputeProperties := expandAksComputeProperties(aksProperties, &aks, location, description)
	aksCompute, isAks := (machinelearningservices.BasicCompute).AsAKS(aksComputeProperties)
	if !isAks {
		return fmt.Errorf("error: No AKS cluster")
	}

	inferenceClusterParameters := machinelearningservices.ComputeResource{
		// Properties - Compute properties
		Properties: aksCompute,
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
		aksCompute.Description = utils.String(v.(string))
	}

	future, err := mlComputeClient.CreateOrUpdate(ctx, workspaceID.ResourceGroup, workspaceID.Name, name, inferenceClusterParameters)
	if err != nil {
		return fmt.Errorf("error creating Inference Cluster %q in workspace %q (Resource Group %q): %+v", name, workspaceID.Name, workspaceID.ResourceGroup, err)
	}
	if err := future.WaitForCompletionRef(ctx, mlComputeClient.Client); err != nil {
		return fmt.Errorf("error waiting for creation of Inference Cluster %q in workspace %q (Resource Group %q): %+v", name, workspaceID.Name, workspaceID.ResourceGroup, err)
	}
	resp, err := mlComputeClient.Get(ctx, workspaceID.ResourceGroup, workspaceID.Name, name)
	if err != nil {
		return fmt.Errorf("error retrieving Inference Cluster Compute %q in workspace %q (Resource Group %q): %+v", name, workspaceID.Name, workspaceID.ResourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("cannot read Inference Cluster ID %q in workspace %q (Resource Group %q) ID", name, workspaceID.Name, workspaceID.ResourceGroup)
	}

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := parse.NewInferenceClusterID(subscriptionId, workspaceID.ResourceGroup, workspaceID.Name, name)
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
	mlResp, err := mlWorkspacesClient.Get(ctx, id.ResourceGroup, id.WorkspaceName)
	if err != nil {
		return err
	}
	d.Set("machine_learning_workspace_id", mlResp.ID)

	// Retrieve AKS Cluster ID
	aksResp, err := aksClient.ListByResourceGroup(ctx, id.ResourceGroup)
	if err != nil {
		return err
	}

	unparsedAksId := *(aksResp.Values()[0].ID)

	// Retrieve AKS Cluster name and Node pool name from ID
	aksId, err := parse.KubernetesClusterID(unparsedAksId)
	if err != nil {
		return err
	}

	nodePoolList, _ := poolsClient.List(ctx, id.ResourceGroup, aksId.ManagedClusterName)
	unparsedPoolId := *(nodePoolList.Values()[0].ID)
	poolId, err := parse.NodePoolID(unparsedPoolId)
	if err != nil {
		return err
	}

	d.Set("kubernetes_cluster_id", aksId)
	d.Set("node_pool_name", poolId.AgentPoolName)

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

	future, err := mlComputeClient.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.ComputeName, machinelearningservices.Detach)
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
	sslStatus := "Disabled"

	if !(v["cert"] == "" && v["key"] == "" && v["cname"] == "") {
		sslStatus = "Enabled"
	}

	return &machinelearningservices.SslConfiguration{
		Status:                  machinelearningservices.Status1(sslStatus),
		Cert:                    utils.String(v["cert"].(string)),
		Key:                     utils.String(v["key"].(string)),
		Cname:                   utils.String(v["cname"].(string)),
		LeafDomainLabel:         utils.String(v["leaf_domain_label"].(string)),
		OverwriteExistingDomain: utils.Bool(v["overwrite_existing_domain"].(bool))}
}

func expandAksNetworkingConfiguration(aks *containerservice.ManagedCluster, nodePool *containerservice.AgentPool) *machinelearningservices.AksNetworkingConfiguration {
	subnet_id := *(nodePool.ManagedClusterAgentPoolProfileProperties).VnetSubnetID
	service_cidr := *(aks.NetworkProfile.ServiceCidr)
	dns_service_ip := *(aks.NetworkProfile.DNSServiceIP)
	docker_bridge_cidr := *(aks.NetworkProfile.DockerBridgeCidr)

	return &machinelearningservices.AksNetworkingConfiguration{
		SubnetID:         utils.String(subnet_id),
		ServiceCidr:      utils.String(service_cidr),
		DNSServiceIP:     utils.String(dns_service_ip),
		DockerBridgeCidr: utils.String(docker_bridge_cidr)}
}

func expandAksProperties(aks *containerservice.ManagedCluster, nodePool *containerservice.AgentPool,
	ssl *machinelearningservices.SslConfiguration, clusterPurpose machinelearningservices.ClusterPurpose) *machinelearningservices.AKSProperties {
	fqdn := *(aks.ManagedClusterProperties.Fqdn)
	agentCount := *(nodePool.ManagedClusterAgentPoolProfileProperties).Count
	agentVmSize := string(nodePool.ManagedClusterAgentPoolProfileProperties.VMSize)

	if agentCount < minNumberOfNodesProd && clusterPurpose == "FastProd" {
		minNumberOfCores := int(math.Ceil(float64(minNumberOfNodesProd) / float64(agentCount)))
		err := fmt.Errorf("error: you should pick a VM with at least %d cores", minNumberOfCores)
		fmt.Println(err.Error())
	}

	return &machinelearningservices.AKSProperties{
		// ClusterFqdn - Cluster fully qualified domain name
		ClusterFqdn: utils.String(fqdn),
		// SystemServices - READ-ONLY; System services
		// AgentCount - Number of agents
		AgentCount: utils.Int32(agentCount),
		// AgentVMSize - Agent virtual machine size
		AgentVMSize: utils.String(agentVmSize),
		// SslConfiguration - SSL configuration
		SslConfiguration: ssl,
		// AksNetworkingConfiguration - AKS networking configuration for vnet
		AksNetworkingConfiguration: expandAksNetworkingConfiguration(aks, nodePool),
		// ClusterPurpose - Possible values include: 'FastProd', 'DenseProd', 'DevTest'
		ClusterPurpose: clusterPurpose,
	}
}

func expandAksComputeProperties(aksProperties *machinelearningservices.AKSProperties, aks *containerservice.ManagedCluster, location string, description string) machinelearningservices.AKS {
	return machinelearningservices.AKS{
		// Properties - AKS properties
		Properties: aksProperties,
		// ComputeLocation - Location for the underlying compute
		ComputeLocation: &location,
		// ProvisioningState - READ-ONLY; The provision state of the cluster. Valid values are Unknown, Updating, Provisioning, Succeeded, and Failed. Possible values include: 'ProvisioningStateUnknown', 'ProvisioningStateUpdating', 'ProvisioningStateCreating', 'ProvisioningStateDeleting', 'ProvisioningStateSucceeded', 'ProvisioningStateFailed', 'ProvisioningStateCanceled'
		// Description - The description of the Machine Learning compute.
		Description: &description,
		// CreatedOn - READ-ONLY; The date and time when the compute was created.
		// ModifiedOn - READ-ONLY; The date and time when the compute was last modified.
		// ResourceID - ARM resource id of the underlying compute
		ResourceID: aks.ID,
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
