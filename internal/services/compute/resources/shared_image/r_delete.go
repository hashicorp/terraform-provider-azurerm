// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package shared_image

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryimages"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	computeClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSharedImageDelete(d *pluginsdk.ResourceData, meta any) error {
	client := computeClient.From(meta).GalleryImagesClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := galleryimages.ParseGalleryImageID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	// Deletion is eventually consistent, https://github.com/Azure/azure-sdk-for-go/issues/8314
	poller := custompollers.NewEventualConsistencyPoller(10, func(pollerCtx context.Context) (*http.Response, error) {
		resp, err := client.Get(pollerCtx, *id)
		return resp.HttpResponse, err
	}, custompollers.DefaultDeletionEventualConsistencyPollerOptions())
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling for deletion of %s: %+v", id, err)
	}

	return nil
}
