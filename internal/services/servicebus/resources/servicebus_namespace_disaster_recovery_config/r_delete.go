// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_namespace_disaster_recovery_config

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/armdisasterrecoveries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceBusNamespaceDisasterRecoveryConfigDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ServiceBus.ArmDisasterRecoveriesClient
	namespacesClient := meta.(*clients.Client).ServiceBus.NamespacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := armdisasterrecoveries.ParseDisasterRecoveryConfigID(d.Id())
	if err != nil {
		return err
	}

	// @tombuildsstuff: whilst we previously checked the 200 response, since that's the only valid status
	// code defined in the Swagger, anything else would raise an error thus the check is superfluous
	if _, err := client.DisasterRecoveryConfigsBreakPairing(ctx, *id); err != nil {
		return fmt.Errorf("breaking pairing %s: %+v", id, err)
	}

	if err := resourceServiceBusNamespaceDisasterRecoveryConfigWaitForState(ctx, client, *id); err != nil {
		return fmt.Errorf("waiting for the pairing to break for %s: %+v", *id, err)
	}

	if _, err := client.DisasterRecoveryConfigsDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	// no future for deletion so wait for it to vanish
	deletePoller := custompollers.NewEventualConsistencyPoller(1, func(pollerCtx context.Context) (*http.Response, error) {
		resp, err := client.DisasterRecoveryConfigsGet(pollerCtx, *id)
		return resp.HttpResponse, err
	}, &custompollers.EventualConsistencyPollerOptions{
		Interval:         30 * time.Second,
		TargetStatusCode: pointer.To(http.StatusNotFound),
	})
	if err := deletePoller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting the deletion of %s: %v", *id, err)
	}

	namespaceId := namespaces.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
	// it can take some time for the name to become available again
	// this is mainly here 	to enable updating the resource in place
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	nameFreeWait := &pluginsdk.StateChangeConf{
		Pending:    []string{"NameInUse"},
		Target:     []string{"None"},
		MinTimeout: 30 * time.Second,
		Timeout:    time.Until(deadline),
		Refresh: func() (any, string, error) {
			resp, err := namespacesClient.DisasterRecoveryConfigsCheckNameAvailability(ctx, namespaceId, namespaces.CheckNameAvailability{
				Name: id.DisasterRecoveryConfigName,
			})
			if err != nil {
				return resp, "Error", fmt.Errorf("checking for the status of %s: %+v", *id, err)
			}

			reason := ""
			if model := resp.Model; model != nil {
				if v := model.Reason; v != nil {
					reason = string(*v)
				}
			}
			return resp, reason, nil
		},
	}

	if _, err := nameFreeWait.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("checking if the name for %s has become free: %v", *id, err)
	}

	return nil
}
