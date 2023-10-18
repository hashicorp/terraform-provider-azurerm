package blobcontainers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ObjectLevelWormOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ObjectLevelWorm ...
func (c BlobContainersClient) ObjectLevelWorm(ctx context.Context, id commonids.StorageContainerId) (result ObjectLevelWormOperationResponse, err error) {
	req, err := c.preparerForObjectLevelWorm(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "ObjectLevelWorm", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForObjectLevelWorm(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "ObjectLevelWorm", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ObjectLevelWormThenPoll performs ObjectLevelWorm then polls until it's completed
func (c BlobContainersClient) ObjectLevelWormThenPoll(ctx context.Context, id commonids.StorageContainerId) error {
	result, err := c.ObjectLevelWorm(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ObjectLevelWorm: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ObjectLevelWorm: %+v", err)
	}

	return nil
}

// preparerForObjectLevelWorm prepares the ObjectLevelWorm request.
func (c BlobContainersClient) preparerForObjectLevelWorm(ctx context.Context, id commonids.StorageContainerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/migrate", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForObjectLevelWorm sends the ObjectLevelWorm request. The method will close the
// http.Response Body if it receives an error.
func (c BlobContainersClient) senderForObjectLevelWorm(ctx context.Context, req *http.Request) (future ObjectLevelWormOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
