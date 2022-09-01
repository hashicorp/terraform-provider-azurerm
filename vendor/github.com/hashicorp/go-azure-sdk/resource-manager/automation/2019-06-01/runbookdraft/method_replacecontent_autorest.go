package runbookdraft

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

type ReplaceContentOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ReplaceContent ...
func (c RunbookDraftClient) ReplaceContent(ctx context.Context, id RunbookId, input interface{}) (result ReplaceContentOperationResponse, err error) {
	req, err := c.preparerForReplaceContent(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "runbookdraft.RunbookDraftClient", "ReplaceContent", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForReplaceContent(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "runbookdraft.RunbookDraftClient", "ReplaceContent", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ReplaceContentThenPoll performs ReplaceContent then polls until it's completed
func (c RunbookDraftClient) ReplaceContentThenPoll(ctx context.Context, id RunbookId, input interface{}) error {
	result, err := c.ReplaceContent(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ReplaceContent: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ReplaceContent: %+v", err)
	}

	return nil
}

// preparerForReplaceContent prepares the ReplaceContent request.
func (c RunbookDraftClient) preparerForReplaceContent(ctx context.Context, id RunbookId, input interface{}) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/draft/content", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForReplaceContent sends the ReplaceContent request. The method will close the
// http.Response Body if it receives an error.
func (c RunbookDraftClient) senderForReplaceContent(ctx context.Context, req *http.Request) (future ReplaceContentOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
