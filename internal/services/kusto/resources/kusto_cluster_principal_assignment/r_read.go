// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_cluster_principal_assignment

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/clusterprincipalassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoClusterPrincipalAssignmentRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.ClusterPrincipalAssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusterprincipalassignments.ParsePrincipalAssignmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.PrincipalAssignmentName)
	d.Set("cluster_name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if resp.Model != nil {
		props := resp.Model.Properties

		if props != nil {
			d.Set("principal_id", props.PrincipalId)

			d.Set("principal_name", pointer.From(props.PrincipalName))

			d.Set("principal_type", string(props.PrincipalType))
			d.Set("role", string(props.Role))

			d.Set("tenant_id", pointer.From(props.TenantId))

			d.Set("tenant_name", pointer.From(props.TenantName))
		}
	}

	return nil
}
