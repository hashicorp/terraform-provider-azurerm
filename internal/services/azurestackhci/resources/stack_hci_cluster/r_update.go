// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_cluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofilehciassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmStackHCIClusterUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AzureStackHCI.Clusters
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	cluster := clusters.ClusterPatch{}

	if d.HasChange("tags") {
		cluster.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if d.HasChange("identity") {
		cluster.Identity = expandSystemAssigned(d.Get("identity").([]any))
	}

	if _, err := client.Update(ctx, *id, cluster); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if d.HasChange("automanage_configuration_id") {
		hciAssignmentClient := meta.(*clients.Client).Automanage.ConfigurationProfileHCIAssignmentsClient
		configurationProfilesClient := meta.(*clients.Client).Automanage.ConfigurationProfilesClient
		hciAssignmentId := configurationprofilehciassignments.NewConfigurationProfileAssignmentID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, "default")

		if v, ok := d.GetOk("automanage_configuration_id"); ok {
			configurationProfileId, err := configurationprofiles.ParseConfigurationProfileID(v.(string))
			if err != nil {
				return err
			}

			if _, err = configurationProfilesClient.Get(ctx, *configurationProfileId); err != nil {
				return fmt.Errorf("checking for existing %s: %+v", configurationProfileId, err)
			}

			properties := configurationprofilehciassignments.ConfigurationProfileAssignment{
				Properties: &configurationprofilehciassignments.ConfigurationProfileAssignmentProperties{
					ConfigurationProfile: pointer.To(configurationProfileId.ID()),
				},
			}

			if _, err := hciAssignmentClient.CreateOrUpdate(ctx, hciAssignmentId, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", hciAssignmentId, err)
			}
		} else {
			assignmentResp, err := hciAssignmentClient.Get(ctx, hciAssignmentId)
			if err != nil && !response.WasNotFound(assignmentResp.HttpResponse) {
				return err
			}

			if !response.WasNotFound(assignmentResp.HttpResponse) {
				if _, err := hciAssignmentClient.Delete(ctx, hciAssignmentId); err != nil {
					return fmt.Errorf("deleting %s: %+v", id, err)
				}
			}
		}
	}

	return resourceArmStackHCIClusterRead(d, meta)
}
