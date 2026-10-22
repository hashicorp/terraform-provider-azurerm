// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module_role_assignment

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/managedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedHSMs.DataPlaneRoleAssignmentsClient
			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			var config KeyVaultManagedHSMRoleAssignmentModel
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
					return fmt.Errorf("parsing the Data Plane Endpoint %q: %+v", pointer.From(endpoint), err)
				}
			}

			locks.ByName(managedHsmId.ID(), "azurerm_key_vault_managed_hardware_security_module")
			defer locks.UnlockByName(managedHsmId.ID(), "azurerm_key_vault_managed_hardware_security_module")

			id := parse.NewManagedHSMDataPlaneRoleAssignmentID(endpoint.ManagedHSMName, endpoint.DomainSuffix, config.Scope, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, endpoint.BaseURI(), config.Scope, config.Name)
				if !response.WasNotFound(existing.Response.Response) {
					if err != nil {
						return fmt.Errorf("retrieving %s: %v", id.ID(), err)
					}
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			var param keyvault.RoleAssignmentCreateParameters
			param.Properties = &keyvault.RoleAssignmentProperties{
				PrincipalID: pointer.To(config.PrincipalId),
				// the role definition id may have '/' prefix, but the api doesn't accept it
				RoleDefinitionID: pointer.To(strings.TrimPrefix(config.RoleDefinitionId, "/")),
			}

			//nolint:misspell
			// TODO: @manicminer: when migrating to go-azure-sdk, the SDK should auto-retry on 400 responses with code "BadParameter" and message "Unkown role definition" (note the misspelling)

			if _, err = client.Create(ctx, endpoint.BaseURI(), config.Scope, config.Name, param); err != nil {
				return fmt.Errorf("creating %s: %v", id.ID(), err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
