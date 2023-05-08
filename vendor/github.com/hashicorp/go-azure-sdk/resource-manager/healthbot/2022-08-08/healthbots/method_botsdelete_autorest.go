package healthbots

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

type BotsDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// BotsDelete ...
func (c HealthbotsClient) BotsDelete(ctx context.Context, id HealthBotId) (result BotsDeleteOperationResponse, err error) {
	req, err := c.preparerForBotsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForBotsDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// BotsDeleteThenPoll performs BotsDelete then polls until it's completed
func (c HealthbotsClient) BotsDeleteThenPoll(ctx context.Context, id HealthBotId) error {
	result, err := c.BotsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing BotsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after BotsDelete: %+v", err)
	}

	return nil
}

// preparerForBotsDelete prepares the BotsDelete request.
func (c HealthbotsClient) preparerForBotsDelete(ctx context.Context, id HealthBotId) (*http.Request, error) {
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

// senderForBotsDelete sends the BotsDelete request. The method will close the
// http.Response Body if it receives an error.
func (c HealthbotsClient) senderForBotsDelete(ctx context.Context, req *http.Request) (future BotsDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
