// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module_role_definition

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: time.Minute * 10,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) (err error) {
			client := metadata.Client.ManagedHSMs.DataPlaneRoleDefinitionsClient
			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			id, err := parse.ManagedHSMDataPlaneRoleDefinitionID(metadata.ResourceData.Id(), domainSuffix)
			if err != nil {
				return err
			}

			subscriptionId := commonids.NewSubscriptionID(metadata.Client.Account.SubscriptionId)
			managedHsmId, err := metadata.Client.ManagedHSMs.ManagedHSMIDFromBaseUrl(ctx, subscriptionId, id.BaseURI(), domainSuffix)
			if err != nil {
				return fmt.Errorf("determining the Managed HSM ID from the Base URI %q: %+v", id.BaseURI(), err)
			}
			if managedHsmId == nil {
				return fmt.Errorf("unable to determine the Managed HSM ID from the Base URI %q: %+v", id.BaseURI(), err)
			}

			locks.ByName(managedHsmId.ID(), "azurerm_key_vault_managed_hardware_security_module")
			defer locks.UnlockByName(managedHsmId.ID(), "azurerm_key_vault_managed_hardware_security_module")

			var model KeyVaultMHSMRoleDefinitionModel
			if err = metadata.Decode(&model); err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.BaseURI(), id.Scope, id.RoleDefinitionName)
			if err != nil {
				if response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("not found resource to update: %s", id)
				}
				return fmt.Errorf("retrieving role definition by name %s: %v", model.Name, err)
			}

			payload := keyvault.RoleDefinitionCreateParameters{
				Properties: &keyvault.RoleDefinitionProperties{
					RoleName:    pointer.To(model.RoleName),
					Description: pointer.To(model.Description),
					RoleType:    keyvault.RoleTypeCustomRole,
					Permissions: expandKeyVaultMHSMRolePermissions(model.Permission),
				},
			}

			if _, err = client.CreateOrUpdate(ctx, id.BaseURI(), id.Scope, id.RoleDefinitionName, payload); err != nil {
				return fmt.Errorf("updating %s: %v", id.ID(), err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
