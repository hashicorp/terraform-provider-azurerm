// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stack_hci_cluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofilehciassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	automanageClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/client"
	azurestackhciClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmStackHCIClusterDelete(d *pluginsdk.ResourceData, meta any) error {
	client := azurestackhciClient.From(meta).Clusters
	hciAssignmentClient := automanageClient.From(meta).ConfigurationProfileHCIAssignmentsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	hciAssignmentId := configurationprofilehciassignments.NewConfigurationProfileAssignmentID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, "default")
	assignmentResp, err := hciAssignmentClient.Get(ctx, hciAssignmentId)
	if err != nil && !response.WasNotFound(assignmentResp.HttpResponse) {
		return err
	}

	if !response.WasNotFound(assignmentResp.HttpResponse) {
		if _, err := hciAssignmentClient.Delete(ctx, hciAssignmentId); err != nil {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
