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
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roledefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
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

			result, err := client.Get(ctx, id.BaseURI(), id.Scope, id.RoleDefinitionName)
			if err != nil {
				if response.WasNotFound(result.Response.Response) {
					return metadata.MarkAsGone(id)
				}
				return err
			}

			state := KeyVaultMHSMRoleDefinitionModel{
				Name:         pointer.From(result.Name),
				ManagedHSMID: managedHsmId.ID(),
			}

			if v := pointer.From(result.ID); v != "" {
				roleID, err := roledefinitions.ParseScopedRoleDefinitionIDInsensitively(v)
				if err != nil {
					return fmt.Errorf("paring role definition id %q: %+v", v, err)
				}
				state.ResourceManagerId = roleID.ID()
			}

			if prop := result.RoleDefinitionProperties; prop != nil {
				state.Description = pointer.From(prop.Description)
				state.RoleType = string(prop.RoleType)
				state.RoleName = pointer.From(prop.RoleName)
				state.Permission = flattenKeyVaultMHSMRolePermission(prop.Permissions)
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}

func flattenKeyVaultMHSMRolePermission(perms *[]keyvault.Permission) []Permission {
	if perms == nil {
		return make([]Permission, 0)
	}

	res := make([]Permission, 0, len(*perms))
	for _, perm := range *perms {
		var data, notData []string
		for _, item := range pointer.From(perm.DataActions) {
			data = append(data, string(item))
		}
		for _, item := range pointer.From(perm.NotDataActions) {
			notData = append(notData, string(item))
		}

		res = append(res, Permission{
			Actions:        pointer.From(perm.Actions),
			NotActions:     pointer.From(perm.NotActions),
			DataActions:    data,
			NotDataActions: notData,
		})
	}
	return res
}

type Permission struct {
	Actions        []string `tfschema:"actions"`
	NotActions     []string `tfschema:"not_actions"`
	DataActions    []string `tfschema:"data_actions"`
	NotDataActions []string `tfschema:"not_data_actions"`
}
