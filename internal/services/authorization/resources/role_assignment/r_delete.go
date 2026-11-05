// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package role_assignment

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roleassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmRoleAssignmentDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Authorization.ScopedRoleAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ScopedRoleAssignmentID(d.Id())
	if err != nil {
		return err
	}

	options := roleassignments.DefaultDeleteOperationOptions()
	if id.TenantId != "" {
		options.TenantId = pointer.To(id.TenantId)
	}

	resp, err := client.Delete(ctx, id.ScopedId, options)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return err
		}
	}

	return nil
}
