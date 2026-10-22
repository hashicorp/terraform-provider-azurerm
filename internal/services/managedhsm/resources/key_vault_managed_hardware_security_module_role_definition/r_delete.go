// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module_role_definition

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
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

			// TODO: @manicminer: when migrating to go-azure-sdk, the SDK should auto-retry on 409 responses
			// (these occur when a recently deleted assignment for the role has not yet fully replicated)

			if _, err = client.Delete(ctx, id.BaseURI(), id.Scope, id.RoleDefinitionName); err != nil {
				return fmt.Errorf("deleting %+v: %v", id, err)
			}

			poller := custompollers.NewEventualConsistencyPoller(5, func(pollerCtx context.Context) (*http.Response, error) {
				result, err := client.Get(pollerCtx, id.BaseURI(), id.Scope, id.RoleDefinitionName)
				return result.Response.Response, err
			}, &custompollers.EventualConsistencyPollerOptions{
				Interval:         5 * time.Second,
				TargetStatusCode: pointer.To(http.StatusNotFound),
			})
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
			}

			return nil
		},
	}
}
