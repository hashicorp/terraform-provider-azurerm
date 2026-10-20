// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_compute_instance

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/machinelearningcomputes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceComputeInstanceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MachineLearning.MachineLearningComputes
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := machinelearningcomputes.ParseComputeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ComputeGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Machine Learning Compute %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Machine Learning %s: %+v", id, err)
	}

	workspaceId := workspaces.NewWorkspaceID(subscriptionId, id.ResourceGroupName, id.WorkspaceName)

	model := resp.Model
	if model == nil {
		return fmt.Errorf("machine learning %s: model is nil", id)
	}

	props := model.Properties.(machinelearningcomputes.ComputeInstance)

	d.Set("name", id.ComputeName)
	d.Set("machine_learning_workspace_id", workspaceId.ID())

	identity, err := helpers.FlattenIdentity(model.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}

	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if props.DisableLocalAuth != nil {
		d.Set("local_auth_enabled", !*props.DisableLocalAuth)
	}

	d.Set("description", props.Description)

	if props.Properties != nil {
		d.Set("virtual_machine_size", props.Properties.VMSize)
		d.Set("authorization_type", pointer.FromEnum(props.Properties.ComputeInstanceAuthorizationType))
		d.Set("ssh", flattenComputeSSHSetting(props.Properties.SshSettings))
		d.Set("assign_to_user", flattenComputePersonalComputeInstanceSetting(props.Properties.PersonalComputeInstanceSettings))

		// Only save the subnet ID to state if the instance's workspace is NOT using managed networking.
		// Azure does not allow instances to have a user-defined subnet ID if they are to run in a workspace with managed networking.
		if props.Properties.Subnet != nil {
			if workspace, err := meta.(*clients.Client).MachineLearning.Workspaces.Get(ctx, workspaceId); err != nil || !hasManagedNetwork(workspace.Model) {
				d.Set("subnet_resource_id", props.Properties.Subnet.Id)
			}
		}

		d.Set("node_public_ip_enabled", pointer.From(props.Properties.EnableNodePublicIP))
	}

	return tags.FlattenAndSet(d, resp.Model.Tags)
}

func flattenComputePersonalComputeInstanceSetting(settings *machinelearningcomputes.PersonalComputeInstanceSettings) any {
	if settings == nil || settings.AssignedUser == nil {
		return []any{}
	}
	return []any{
		map[string]any{
			"tenant_id": settings.AssignedUser.TenantId,
			"object_id": settings.AssignedUser.ObjectId,
		},
	}
}

func flattenComputeSSHSetting(settings *machinelearningcomputes.ComputeInstanceSshSettings) any {
	if settings == nil || strings.EqualFold(string(*settings.SshPublicAccess), string(machinelearningcomputes.SshPublicAccessDisabled)) {
		return []any{}
	}

	return []any{
		map[string]any{
			"public_key": settings.AdminPublicKey,
			"username":   settings.AdminUserName,
			"port":       settings.SshPort,
		},
	}
}

func hasManagedNetwork(workspace *workspaces.Workspace) bool {
	if workspace == nil || workspace.Properties == nil || workspace.Properties.ManagedNetwork == nil {
		return false
	}
	return slices.Contains([]workspaces.IsolationMode{workspaces.IsolationModeAllowInternetOutbound, workspaces.IsolationModeAllowOnlyApprovedOutbound}, pointer.From(workspace.Properties.ManagedNetwork.IsolationMode))
}
