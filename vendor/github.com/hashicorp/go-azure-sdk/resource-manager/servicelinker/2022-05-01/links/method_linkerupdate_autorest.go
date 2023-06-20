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

type LinkerUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// LinkerUpdate ...
func (c LinksClient) LinkerUpdate(ctx context.Context, id ScopedLinkerId, input LinkerPatch) (result LinkerUpdateOperationResponse, err error) {
	req, err := c.preparerForLinkerUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "links.LinksClient", "LinkerUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForLinkerUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "links.LinksClient", "LinkerUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// LinkerUpdateThenPoll performs LinkerUpdate then polls until it's completed
func (c LinksClient) LinkerUpdateThenPoll(ctx context.Context, id ScopedLinkerId, input LinkerPatch) error {
	result, err := c.LinkerUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing LinkerUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after LinkerUpdate: %+v", err)
	}

	return nil
}

// preparerForLinkerUpdate prepares the LinkerUpdate request.
func (c LinksClient) preparerForLinkerUpdate(ctx context.Context, id ScopedLinkerId, input LinkerPatch) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForLinkerUpdate sends the LinkerUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c LinksClient) senderForLinkerUpdate(ctx context.Context, req *http.Request) (future LinkerUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
