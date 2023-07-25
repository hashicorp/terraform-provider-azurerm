package updateruns

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

type UpdateRunsDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// UpdateRunsDelete ...
func (c UpdateRunsClient) UpdateRunsDelete(ctx context.Context, id UpdateRunId) (result UpdateRunsDeleteOperationResponse, err error) {
	req, err := c.preparerForUpdateRunsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "updateruns.UpdateRunsClient", "UpdateRunsDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForUpdateRunsDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "updateruns.UpdateRunsClient", "UpdateRunsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// UpdateRunsDeleteThenPoll performs UpdateRunsDelete then polls until it's completed
func (c UpdateRunsClient) UpdateRunsDeleteThenPoll(ctx context.Context, id UpdateRunId) error {
	result, err := c.UpdateRunsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing UpdateRunsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after UpdateRunsDelete: %+v", err)
	}

	return nil
}

// preparerForUpdateRunsDelete prepares the UpdateRunsDelete request.
func (c UpdateRunsClient) preparerForUpdateRunsDelete(ctx context.Context, id UpdateRunId) (*http.Request, error) {
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

// senderForUpdateRunsDelete sends the UpdateRunsDelete request. The method will close the
// http.Response Body if it receives an error.
func (c UpdateRunsClient) senderForUpdateRunsDelete(ctx context.Context, req *http.Request) (future UpdateRunsDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
