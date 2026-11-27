// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kusto_database_principal_assignment

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/databaseprincipalassignments"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	kustoClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoDatabasePrincipalAssignmentCreate(d *pluginsdk.ResourceData, meta any) error {
	client := kustoClient.From(meta).DatabasePrincipalAssignmentsClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	id := databaseprincipalassignments.NewDatabasePrincipalAssignmentID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cluster_name").(string), d.Get("database_name").(string), d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_kusto_database_principal_assignment", id.ID())
		}
	}

	principalAssignment := databaseprincipalassignments.DatabasePrincipalAssignment{
		Properties: &databaseprincipalassignments.DatabasePrincipalProperties{
			TenantId:      pointer.To(d.Get("tenant_id").(string)),
			PrincipalId:   d.Get("principal_id").(string),
			PrincipalType: databaseprincipalassignments.PrincipalType(d.Get("principal_type").(string)),
			Role:          databaseprincipalassignments.DatabasePrincipalRole(d.Get("role").(string)),
		},
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, principalAssignment, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceKustoDatabasePrincipalAssignmentRead(d, meta)
}
