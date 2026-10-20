// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_compute_instance

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
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

func resourceComputeInstanceCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MachineLearning.MachineLearningComputes
	mlWorkspacesClient := meta.(*clients.Client).MachineLearning.Workspaces
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceID, _ := workspaces.ParseWorkspaceID(d.Get("machine_learning_workspace_id").(string))
	id := machinelearningcomputes.NewComputeID(subscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.ComputeGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing Machine Learning Compute (%q): %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_machine_learning_compute_instance", id.ID())
		}
	}

	identity, err := helpers.ExpandIdentity(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	var subnet *machinelearningcomputes.ResourceId
	if subnetId, ok := d.GetOk("subnet_resource_id"); ok {
		subnet = &machinelearningcomputes.ResourceId{
			Id: subnetId.(string),
		}
	}

	// NOTE: The 'ComputeResource' struct contains the information
	// which is related to the parent resource of the instance that is
	// to be deployed (e.g., the workspace), which is why we need to
	// GET the workspace to discover the location it has been deployed to.
	// If we do not set the correct location, the identity will be created
	// and then orphaned in the incorrect region.
	workspace, err := mlWorkspacesClient.Get(ctx, *workspaceID)
	if err != nil {
		return err
	}

	model := workspace.Model
	if model == nil {
		return fmt.Errorf("machine learning %s Workspace: model is nil", id)
	}

	if model.Location == nil {
		return fmt.Errorf("machine learning %s Workspace: model `Location` is nil", id)
	}

	if model.Properties == nil {
		return fmt.Errorf("machine learning %s Workspace: model `Properties` is nil", id)
	}
	if !d.Get("node_public_ip_enabled").(bool) && d.Get("subnet_resource_id").(string) == "" && !hasManagedNetwork(model) {
		return fmt.Errorf("`subnet_resource_id` must be set if `node_public_ip_enabled` is set to `false` and the workspace is NOT using a managed network (isolation_mode = \"Disabled\")")
	}

	parameters := machinelearningcomputes.ComputeResource{
		Identity: identity,
		Location: pointer.To(location.Normalize(*model.Location)),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}

	props := machinelearningcomputes.ComputeInstance{
		Properties: &machinelearningcomputes.ComputeInstanceProperties{
			VMSize:                          pointer.To(d.Get("virtual_machine_size").(string)),
			Subnet:                          subnet,
			SshSettings:                     expandComputeSSHSetting(d.Get("ssh").([]any)),
			PersonalComputeInstanceSettings: expandComputePersonalComputeInstanceSetting(d.Get("assign_to_user").([]any)),
			EnableNodePublicIP:              pointer.To(d.Get("node_public_ip_enabled").(bool)),
		},
		Description:      pointer.To(d.Get("description").(string)),
		DisableLocalAuth: pointer.To(!d.Get("local_auth_enabled").(bool)),
	}

	// NOTE: The 'location' field is not supported for instances, "Compute clusters can be created in
	// a different region than your workspace. This functionality is only available for compute
	// clusters, not compute instances"
	//
	// https://learn.microsoft.com/azure/machine-learning/how-to-create-attach-compute-cluster?view=azureml-api-2&tabs=python#limitations

	if v, ok := d.GetOk("authorization_type"); ok {
		props.Properties.ComputeInstanceAuthorizationType = pointer.ToEnum[machinelearningcomputes.ComputeInstanceAuthorizationType](v.(string))
	}

	parameters.Properties = props

	if err := client.ComputeCreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating Machine Learning %s: %+v", id, err)
	}
	d.SetId(id.ID())

	return resourceComputeInstanceRead(d, meta)
}

func expandComputePersonalComputeInstanceSetting(input []any) *machinelearningcomputes.PersonalComputeInstanceSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	value := input[0].(map[string]any)
	return &machinelearningcomputes.PersonalComputeInstanceSettings{
		AssignedUser: &machinelearningcomputes.AssignedUser{
			ObjectId: value["object_id"].(string),
			TenantId: value["tenant_id"].(string),
		},
	}
}

func expandComputeSSHSetting(input []any) *machinelearningcomputes.ComputeInstanceSshSettings {
	if len(input) == 0 {
		return &machinelearningcomputes.ComputeInstanceSshSettings{
			SshPublicAccess: pointer.To(machinelearningcomputes.SshPublicAccessDisabled),
		}
	}
	value := input[0].(map[string]any)
	return &machinelearningcomputes.ComputeInstanceSshSettings{
		SshPublicAccess: pointer.To(machinelearningcomputes.SshPublicAccessEnabled),
		AdminPublicKey:  pointer.To(value["public_key"].(string)),
	}
}
