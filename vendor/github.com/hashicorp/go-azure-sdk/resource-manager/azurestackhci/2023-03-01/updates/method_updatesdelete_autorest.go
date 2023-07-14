package updates

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

type UpdatesDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// UpdatesDelete ...
func (c UpdatesClient) UpdatesDelete(ctx context.Context, id UpdateId) (result UpdatesDeleteOperationResponse, err error) {
	req, err := c.preparerForUpdatesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "updates.UpdatesClient", "UpdatesDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForUpdatesDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "updates.UpdatesClient", "UpdatesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// UpdatesDeleteThenPoll performs UpdatesDelete then polls until it's completed
func (c UpdatesClient) UpdatesDeleteThenPoll(ctx context.Context, id UpdateId) error {
	result, err := c.UpdatesDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing UpdatesDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after UpdatesDelete: %+v", err)
	}

	return nil
}

// preparerForUpdatesDelete prepares the UpdatesDelete request.
func (c UpdatesClient) preparerForUpdatesDelete(ctx context.Context, id UpdateId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForUpdatesDelete sends the UpdatesDelete request. The method will close the
// http.Response Body if it receives an error.
func (c UpdatesClient) senderForUpdatesDelete(ctx context.Context, req *http.Request) (future UpdatesDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
