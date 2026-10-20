// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_cluster_principal_assignment

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/clusterprincipalassignments"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoClusterPrincipalAssignmentCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.ClusterPrincipalAssignmentsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := clusterprincipalassignments.NewPrincipalAssignmentID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cluster_name").(string), d.Get("name").(string))
	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		principalAssignment, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(principalAssignment.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(principalAssignment.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_kusto_cluster_principal_assignment", id.ID())
		}
	}

	tenantID := d.Get("tenant_id").(string)
	principalID := d.Get("principal_id").(string)
	principalType := d.Get("principal_type").(string)
	role := d.Get("role").(string)

	principalAssignment := clusterprincipalassignments.ClusterPrincipalAssignment{
		Properties: &clusterprincipalassignments.ClusterPrincipalProperties{
			TenantId:      pointer.To(tenantID),
			PrincipalId:   principalID,
			PrincipalType: clusterprincipalassignments.PrincipalType(principalType),
			Role:          clusterprincipalassignments.ClusterPrincipalRole(role),
		},
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, principalAssignment, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceKustoClusterPrincipalAssignmentRead(d, meta)
}
