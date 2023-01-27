package servicelinker

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

type LinkerCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// LinkerCreateOrUpdate ...
func (c ServiceLinkerClient) LinkerCreateOrUpdate(ctx context.Context, id ScopedLinkerId, input LinkerResource) (result LinkerCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForLinkerCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicelinker.ServiceLinkerClient", "LinkerCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForLinkerCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicelinker.ServiceLinkerClient", "LinkerCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// LinkerCreateOrUpdateThenPoll performs LinkerCreateOrUpdate then polls until it's completed
func (c ServiceLinkerClient) LinkerCreateOrUpdateThenPoll(ctx context.Context, id ScopedLinkerId, input LinkerResource) error {
	result, err := c.LinkerCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing LinkerCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after LinkerCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForLinkerCreateOrUpdate prepares the LinkerCreateOrUpdate request.
func (c ServiceLinkerClient) preparerForLinkerCreateOrUpdate(ctx context.Context, id ScopedLinkerId, input LinkerResource) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForLinkerCreateOrUpdate sends the LinkerCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c ServiceLinkerClient) senderForLinkerCreateOrUpdate(ctx context.Context, req *http.Request) (future LinkerCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
