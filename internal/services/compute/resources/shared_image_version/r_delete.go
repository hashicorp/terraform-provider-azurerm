// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package shared_image_version

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	computeClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSharedImageVersionDelete(d *pluginsdk.ResourceData, meta any) error {
	client := computeClient.From(meta).GalleryImageVersionsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := galleryimageversions.ParseImageVersionID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	// @tombuildsstuff: there appears to be an eventual consistency issue here
	log.Printf("[DEBUG] Waiting for %s to be eventually deleted", *id)
	poller := custompollers.NewEventualConsistencyPoller(10, func(pollerCtx context.Context) (*http.Response, error) {
		resp, err := client.Get(pollerCtx, *id, galleryimageversions.DefaultGetOperationOptions())
		return resp.HttpResponse, err
	}, custompollers.DefaultDeletionEventualConsistencyPollerOptions())
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", *id, err)
	}

	return nil
}
