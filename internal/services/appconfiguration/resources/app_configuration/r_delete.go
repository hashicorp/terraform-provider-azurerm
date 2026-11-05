// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_configuration

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/configurationstores"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/deletedconfigurationstores"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppConfigurationDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppConfiguration.ConfigurationStoresClient
	deletedConfigurationStoresClient := meta.(*clients.Client).AppConfiguration.DeletedConfigurationStoresClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	id, err := configurationstores.ParseConfigurationStoreID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %q: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %q: `properties` was nil", *id)
	}

	purgeProtectionEnabled := pointer.From(existing.Model.Properties.EnablePurgeProtection)
	softDeleteEnabled := false
	if sde := existing.Model.Properties.SoftDeleteRetentionInDays; sde != nil && *sde > 0 {
		softDeleteEnabled = true
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if meta.(*clients.Client).Features.AppConfiguration.PurgeSoftDeleteOnDestroy && softDeleteEnabled {
		deletedId := deletedconfigurationstores.NewDeletedConfigurationStoreID(subscriptionId, existing.Model.Location, id.ConfigurationStoreName)

		// AppConfiguration with Purge Protection Enabled cannot be deleted unless done by Azure
		if purgeProtectionEnabled {
			deletedInfo, err := deletedConfigurationStoresClient.ConfigurationStoresGetDeleted(ctx, deletedId)
			if err != nil {
				return fmt.Errorf("while purging the soft-deleted, retrieving the Deletion Details for %s: %+v", *id, err)
			}

			if deletedInfo.Model != nil && deletedInfo.Model.Properties != nil && deletedInfo.Model.Properties.DeletionDate != nil && deletedInfo.Model.Properties.ScheduledPurgeDate != nil {
				log.Printf("[DEBUG] The App Configuration %q has Purge Protection Enabled and was deleted on %q. Azure will purge this on %q",
					id.ConfigurationStoreName, *deletedInfo.Model.Properties.DeletionDate, *deletedInfo.Model.Properties.ScheduledPurgeDate)
			} else {
				log.Printf("[DEBUG] The App Configuration %q has Purge Protection Enabled and will be purged automatically by Azure", id.ConfigurationStoreName)
			}
			return nil
		}

		log.Printf("[DEBUG]  %q marked for purge - executing purge", id.ConfigurationStoreName)

		if _, err := deletedConfigurationStoresClient.ConfigurationStoresPurgeDeleted(ctx, deletedId); err != nil {
			return fmt.Errorf("purging %s: %+v", *id, err)
		}

		// The PurgeDeleted API is a POST which returns a 200 with no body and nothing to poll on, so we'll need
		// a custom poller to poll until the LRO returns a 404
		pollerType := &purgeDeletedPoller{
			client: deletedConfigurationStoresClient,
			id:     deletedId,
		}
		poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("polling after purging for %s: %+v", *id, err)
		}

		// retry checkNameAvailability until the name is released by purged app configuration, see https://github.com/Azure/AppConfiguration/issues/677
		operationsClient := meta.(*clients.Client).AppConfiguration.OperationsClient
		if err = resourceConfigurationStoreWaitForNameAvailable(ctx, operationsClient, *id); err != nil {
			return err
		}
		log.Printf("[DEBUG] Purged AppConfiguration %q.", id.ConfigurationStoreName)
	}

	meta.(*clients.Client).AppConfiguration.RemoveFromCache(*id)

	return nil
}
