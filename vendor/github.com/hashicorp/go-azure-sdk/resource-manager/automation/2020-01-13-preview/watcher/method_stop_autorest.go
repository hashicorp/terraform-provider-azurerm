package watcher

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StopOperationResponse struct {
	HttpResponse *http.Response
}

// Stop ...
func (c WatcherClient) Stop(ctx context.Context, id WatcherId) (result StopOperationResponse, err error) {
	req, err := c.preparerForStop(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "watcher.WatcherClient", "Stop", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "watcher.WatcherClient", "Stop", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForStop(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "watcher.WatcherClient", "Stop", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForStop prepares the Stop request.
func (c WatcherClient) preparerForStop(ctx context.Context, id WatcherId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/stop", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForStop handles the response to the Stop request. The method always
// closes the http.Response Body.
func (c WatcherClient) responderForStop(resp *http.Response) (result StopOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
