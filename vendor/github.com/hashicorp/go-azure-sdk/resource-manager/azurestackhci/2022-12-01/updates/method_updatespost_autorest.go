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

type UpdatesPostOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// UpdatesPost ...
func (c UpdatesClient) UpdatesPost(ctx context.Context, id UpdateId) (result UpdatesPostOperationResponse, err error) {
	req, err := c.preparerForUpdatesPost(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "updates.UpdatesClient", "UpdatesPost", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForUpdatesPost(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "updates.UpdatesClient", "UpdatesPost", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// UpdatesPostThenPoll performs UpdatesPost then polls until it's completed
func (c UpdatesClient) UpdatesPostThenPoll(ctx context.Context, id UpdateId) error {
	result, err := c.UpdatesPost(ctx, id)
	if err != nil {
		return fmt.Errorf("performing UpdatesPost: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after UpdatesPost: %+v", err)
	}

	return nil
}

// preparerForUpdatesPost prepares the UpdatesPost request.
func (c UpdatesClient) preparerForUpdatesPost(ctx context.Context, id UpdateId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/apply", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForUpdatesPost sends the UpdatesPost request. The method will close the
// http.Response Body if it receives an error.
func (c UpdatesClient) senderForUpdatesPost(ctx context.Context, req *http.Request) (future UpdatesPostOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
