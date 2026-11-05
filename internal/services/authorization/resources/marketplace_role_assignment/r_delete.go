// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package marketplace_role_assignment

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roleassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.ScopedRoleAssignmentsClient

			id, err := parse.ScopedRoleAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			options := roleassignments.DefaultDeleteOperationOptions()
			if id.TenantId != "" {
				options.TenantId = &id.TenantId
			}

			resp, err := client.Delete(ctx, id.ScopedId, options)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return err
				}
			}

			return nil
		},
	}
}
