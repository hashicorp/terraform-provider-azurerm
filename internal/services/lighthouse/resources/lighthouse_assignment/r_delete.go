// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lighthouse_assignment

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedservices/2022-10-01/registrationassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	lighthouseClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/lighthouse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLighthouseAssignmentDelete(d *pluginsdk.ResourceData, meta any) error {
	client := lighthouseClient.From(meta).AssignmentsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := registrationassignments.ParseScopedRegistrationAssignmentID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	poller := custompollers.NewEventualConsistencyPoller(1, func(pollerCtx context.Context) (*http.Response, error) {
		resp, err := client.Get(pollerCtx, *id, registrationassignments.GetOperationOptions{ExpandRegistrationDefinition: pointer.To(true)})
		return resp.HttpResponse, err
	}, &custompollers.EventualConsistencyPollerOptions{
		Interval:         15 * time.Second,
		TargetStatusCode: pointer.To(http.StatusNotFound),
	})
	if err := poller.PollUntilDone(ctx); err != nil {
		return err
	}

	return nil
}
