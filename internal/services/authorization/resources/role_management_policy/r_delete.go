// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package role_management_policy

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if _, err := parse.RoleManagementPolicyID(metadata.ResourceData.Id()); err != nil {
				return err
			}

			// Role Management Policies cannot be deleted, so we'll just return without doing anything
			return nil
		},
	}
}
