// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package maps_account

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maps/2023-06-01/accounts"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maps/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMapsAccountUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Maps.AccountsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := accounts.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("local_authentication_enabled") {
		payload.Properties.DisableLocalAuth = pointer.To(!d.Get("local_authentication_enabled").(bool))
	}

	if d.HasChange("cors") {
		payload.Properties.Cors = expandCors(d.Get("cors").([]any))
	}

	if d.HasChange("data_store") {
		dataStores, err := expandDataStore(d.Get("data_store").([]any))
		if err != nil {
			return fmt.Errorf("expanding `data_store`: %+v", err)
		}
		payload.Properties.LinkedResources = dataStores
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if _, err := client.CreateOrUpdate(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	// These should actually be LROs, but they're not, custom poller is required until https://github.com/Azure/azure-rest-api-specs/issues/29501 is resolved
	pollerType := custompollers.NewMapsAccountPoller(client, *id)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return err
	}

	return resourceMapsAccountRead(d, meta)
}
