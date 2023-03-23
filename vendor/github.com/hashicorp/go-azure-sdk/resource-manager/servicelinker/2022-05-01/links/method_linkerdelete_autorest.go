package links

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

type LinkerDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// LinkerDelete ...
func (c LinksClient) LinkerDelete(ctx context.Context, id ScopedLinkerId) (result LinkerDeleteOperationResponse, err error) {
	req, err := c.preparerForLinkerDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "links.LinksClient", "LinkerDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForLinkerDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "links.LinksClient", "LinkerDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// LinkerDeleteThenPoll performs LinkerDelete then polls until it's completed
func (c LinksClient) LinkerDeleteThenPoll(ctx context.Context, id ScopedLinkerId) error {
	result, err := c.LinkerDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing LinkerDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after LinkerDelete: %+v", err)
	}

	return nil
}

// preparerForLinkerDelete prepares the LinkerDelete request.
func (c LinksClient) preparerForLinkerDelete(ctx context.Context, id ScopedLinkerId) (*http.Request, error) {
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

// senderForLinkerDelete sends the LinkerDelete request. The method will close the
// http.Response Body if it receives an error.
func (c LinksClient) senderForLinkerDelete(ctx context.Context, req *http.Request) (future LinkerDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
