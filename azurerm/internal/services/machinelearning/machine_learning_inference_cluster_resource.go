package machinelearning

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2021-02-01/containerservice"
	"github.com/Azure/azure-sdk-for-go/services/machinelearningservices/mgmt/2020-04-01/machinelearningservices"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/machinelearning/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/machinelearning/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"

	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

const minNumberOfNodesProd int32 = 12

func resourceAksInferenceCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAksInferenceClusterCreate,
		Read:   resourceAksInferenceClusterRead,
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

			"kubernetes_cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.KubernetesClusterID,
				// remove in 3.0 of the provider
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"location": azure.SchemaLocation(),

			"machine_learning_workspace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cluster_purpose": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "FastProd",
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
							ForceNew: true,
							Default:  "",
						},
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "",
						},
						"cname": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "",
						},
					},
				},
			},

			"tags": tags.ForceNewSchema(),
		},
	}
}

func resourceAksInferenceClusterCreate(d *schema.ResourceData, meta interface{}) error {
	mlComputeClient := meta.(*clients.Client).MachineLearning.MachineLearningComputeClient
	aksClient := meta.(*clients.Client).Containers.KubernetesClustersClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// Define Inference Cluster Name
	name := d.Get("name").(string)

	// Get Machine Learning Workspace Name and Resource Group from ID
	workspaceID, err := parse.WorkspaceID(d.Get("machine_learning_workspace_id").(string))
	if err != nil {
		return err
	}

	// Check if Inference Cluster already exists
	existing, err := mlComputeClient.Get(ctx, workspaceID.ResourceGroup, workspaceID.Name, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing Inference Cluster %q in Workspace %q (Resource Group %q): %s", name, workspaceID.Name, workspaceID.ResourceGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_machine_learning_inference_cluster", *existing.ID)
	}

	// Get Kubernetes Cluster Name and Resource Group from ID
	aksID, err := parse.KubernetesClusterID(d.Get("kubernetes_cluster_id").(string))
	if err != nil {
		return err
	}

	// Get Existing AKS
	aks, err := aksClient.Get(ctx, aksID.ResourceGroup, aksID.ManagedClusterName)
	if err != nil {
		return err
	}

	ssl := expandSSLConfig(d.Get("ssl").([]interface{}))

	clusterPurpose := machinelearningservices.ClusterPurpose(d.Get("cluster_purpose").(string))
	aksProperties := expandAksProperties(&aks, ssl, clusterPurpose)
	location := azure.NormalizeLocation(d.Get("location").(string))
	description := d.Get("description").(string)
	aksComputeProperties, isAks := (machinelearningservices.BasicCompute).AsAKS(expandAksComputeProperties(aksProperties, &aks, location, description))
	if !isAks {
		return fmt.Errorf("error: No AKS cluster")
	}

	inferenceClusterParameters := machinelearningservices.ComputeResource{
		Properties: aksComputeProperties,
		Location:   &location,
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("description"); ok {
		aksComputeProperties.Description = utils.String(v.(string))
	}

	future, err := mlComputeClient.CreateOrUpdate(ctx, workspaceID.ResourceGroup, workspaceID.Name, name, inferenceClusterParameters)
	if err != nil {
		return fmt.Errorf("creating Inference Cluster %q in workspace %q (Resource Group %q): %+v", name, workspaceID.Name, workspaceID.ResourceGroup, err)
	}
	if err := future.WaitForCompletionRef(ctx, mlComputeClient.Client); err != nil {
		return fmt.Errorf("waiting for creation of Inference Cluster %q in workspace %q (Resource Group %q): %+v", name, workspaceID.Name, workspaceID.ResourceGroup, err)
	}
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := parse.NewInferenceClusterID(subscriptionId, workspaceID.ResourceGroup, workspaceID.Name, name)
	d.SetId(id.ID())

	return resourceAksInferenceClusterRead(d, meta)
}

func resourceAksInferenceClusterRead(d *schema.ResourceData, meta interface{}) error {
	mlWorkspacesClient := meta.(*clients.Client).MachineLearning.WorkspacesClient
	mlComputeClient := meta.(*clients.Client).MachineLearning.MachineLearningComputeClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.InferenceClusterID(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", id.ComputeName)

	// Check that Inference Cluster Response can be read
	computeResource, err := mlComputeClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.ComputeName)
	if err != nil {
		if utils.ResponseWasNotFound(computeResource.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Inference Cluster %q in Workspace %q (Resource Group %q): %+v",
			id.ComputeName, id.WorkspaceName, id.ResourceGroup, err)
	}

	// Retrieve Machine Learning Workspace ID
	mlResp, err := mlWorkspacesClient.Get(ctx, id.ResourceGroup, id.WorkspaceName)
	if err != nil {
		return err
	}
	d.Set("machine_learning_workspace_id", mlResp.ID)

	// use ComputeResource to get to AKS Cluster ID and other properties
	aksComputeProperties, isAks := (machinelearningservices.BasicCompute).AsAKS(computeResource.Properties)
	if !isAks {
		return fmt.Errorf("compute resource %s is not an AKS cluster", id.ComputeName)
	}

	// Retrieve AKS Cluster ID
	aksId, err := parse.KubernetesClusterID(*aksComputeProperties.ResourceID)
	if err != nil {
		return err
	}
	d.Set("kubernetes_cluster_id", aksId.ID())
	d.Set("cluster_purpose", string(aksComputeProperties.Properties.ClusterPurpose))
	d.Set("description", aksComputeProperties.Description)
	d.Set("ssl", aksComputeProperties.Properties.SslConfiguration)

	// Retrieve location
	if location := computeResource.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	return tags.FlattenAndSet(d, computeResource.Tags)
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
		return fmt.Errorf("deleting Inference Cluster %q in workspace %q (Resource Group %q): %+v",
			id.ComputeName, id.WorkspaceName, id.ResourceGroup, err)
	}
	if err := future.WaitForCompletionRef(ctx, mlComputeClient.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Inference Cluster %q in workspace %q (Resource Group %q): %+v",
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

	if !(v["cert"].(string) == "" && v["key"].(string) == "" && v["cname"].(string) == "") {
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

func expandAksProperties(aks *containerservice.ManagedCluster,
	ssl *machinelearningservices.SslConfiguration, clusterPurpose machinelearningservices.ClusterPurpose) *machinelearningservices.AKSProperties {
	fqdn := *(aks.ManagedClusterProperties.Fqdn)

	return &machinelearningservices.AKSProperties{
		ClusterFqdn:      utils.String(fqdn),
		SslConfiguration: ssl,
		ClusterPurpose:   clusterPurpose,
	}
}

func expandAksComputeProperties(aksProperties *machinelearningservices.AKSProperties, aks *containerservice.ManagedCluster, location string, description string) machinelearningservices.AKS {
	return machinelearningservices.AKS{
		Properties:      aksProperties,
		ComputeLocation: &location,
		Description:     &description,
		ResourceID:      aks.ID,
		ComputeType:     "ComputeTypeAKS1",
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
