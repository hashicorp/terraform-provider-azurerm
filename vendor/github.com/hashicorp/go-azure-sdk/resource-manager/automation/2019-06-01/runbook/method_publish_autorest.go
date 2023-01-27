package runbook

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

type PublishOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Publish ...
func (c RunbookClient) Publish(ctx context.Context, id RunbookId) (result PublishOperationResponse, err error) {
	req, err := c.preparerForPublish(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "runbook.RunbookClient", "Publish", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPublish(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "runbook.RunbookClient", "Publish", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PublishThenPoll performs Publish then polls until it's completed
func (c RunbookClient) PublishThenPoll(ctx context.Context, id RunbookId) error {
	result, err := c.Publish(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Publish: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Publish: %+v", err)
	}

	return nil
}

// preparerForPublish prepares the Publish request.
func (c RunbookClient) preparerForPublish(ctx context.Context, id RunbookId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/publish", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForPublish sends the Publish request. The method will close the
// http.Response Body if it receives an error.
func (c RunbookClient) senderForPublish(ctx context.Context, req *http.Request) (future PublishOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
