// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_linked_service

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsLinkedServiceDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LogAnalytics.LinkedServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := linkedservices.ParseLinkedServiceID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	// (@WodansSon) - This is a bug in the service API, it returns instantly from the delete call with a 200
	// so we must wait for the state to change before we return from the delete function
	deleteWait := logAnalyticsLinkedServiceDeleteWaitForState(ctx, client, d.Timeout(pluginsdk.TimeoutDelete), *id)

	if _, err := deleteWait.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s: %+v", *id, err)
	}

	return nil
}
