package gallerysharingupdate

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GallerySharingProfileUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// GallerySharingProfileUpdate ...
func (c GallerySharingUpdateClient) GallerySharingProfileUpdate(ctx context.Context, id GalleryId, input SharingUpdate) (result GallerySharingProfileUpdateOperationResponse, err error) {
	req, err := c.preparerForGallerySharingProfileUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "gallerysharingupdate.GallerySharingUpdateClient", "GallerySharingProfileUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForGallerySharingProfileUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "gallerysharingupdate.GallerySharingUpdateClient", "GallerySharingProfileUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// GallerySharingProfileUpdateThenPoll performs GallerySharingProfileUpdate then polls until it's completed
func (c GallerySharingUpdateClient) GallerySharingProfileUpdateThenPoll(ctx context.Context, id GalleryId, input SharingUpdate) error {
	result, err := c.GallerySharingProfileUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing GallerySharingProfileUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after GallerySharingProfileUpdate: %+v", err)
	}

	return nil
}

// preparerForGallerySharingProfileUpdate prepares the GallerySharingProfileUpdate request.
func (c GallerySharingUpdateClient) preparerForGallerySharingProfileUpdate(ctx context.Context, id GalleryId, input SharingUpdate) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/share", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForGallerySharingProfileUpdate sends the GallerySharingProfileUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c GallerySharingUpdateClient) senderForGallerySharingProfileUpdate(ctx context.Context, req *http.Request) (future GallerySharingProfileUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
