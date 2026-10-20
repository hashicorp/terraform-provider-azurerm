// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_inference_cluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/machinelearningcomputes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAksInferenceClusterCreate(d *pluginsdk.ResourceData, meta any) error {
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
	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.ComputeGet(ctx, computeId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing Inference Cluster %q in Workspace %q (Resource Group %q): %s", name, workspaceID.WorkspaceName, workspaceID.ResourceGroupName, err)
			}
		}
		if existing.Model != nil && *existing.Model.Id != "" {
			return tf.ImportAsExistsError("azurerm_machine_learning_inference_cluster", *existing.Model.Id)
		}
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

	identity, err := helpers.ExpandIdentity(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	inferenceClusterParameters := machinelearningcomputes.ComputeResource{
		Properties: expandAksComputeProperties(aksID.ID(), aksModel, d),
		Identity:   identity,
		Location:   pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:       tags.Expand(d.Get("tags").(map[string]any)),
	}

	id := machinelearningcomputes.NewComputeID(meta.(*clients.Client).Account.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)

	if err := client.ComputeCreateOrUpdateCallbackThenPoll(ctx, computeId, inferenceClusterParameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating Inference Cluster %q in workspace %q (Resource Group %q): %+v", name, workspaceID.WorkspaceName, workspaceID.ResourceGroupName, err)
	}
	d.SetId(id.ID())

	return resourceAksInferenceClusterRead(d, meta)
}

func expandAksComputeProperties(aksId string, aks *managedclusters.ManagedCluster, d *pluginsdk.ResourceData) machinelearningcomputes.AKS {
	fqdn := aks.Properties.PrivateFQDN
	if fqdn == nil {
		fqdn = aks.Properties.Fqdn
	}

	return machinelearningcomputes.AKS{
		Properties: &machinelearningcomputes.AKSSchemaProperties{
			ClusterFqdn:      pointer.To(*fqdn),
			SslConfiguration: expandSSLConfig(d.Get("ssl").([]any)),
			ClusterPurpose:   pointer.ToEnum[machinelearningcomputes.ClusterPurpose](d.Get("cluster_purpose").(string)),
		},
		ComputeLocation: pointer.To(aks.Location),
		Description:     pointer.To(d.Get("description").(string)),
		ResourceId:      pointer.To(aksId),
	}
}

func expandSSLConfig(input []any) *machinelearningcomputes.SslConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	// SSL Certificate default values
	sslStatus := "Disabled"

	if v["cert"].(string) != "" || v["key"].(string) != "" || v["cname"].(string) != "" {
		sslStatus = "Enabled"
	}

	if v["leaf_domain_label"].(string) != "" {
		sslStatus = "Auto"
		v["cname"] = ""
	}

	return &machinelearningcomputes.SslConfiguration{
		Status:                  pointer.ToEnum[machinelearningcomputes.SslConfigStatus](sslStatus),
		Cert:                    pointer.To(v["cert"].(string)),
		Key:                     pointer.To(v["key"].(string)),
		Cname:                   pointer.To(v["cname"].(string)),
		LeafDomainLabel:         pointer.To(v["leaf_domain_label"].(string)),
		OverwriteExistingDomain: pointer.To(v["overwrite_existing_domain"].(bool)),
	}
}
