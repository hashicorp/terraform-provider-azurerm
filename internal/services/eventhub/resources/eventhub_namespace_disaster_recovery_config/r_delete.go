// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace_disaster_recovery_config

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/checknameavailabilitydisasterrecoveryconfigs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/disasterrecoveryconfigs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventHubNamespaceDisasterRecoveryConfigDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Eventhub.DisasterRecoveryConfigsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := disasterrecoveryconfigs.ParseDisasterRecoveryConfigID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.NamespaceName, eventHubNamespaceResourceName)
	defer locks.UnlockByName(id.NamespaceName, eventHubNamespaceResourceName)

	pairingStatus, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("checking the status of eventhub disaster recovery error: %+v", err)
	}

	// need to check if DCR needs pair-breaking first
	breakPairFirst := false
	if model := pairingStatus.Model; model != nil {
		if model.Properties != nil {
			if model.Properties.PartnerNamespace != nil && *model.Properties.PartnerNamespace != "" {
				breakPairFirst = true
			}
		}
	}

	if breakPairFirst {
		if _, err := client.BreakPairing(ctx, *id); err != nil {
			return fmt.Errorf("breaking pairing of %s: %+v", *id, err)
		}
		if err := resourceEventHubNamespaceDisasterRecoveryConfigWaitForState(ctx, client, *id); err != nil {
			return fmt.Errorf("waiting for pairing to break for %s: %+v", *id, err)
		}
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	// no future for deletion so wait for it to vanish
	deletePoller := custompollers.NewEventualConsistencyPoller(1, func(pollerCtx context.Context) (*http.Response, error) {
		resp, err := client.Get(pollerCtx, *id)
		return resp.HttpResponse, err
	}, &custompollers.EventualConsistencyPollerOptions{
		Interval:         30 * time.Second,
		TargetStatusCode: pointer.To(http.StatusNotFound),
	})
	if err := deletePoller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting the deletion of %s: %+v", *id, err)
	}

	// it can take some time for the name to become available again
	// this is mainly here	to enable updating the resource in place
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context has no deadline")
	}
	parentNamespaceId := checknameavailabilitydisasterrecoveryconfigs.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
	availabilityClient := meta.(*clients.Client).Eventhub.DisasterRecoveryNameAvailabilityClient
	nameFreeWait := &pluginsdk.StateChangeConf{
		Pending:    []string{"NameInUse"},
		Target:     []string{"None"},
		MinTimeout: 30 * time.Second,
		Timeout:    time.Until(deadline),
		Refresh: func() (any, string, error) {
			input := checknameavailabilitydisasterrecoveryconfigs.CheckNameAvailabilityParameter{
				Name: id.DisasterRecoveryConfigName,
			}
			resp, err := availabilityClient.DisasterRecoveryConfigsCheckNameAvailability(ctx, parentNamespaceId, input)
			if err != nil {
				return resp, "Error", fmt.Errorf("waiting for the name of %s to become free: %v", *id, err)
			}
			if resp.Model == nil || resp.Model.Reason == nil {
				return resp, "Error", fmt.Errorf("`model` or `model.Reason` was nil")
			}
			return resp, string(*resp.Model.Reason), nil
		},
	}

	if _, err := nameFreeWait.WaitForStateContext(ctx); err != nil {
		return err
	}

	return nil
}
