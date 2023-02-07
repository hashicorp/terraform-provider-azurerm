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

type MediaservicesCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// MediaservicesCreateOrUpdate ...
func (c AccountsClient) MediaservicesCreateOrUpdate(ctx context.Context, id MediaServiceId, input MediaService) (result MediaservicesCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForMediaservicesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForMediaservicesCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// MediaservicesCreateOrUpdateThenPoll performs MediaservicesCreateOrUpdate then polls until it's completed
func (c AccountsClient) MediaservicesCreateOrUpdateThenPoll(ctx context.Context, id MediaServiceId, input MediaService) error {
	result, err := c.MediaservicesCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing MediaservicesCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after MediaservicesCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForMediaservicesCreateOrUpdate prepares the MediaservicesCreateOrUpdate request.
func (c AccountsClient) preparerForMediaservicesCreateOrUpdate(ctx context.Context, id MediaServiceId, input MediaService) (*http.Request, error) {
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

// senderForMediaservicesCreateOrUpdate sends the MediaservicesCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c AccountsClient) senderForMediaservicesCreateOrUpdate(ctx context.Context, req *http.Request) (future MediaservicesCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
