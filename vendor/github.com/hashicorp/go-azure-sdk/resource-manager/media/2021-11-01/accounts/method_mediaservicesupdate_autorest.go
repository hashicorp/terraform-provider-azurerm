package accounts

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

type MediaservicesUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// MediaservicesUpdate ...
func (c AccountsClient) MediaservicesUpdate(ctx context.Context, id MediaServiceId, input MediaServiceUpdate) (result MediaservicesUpdateOperationResponse, err error) {
	req, err := c.preparerForMediaservicesUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForMediaservicesUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// MediaservicesUpdateThenPoll performs MediaservicesUpdate then polls until it's completed
func (c AccountsClient) MediaservicesUpdateThenPoll(ctx context.Context, id MediaServiceId, input MediaServiceUpdate) error {
	result, err := c.MediaservicesUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing MediaservicesUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after MediaservicesUpdate: %+v", err)
	}

	return nil
}

// preparerForMediaservicesUpdate prepares the MediaservicesUpdate request.
func (c AccountsClient) preparerForMediaservicesUpdate(ctx context.Context, id MediaServiceId, input MediaServiceUpdate) (*http.Request, error) {
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

// senderForMediaservicesUpdate sends the MediaservicesUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c AccountsClient) senderForMediaservicesUpdate(ctx context.Context, req *http.Request) (future MediaservicesUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
