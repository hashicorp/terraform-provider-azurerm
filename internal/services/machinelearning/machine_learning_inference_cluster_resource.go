// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-05-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/machinelearningcomputes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
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
			_, err := machinelearningcomputes.ParseComputeID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
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
				ValidateFunc: commonids.ValidateKubernetesClusterID,
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
	client := meta.(*clients.Client).MachineLearning.MachineLearningComputes
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
	aksID, err := commonids.ParseKubernetesClusterID(d.Get("kubernetes_cluster_id").(string))
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

	identity, err := expandIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	inferenceClusterParameters := machinelearningcomputes.ComputeResource{
		Properties: expandAksComputeProperties(aksID.ID(), aksModel, d),
		Identity:   identity,
		Location:   utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.ComputeCreateOrUpdate(ctx, computeId, inferenceClusterParameters)
	if err != nil {
		return fmt.Errorf("creating Inference Cluster %q in workspace %q (Resource Group %q): %+v", name, workspaceID.WorkspaceName, workspaceID.ResourceGroupName, err)
	}
	if err := future.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for creation of Inference Cluster %q in workspace %q (Resource Group %q): %+v", name, workspaceID.ResourceGroupName, workspaceID.ResourceGroupName, err)
	}
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := machinelearningcomputes.NewComputeID(subscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)
	d.SetId(id.ID())

	return resourceAksInferenceClusterRead(d, meta)
}

func resourceAksInferenceClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.MachineLearningComputes
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := machinelearningcomputes.ParseComputeID(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", id.ComputeName)

	// Check that Inference Cluster Response can be read
	computeResource, err := client.ComputeGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(computeResource.HttpResponse) {
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

	// use ComputeResource to get to AKS Cluster ID and other properties
	aksComputeProperties := computeResource.Model.Properties.(machinelearningcomputes.AKS)

	// Retrieve AKS Cluster ID
	aksId, err := commonids.ParseKubernetesClusterIDInsensitively(*aksComputeProperties.ResourceId)
	if err != nil {
		return err
	}
	d.Set("kubernetes_cluster_id", aksId.ID())
	clusterPurpose := ""
	if aksComputeProperties.Properties != nil {
		clusterPurpose = string(pointer.From(aksComputeProperties.Properties.ClusterPurpose))
	}
	d.Set("cluster_purpose", clusterPurpose)
	d.Set("description", aksComputeProperties.Description)

	// Retrieve location
	if location := computeResource.Model.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	identity, err := flattenIdentity(computeResource.Model.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	return tags.FlattenAndSet(d, computeResource.Model.Tags)
}

func resourceAksInferenceClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.MachineLearningComputes
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := machinelearningcomputes.ParseComputeID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.ComputeDelete(ctx, *id, machinelearningcomputes.ComputeDeleteOperationOptions{
		UnderlyingResourceAction: pointer.To(machinelearningcomputes.UnderlyingResourceActionDetach),
	})
	if err != nil {
		return fmt.Errorf("deleting Inference Cluster %q in workspace %q (Resource Group %q): %+v",
			id.ComputeName, id.WorkspaceName, id.ResourceGroupName, err)
	}
	if err := future.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for deletion of Inference Cluster %q in workspace %q (Resource Group %q): %+v",
			id.ComputeName, id.WorkspaceName, id.ResourceGroupName, err)
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
			ClusterPurpose:   pointer.To(machinelearningcomputes.ClusterPurpose(d.Get("cluster_purpose").(string))),
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
		Status:                  pointer.To(machinelearningcomputes.SslConfigStatus(sslStatus)),
		Cert:                    utils.String(v["cert"].(string)),
		Key:                     utils.String(v["key"].(string)),
		Cname:                   utils.String(v["cname"].(string)),
		LeafDomainLabel:         utils.String(v["leaf_domain_label"].(string)),
		OverwriteExistingDomain: utils.Bool(v["overwrite_existing_domain"].(bool)),
	}
}
