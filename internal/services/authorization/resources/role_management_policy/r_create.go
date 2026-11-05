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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.RoleManagementPoliciesClient

			var config RoleManagementPolicyModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			// Resources already have a default policy, so retrieve its ID in order that we can update it
			policyId, err := FindRoleManagementPolicyId(ctx, metadata.Client.Authorization.RoleManagementPoliciesClient, config.Scope, config.RoleDefinitionId)
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *policyId)
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", policyId, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving existing %s: model was nil", policyId)
			}

			model, err := buildRoleManagementPolicyForUpdate(pointer.To(metadata), existing.Model)
			if err != nil {
				return err
			}

			if _, err = client.Update(ctx, *policyId, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", policyId, err)
			}

			// We are using a custom type parse.RoleManagementPolicyId as the ID type for this resource, because the actual
			// resource ID type (ScopedRoleManagementPolicyId) changes each time the policy is updated, so this allows us
			// to search for the latest policy at Read time.
			metadata.SetID(parse.NewRoleManagementPolicyId(config.RoleDefinitionId, config.Scope))
			return nil
		},
	}
}
