// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_database_principal_assignment

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/databaseprincipalassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoDatabasePrincipalAssignmentDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.DatabasePrincipalAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := databaseprincipalassignments.ParseDatabasePrincipalAssignmentID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
