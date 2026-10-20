// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_compute_cluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/machinelearningcomputes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceComputeClusterRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MachineLearning.MachineLearningComputes
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := machinelearningcomputes.ParseComputeID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Compute Cluster ID `%q`: %+v", d.Id(), err)
	}

	computeResource, err := client.ComputeGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(computeResource.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ComputeName)

	workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)
	d.Set("machine_learning_workspace_id", workspaceId.ID())

	// use ComputeResource to get to AKS Cluster ID and other properties
	computeCluster := computeResource.Model.Properties.(machinelearningcomputes.AmlCompute)

	if computeCluster.DisableLocalAuth != nil {
		d.Set("local_auth_enabled", !*computeCluster.DisableLocalAuth)
	}
	d.Set("description", computeCluster.Description)
	if props := computeCluster.Properties; props != nil {
		d.Set("vm_size", props.VMSize)
		d.Set("vm_priority", pointer.FromEnum(props.VMPriority))
		d.Set("scale_settings", flattenScaleSettings(props.ScaleSettings))
		d.Set("ssh", flattenUserAccountCredentials(props.UserAccountCredentials))
		enableNodePublicIP := true
		if props.EnableNodePublicIP != nil {
			enableNodePublicIP = *props.EnableNodePublicIP
		}
		d.Set("node_public_ip_enabled", enableNodePublicIP)
		if props.Subnet != nil {
			d.Set("subnet_resource_id", props.Subnet.Id)
		}

		switch *props.RemoteLoginPortPublicAccess {
		case machinelearningcomputes.RemoteLoginPortPublicAccessNotSpecified:
			d.Set("ssh_public_access_enabled", nil)
		case machinelearningcomputes.RemoteLoginPortPublicAccessEnabled:
			d.Set("ssh_public_access_enabled", true)
		case machinelearningcomputes.RemoteLoginPortPublicAccessDisabled:
			d.Set("ssh_public_access_enabled", false)
		}
	}

	if loc := computeResource.Model.Location; loc != nil {
		d.Set("location", location.Normalize(*loc))
	}

	identity, err := helpers.FlattenIdentity(computeResource.Model.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	return tags.FlattenAndSet(d, computeResource.Model.Tags)
}

func flattenScaleSettings(scaleSettings *machinelearningcomputes.ScaleSettings) []any {
	if scaleSettings == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"max_node_count":                       scaleSettings.MaxNodeCount,
			"min_node_count":                       scaleSettings.MinNodeCount,
			"scale_down_nodes_after_idle_duration": scaleSettings.NodeIdleTimeBeforeScaleDown,
		},
	}
}

func flattenUserAccountCredentials(credentials *machinelearningcomputes.UserAccountCredentials) any {
	if credentials == nil {
		return []any{}
	}
	var username string
	if credentials.AdminUserName != "" {
		username = credentials.AdminUserName
	}

	return []any{
		map[string]any{
			"admin_username": username,
			"admin_password": pointer.From(credentials.AdminUserPassword),
			"key_value":      pointer.From(credentials.AdminUserSshPublicKey),
		},
	}
}
