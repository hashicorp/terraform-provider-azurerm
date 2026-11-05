// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package role_management_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.RoleManagementPoliciesClient

			id, err := parse.RoleManagementPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			policyId, err := FindRoleManagementPolicyId(ctx, metadata.Client.Authorization.RoleManagementPoliciesClient, id.Scope, id.RoleDefinitionId)
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *policyId)
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", policyId, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			model, err := buildRoleManagementPolicyForUpdate(pointer.To(metadata), existing.Model)
			if err != nil {
				return fmt.Errorf("could not build update request, %+v", err)
			}

			if _, err = client.Update(ctx, *policyId, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", policyId, err)
			}

			return nil
		},
	}
}
