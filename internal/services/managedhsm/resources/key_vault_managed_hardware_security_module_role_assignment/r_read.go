// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module_role_assignment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roledefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedHSMs.DataPlaneRoleAssignmentsClient
			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			id, err := parse.ManagedHSMDataPlaneRoleAssignmentID(metadata.ResourceData.Id(), domainSuffix)
			if err != nil {
				return err
			}

			subscriptionId := commonids.NewSubscriptionID(metadata.Client.Account.SubscriptionId)
			resourceManagerId, err := metadata.Client.ManagedHSMs.ManagedHSMIDFromBaseUrl(ctx, subscriptionId, id.BaseURI(), domainSuffix)
			if err != nil {
				return fmt.Errorf("determining Resource Manager ID for %q: %+v", id, err)
			}
			if resourceManagerId == nil {
				return fmt.Errorf("unable to determine the Resource Manager ID for %s", id)
			}

			resp, err := client.Get(ctx, id.BaseURI(), id.Scope, id.RoleAssignmentName)
			if err != nil {
				if response.WasNotFound(resp.Response.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := KeyVaultManagedHSMRoleAssignmentModel{
				ManagedHSMID: resourceManagerId.ID(),
				Name:         id.RoleAssignmentName,
				Scope:        id.Scope,
			}

			if props := resp.Properties; props != nil {
				model.PrincipalId = pointer.From(props.PrincipalID)
				model.ResourceId = pointer.From(resp.ID) // TODO: verify if we should deprecate this

				if roleDefinitionId := pointer.From(props.RoleDefinitionID); roleDefinitionId != "" {
					parsed, err := roledefinitions.ParseScopedRoleDefinitionIDInsensitively(roleDefinitionId)
					if err != nil {
						return fmt.Errorf("parsing role definition id %q: %v", roleDefinitionId, err)
					}

					model.RoleDefinitionId = parsed.ID()
				}
			}

			return metadata.Encode(&model)
		},
	}
}
