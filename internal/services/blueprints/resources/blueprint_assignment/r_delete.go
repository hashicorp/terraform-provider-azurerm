// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package blueprint_assignment

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/resource-manager/blueprints/2018-11-01-preview/assignment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBlueprintAssignmentDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Blueprints.AssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := assignment.ParseScopedBlueprintAssignmentID(d.Id())
	if err != nil {
		return err
	}

	// We use none here to align the previous behaviour of the blueprint resource
	// TODO: we could add a features flag for the blueprint to empower terraform when deleting the blueprint to delete all the generated resources as well
	if _, err := client.Delete(ctx, *id, assignment.DeleteOperationOptions{}); err != nil {
		return fmt.Errorf("failed to delete Blueprint Assignment %q from scope %q: %+v", id.BlueprintAssignmentName, id.ResourceScope, err)
	}

	poller := custompollers.NewEventualConsistencyPoller(1, func(pollerCtx context.Context) (*http.Response, error) {
		resp, err := client.Get(pollerCtx, *id)
		return resp.HttpResponse, err
	}, custompollers.DefaultDeletionEventualConsistencyPollerOptions())
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}

	return nil
}
