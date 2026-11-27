// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kusto_cluster_principal_assignment

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/clusterprincipalassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	kustoClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoClusterPrincipalAssignmentDelete(d *pluginsdk.ResourceData, meta any) error {
	client := kustoClient.From(meta).ClusterPrincipalAssignmentsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := clusterprincipalassignments.ParsePrincipalAssignmentID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
