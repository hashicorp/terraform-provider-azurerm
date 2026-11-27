// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cognitive_account

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	clientsAccount "github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/client"
	networkClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCognitiveAccountDelete(d *pluginsdk.ResourceData, meta any) error {
	accountsClient := client.From(meta).AccountsClient
	deletedAccountsClient := client.From(meta).AccountsClient
	ctx, cancel := timeouts.ForDelete(clientsAccount.StopContext(meta), d)
	defer cancel()

	id, err := cognitiveservicesaccounts.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	// first we need to retrieve it, since we need the location to be able to purge it
	account, err := accountsClient.AccountsGet(ctx, *id)
	if err != nil || account.Model == nil || account.Model.Location == nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	deletedAccountId := cognitiveservicesaccounts.NewDeletedAccountID(id.SubscriptionId, *account.Model.Location, id.ResourceGroupName, id.AccountName)

	if err := accountsClient.AccountsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if features.From(meta).CognitiveAccount.PurgeSoftDeleteOnDestroy {
		log.Printf("[DEBUG] Purging %s..", *id)
		if err := deletedAccountsClient.DeletedAccountsPurgeThenPoll(ctx, deletedAccountId); err != nil {
			return fmt.Errorf("purging %s: %+v", *id, err)
		}
	} else {
		log.Printf("[DEBUG] Skipping Purge of %s", *id)
	}

	// If `network_injection` is configured, wait for Service Association Link (SAL) to be removed from the agent subnet
	// This is a service issue workaround as the SAL is not removed immediately after the Cognitive Account is deleted
	var subnetId *commonids.SubnetId
	if account.Model.Properties != nil && account.Model.Properties.NetworkInjections != nil {
		networkInjections := *account.Model.Properties.NetworkInjections
		if len(networkInjections) > 0 {
			if subnetIdStr := pointer.From(networkInjections[0].SubnetArmId); subnetIdStr != "" {
				parsedSubnetId, err := commonids.ParseSubnetIDInsensitively(subnetIdStr)
				if err != nil {
					return err
				}
				subnetId = parsedSubnetId
			}
		}
	}

	if subnetId != nil {
		deadline, ok := ctx.Deadline()
		if !ok {
			return errors.New("internal-error: context had no deadline")
		}

		subnetClient := networkClient.From(meta).Subnets
		log.Printf("[DEBUG] Waiting for Service Association Links to be removed from subnet %s", subnetId.ID())

		stateConf := &pluginsdk.StateChangeConf{
			Pending:    []string{"SALExists"},
			Target:     []string{"SALDeleted"},
			Refresh:    serviceAssociationLinkStateRefreshFunc(ctx, subnetClient, *subnetId),
			MinTimeout: 15 * time.Second,
			Timeout:    time.Until(deadline),
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for Service Association Links to be removed from %s: %+v", *subnetId, err)
		}
	}

	return nil
}
