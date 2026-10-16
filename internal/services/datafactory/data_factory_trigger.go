// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

// startDataFactoryEventTrigger provisions the Event Grid subscription and waits for it to reach Enabled
// before starting the trigger, because Start fails while the subscription is still provisioning
func startDataFactoryEventTrigger(ctx context.Context, client datafactory.TriggersClient, id parse.TriggerId) error {
	subscribeFuture, err := client.SubscribeToEvents(ctx, id.ResourceGroup, id.FactoryName, id.Name)
	if err != nil {
		return fmt.Errorf("subscribing %s to events: %+v", id, err)
	}
	if err = subscribeFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to subscribe to events: %+v", id, err)
	}

	poller := pollers.NewPoller(custompollers.NewDataFactoryEventSubscriptionPoller(client, id), 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for event subscription of %s to become enabled: %+v", id, err)
	}
	startFuture, err := client.Start(ctx, id.ResourceGroup, id.FactoryName, id.Name)
	if err != nil {
		return fmt.Errorf("starting %s: %+v", id, err)
	}
	if err := startFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for start of %s: %+v", id, err)
	}

	return nil
}
