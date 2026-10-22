// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module_role_definition

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/managedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedHSMs.DataPlaneRoleDefinitionsClient
			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			var config KeyVaultMHSMRoleDefinitionModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

			var managedHsmId *managedhsms.ManagedHSMId
			var endpoint *parse.ManagedHSMDataPlaneEndpoint
			var err error
			if config.ManagedHSMID != "" {
				managedHsmId, err = managedhsms.ParseManagedHSMID(config.ManagedHSMID)
				if err != nil {
					return err
				}
				baseUri, err := metadata.Client.ManagedHSMs.BaseUriForManagedHSM(ctx, *managedHsmId)
				if err != nil {
					return fmt.Errorf("determining the Data Plane Endpoint for %s: %+v", *managedHsmId, err)
				}
				if baseUri == nil {
					return fmt.Errorf("unable to determine the Data Plane Endpoint for %q", *managedHsmId)
				}
				endpoint, err = parse.ManagedHSMEndpoint(*baseUri, domainSuffix)
				if err != nil {
					return fmt.Errorf("parsing the Data Plane Endpoint %q: %+v", *endpoint, err)
				}
			}

			// need a lock for hsm subresource create/update/delete, or API may respond error as below
			// Status=409 Code="Conflict" Message="There was a conflict while trying to delete the role assignment.
			locks.ByName(managedHsmId.ID(), "azurerm_key_vault_managed_hardware_security_module")
			defer locks.UnlockByName(managedHsmId.ID(), "azurerm_key_vault_managed_hardware_security_module")

			scope := keyvault.RoleScopeGlobal
			id := parse.NewManagedHSMDataPlaneRoleDefinitionID(endpoint.ManagedHSMName, endpoint.DomainSuffix, string(scope), config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id.BaseURI(), id.Scope, id.ManagedHSMName)
				if !response.WasNotFound(existing.Response.Response) {
					if err != nil {
						return fmt.Errorf("checking for the existence of an existing %q: %+v", id, err)
					}
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := keyvault.RoleDefinitionCreateParameters{
				Properties: &keyvault.RoleDefinitionProperties{
					RoleName:    pointer.To(config.RoleName),
					Description: pointer.To(config.Description),
					RoleType:    keyvault.RoleTypeCustomRole,
					Permissions: expandKeyVaultMHSMRolePermissions(config.Permission),
					AssignableScopes: pointer.To([]keyvault.RoleScope{
						scope,
					}),
				},
			}

			// TODO: @manicminer: when migrating to go-azure-sdk, the SDK should auto-retry on 409 responses and should consider manually polling afterwards

			if _, err = client.CreateOrUpdate(ctx, id.BaseURI(), id.Scope, id.RoleDefinitionName, payload); err != nil {
				return fmt.Errorf("creating %s: %v", id.ID(), err)
			}

			poller := custompollers.NewEventualConsistencyPoller(5, func(pollerCtx context.Context) (*http.Response, error) {
				result, err := client.Get(pollerCtx, id.BaseURI(), id.Scope, id.RoleDefinitionName)
				return result.Response.Response, err
			}, &custompollers.EventualConsistencyPollerOptions{
				Interval:              5 * time.Second,
				RetryErrorStatusCodes: []int{http.StatusNotFound},
			})
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandKeyVaultMHSMRolePermissions(perms []Permission) *[]keyvault.Permission {
	res := make([]keyvault.Permission, 0, len(perms))
	for _, perm := range perms {
		dataActions, notDataActions := make([]keyvault.DataAction, 0), make([]keyvault.DataAction, 0)
		for _, data := range perm.DataActions {
			dataActions = append(dataActions, keyvault.DataAction(data))
		}
		for _, notData := range perm.NotDataActions {
			notDataActions = append(notDataActions, keyvault.DataAction(notData))
		}

		res = append(res, keyvault.Permission{
			Actions:        pointer.To(perm.Actions),
			NotActions:     pointer.To(perm.NotActions),
			DataActions:    pointer.To(dataActions),
			NotDataActions: pointer.To(notDataActions),
		})
	}
	return &res
}
