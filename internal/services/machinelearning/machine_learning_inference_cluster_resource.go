package machinelearning

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-02-02-preview/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2022-05-01/machinelearningcomputes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2022-05-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAksInferenceCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAksInferenceClusterCreate,
		Read:   resourceAksInferenceClusterRead,
		Delete: resourceAksInferenceClusterDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.InferenceClusterID(id)
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"kubernetes_cluster_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.KubernetesClusterID,
				// remove in 3.0 of the provider
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"location": commonschema.Location(),

			"machine_learning_workspace_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cluster_purpose": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(machinelearningcomputes.ClusterPurposeFastProd),
				ValidateFunc: validation.StringInSlice([]string{
					string(machinelearningcomputes.ClusterPurposeDevTest),
					string(machinelearningcomputes.ClusterPurposeFastProd),
					string(machinelearningcomputes.ClusterPurposeDenseProd),
				}, false),
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptionalForceNew(),

			"ssl": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"cert": {
							Type:          pluginsdk.TypeString,
							Optional:      true,
							ForceNew:      true,
							Default:       "",
							ConflictsWith: []string{"ssl.0.leaf_domain_label", "ssl.0.overwrite_existing_domain"},
						},
						"key": {
							Type:          pluginsdk.TypeString,
							Optional:      true,
							ForceNew:      true,
							Default:       "",
							ConflictsWith: []string{"ssl.0.leaf_domain_label", "ssl.0.overwrite_existing_domain"},
						},
						"cname": {
							Type:          pluginsdk.TypeString,
							Optional:      true,
							ForceNew:      true,
							Default:       "",
							ConflictsWith: []string{"ssl.0.leaf_domain_label", "ssl.0.overwrite_existing_domain"},
						},
						"leaf_domain_label": {
							Type:          pluginsdk.TypeString,
							Optional:      true,
							ForceNew:      true,
							Default:       "",
							ConflictsWith: []string{"ssl.0.cert", "ssl.0.key", "ssl.0.cname"},
						},
						"overwrite_existing_domain": {
							Type:          pluginsdk.TypeBool,
							Optional:      true,
							ForceNew:      true,
							Default:       "",
							ConflictsWith: []string{"ssl.0.cert", "ssl.0.key", "ssl.0.cname"},
						},
					},
				},
			},

			"tags": commonschema.TagsForceNew(),
		},
	}
}

func resourceAksInferenceClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.ComputeClient
	aksClient := meta.(*clients.Client).Containers.KubernetesClustersClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// Define Inference Cluster Name
	name := d.Get("name").(string)

	// Get Machine Learning Workspace Name and Resource Group from ID
	workspaceID, err := workspaces.ParseWorkspaceID(d.Get("machine_learning_workspace_id").(string))
	if err != nil {
		return err
	}
	computeId := machinelearningcomputes.NewComputeID(workspaceID.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, d.Get("name").(string))

	// Check if Inference Cluster already exists
	existing, err := client.ComputeGet(ctx, computeId)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing Inference Cluster %q in Workspace %q (Resource Group %q): %s", name, workspaceID.WorkspaceName, workspaceID.ResourceGroupName, err)
		}
	}
	if existing.Model != nil && *existing.Model.Id != "" {
		return tf.ImportAsExistsError("azurerm_machine_learning_inference_cluster", *existing.Model.Id)
	}

	// Get AKS Compute Properties
	aksID, err := managedclusters.ParseManagedClusterID(d.Get("kubernetes_cluster_id").(string))
	if err != nil {
		return err
	}
	aks, err := aksClient.Get(ctx, *aksID)
	if err != nil {
		return err
	}

	aksModel := aks.Model
	if aksModel == nil {
		return fmt.Errorf("AKS not found")
	}

	identity, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	var compute machinelearningcomputes.Compute = expandAksComputeProperties(aksID.ID(), aksModel, d)
	inferenceClusterParameters := machinelearningcomputes.ComputeResource{
		Properties: pointer.To(compute),
		Identity:   identity,
		Location:   utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.ComputeCreateOrUpdate(ctx, computeId, inferenceClusterParameters)
	if err != nil {
		return fmt.Errorf("creating Inference Cluster %q in workspace %q (Resource Group %q): %+v", name, workspaceID.WorkspaceName, workspaceID.ResourceGroupName, err)
	}
	if err := future.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("waiting for creation of Inference Cluster %q in workspace %q (Resource Group %q): %+v", name, workspaceID.ResourceGroupName, workspaceID.ResourceGroupName, err)
	}
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := machinelearningcomputes.NewComputeID(subscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)
	d.SetId(id.ID())

	return resourceAksInferenceClusterRead(d, meta)
}

func resourceAksInferenceClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.ComputeClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := machinelearningcomputes.ParseComputeID(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", id.ComputeName)

	// Check that Inference Cluster Response can be read
	resp, err := client.ComputeGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Inference Cluster %q in Workspace %q (Resource Group %q): %+v",
			id.ComputeName, id.WorkspaceName, id.ResourceGroupName, err)
	}

	// Retrieve Machine Learning Workspace ID
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	workspaceId := machinelearningcomputes.NewWorkspaceID(subscriptionId, id.ResourceGroupName, id.WorkspaceName)
	d.Set("machine_learning_workspace_id", workspaceId.ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		identity, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}

		if props := model.Properties; props != nil {
			if aksProps, ok := (*props).(machinelearningcomputes.AKS); ok {
				aksId, err := managedclusters.ParseManagedClusterIDInsensitively(*aksProps.ResourceId)
				if err != nil {
					return err
				}
				d.Set("kubernetes_cluster_id", aksId.ID())
				clusterPurpose := ""
				if aksProps.Properties != nil {
					clusterPurpose = string(pointer.From(aksProps.Properties.ClusterPurpose))
				}
				d.Set("cluster_purpose", clusterPurpose)
				d.Set("description", aksProps.Description)
			}
		}
	}

	return nil
}

func resourceAksInferenceClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.ComputeClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := machinelearningcomputes.ParseComputeID(d.Id())
	if err != nil {
		return err
	}

	opts := machinelearningcomputes.ComputeDeleteOperationOptions{
		UnderlyingResourceAction: utils.ToPtr(machinelearningcomputes.UnderlyingResourceActionDetach),
	}
	if err := client.ComputeDeleteThenPoll(ctx, *id, opts); err != nil {
		return fmt.Errorf("deleting Inference Cluster %s: %+v", *id, err)
	}
	return nil
}

func expandAksComputeProperties(aksId string, aks *managedclusters.ManagedCluster, d *pluginsdk.ResourceData) machinelearningcomputes.AKS {
	fqdn := aks.Properties.PrivateFQDN
	if fqdn == nil {
		fqdn = aks.Properties.Fqdn
	}

	return machinelearningcomputes.AKS{
		Properties: &machinelearningcomputes.AKSSchemaProperties{
			ClusterFqdn:      utils.String(*fqdn),
			SslConfiguration: expandSSLConfig(d.Get("ssl").([]interface{})),
			ClusterPurpose:   utils.ToPtr(machinelearningcomputes.ClusterPurpose(d.Get("cluster_purpose").(string))),
		},
		ComputeLocation: utils.String(aks.Location),
		Description:     utils.String(d.Get("description").(string)),
		ResourceId:      utils.String(aksId),
	}
}

func expandSSLConfig(input []interface{}) *machinelearningcomputes.SslConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	// SSL Certificate default values
	sslStatus := "Disabled"

	if !(v["cert"].(string) == "" && v["key"].(string) == "" && v["cname"].(string) == "") {
		sslStatus = "Enabled"
	}

	if !(v["leaf_domain_label"].(string) == "") {
		sslStatus = "Auto"
		v["cname"] = ""
	}

	return &machinelearningcomputes.SslConfiguration{
		Status:                  utils.ToPtr(machinelearningcomputes.SslConfigStatus(sslStatus)),
		Cert:                    utils.String(v["cert"].(string)),
		Key:                     utils.String(v["key"].(string)),
		Cname:                   utils.String(v["cname"].(string)),
		LeafDomainLabel:         utils.String(v["leaf_domain_label"].(string)),
		OverwriteExistingDomain: utils.Bool(v["overwrite_existing_domain"].(bool)),
	}
}
