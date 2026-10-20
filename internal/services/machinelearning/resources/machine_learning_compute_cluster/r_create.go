// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_compute_cluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/machinelearningcomputes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceComputeClusterCreate(d *pluginsdk.ResourceData, meta any) error {
	mlWorkspacesClient := meta.(*clients.Client).MachineLearning.Workspaces
	client := meta.(*clients.Client).MachineLearning.MachineLearningComputes
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceID, err := workspaces.ParseWorkspaceID(d.Get("machine_learning_workspace_id").(string))
	if err != nil {
		return err
	}
	// Get the Machine Learning Workspace...
	id := machinelearningcomputes.NewComputeID(workspaceID.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, d.Get("name").(string))

	workspace, err := mlWorkspacesClient.Get(ctx, *workspaceID)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", workspaceID, err)
	}

	workspaceModel := workspace.Model
	if workspaceModel == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", workspaceID)
	}

	if workspaceModel.Sku == nil || workspaceModel.Sku.Tier == nil || workspaceModel.Sku.Name == "" {
		return fmt.Errorf("retrieving %s: `sku` was nil or empty", workspaceID)
	}

	if workspaceModel.Location == nil {
		return fmt.Errorf("retrieving %s: `location` was nil", workspaceID)
	}

	identity, err := helpers.ExpandIdentity(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.ComputeGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_machine_learning_compute_cluster", id.ID())
		}
	}

	vmPriority := machinelearningcomputes.VMPriority(d.Get("vm_priority").(string))
	computeClusterAmlComputeProperties := machinelearningcomputes.AmlComputeProperties{
		VMSize:                 pointer.To(d.Get("vm_size").(string)),
		VMPriority:             &vmPriority,
		ScaleSettings:          expandScaleSettings(d.Get("scale_settings").([]any)),
		UserAccountCredentials: expandUserAccountCredentials(d.Get("ssh").([]any)),
		EnableNodePublicIP:     pointer.To(d.Get("node_public_ip_enabled").(bool)),
	}

	computeClusterAmlComputeProperties.RemoteLoginPortPublicAccess = pointer.To(machinelearningcomputes.RemoteLoginPortPublicAccessDisabled)
	if d.Get("ssh_public_access_enabled").(bool) {
		computeClusterAmlComputeProperties.RemoteLoginPortPublicAccess = pointer.To(machinelearningcomputes.RemoteLoginPortPublicAccessEnabled)
	}

	if subnetId, ok := d.GetOk("subnet_resource_id"); ok && subnetId.(string) != "" {
		computeClusterAmlComputeProperties.Subnet = &machinelearningcomputes.ResourceId{Id: subnetId.(string)}
	}

	// NOTE: The 'AmlCompute' 'ComputeLocation' field should always point
	// to configuration files 'location' field...
	computeClusterProperties := machinelearningcomputes.AmlCompute{
		Properties:       &computeClusterAmlComputeProperties,
		ComputeLocation:  pointer.To(d.Get("location").(string)),
		Description:      pointer.To(d.Get("description").(string)),
		DisableLocalAuth: pointer.To(!d.Get("local_auth_enabled").(bool)),
	}

	// NOTE: The 'ComputeResource' 'Location' field should always point
	// to the workspace's 'location'...
	computeClusterParameters := machinelearningcomputes.ComputeResource{
		Properties: computeClusterProperties,
		Identity:   identity,
		Location:   workspaceModel.Location,
		Tags:       tags.Expand(d.Get("tags").(map[string]any)),
		Sku: &machinelearningcomputes.Sku{
			Name: workspaceModel.Sku.Name,
			Tier: pointer.ToEnum[machinelearningcomputes.SkuTier](string(*workspaceModel.Sku.Tier)),
		},
	}

	if err := client.ComputeCreateOrUpdateCallbackThenPoll(ctx, id, computeClusterParameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	return resourceComputeClusterRead(d, meta)
}

func expandScaleSettings(input []any) *machinelearningcomputes.ScaleSettings {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	maxNodeCount := int64(v["max_node_count"].(int))

	return &machinelearningcomputes.ScaleSettings{
		MaxNodeCount:                maxNodeCount,
		MinNodeCount:                pointer.To(int64(v["min_node_count"].(int))),
		NodeIdleTimeBeforeScaleDown: pointer.To(v["scale_down_nodes_after_idle_duration"].(string)),
	}
}

func expandUserAccountCredentials(input []any) *machinelearningcomputes.UserAccountCredentials {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]any)

	return &machinelearningcomputes.UserAccountCredentials{
		AdminUserName:         v["admin_username"].(string),
		AdminUserPassword:     pointer.To(v["admin_password"].(string)),
		AdminUserSshPublicKey: pointer.To(v["key_value"].(string)),
	}
}
