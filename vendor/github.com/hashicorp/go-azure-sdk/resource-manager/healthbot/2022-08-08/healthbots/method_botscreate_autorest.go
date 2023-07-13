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

type BotsCreateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// BotsCreate ...
func (c HealthbotsClient) BotsCreate(ctx context.Context, id HealthBotId, input HealthBot) (result BotsCreateOperationResponse, err error) {
	req, err := c.preparerForBotsCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsCreate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForBotsCreate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsCreate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// BotsCreateThenPoll performs BotsCreate then polls until it's completed
func (c HealthbotsClient) BotsCreateThenPoll(ctx context.Context, id HealthBotId, input HealthBot) error {
	result, err := c.BotsCreate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing BotsCreate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after BotsCreate: %+v", err)
	}

	return nil
}

// preparerForBotsCreate prepares the BotsCreate request.
func (c HealthbotsClient) preparerForBotsCreate(ctx context.Context, id HealthBotId, input HealthBot) (*http.Request, error) {
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

// senderForBotsCreate sends the BotsCreate request. The method will close the
// http.Response Body if it receives an error.
func (c HealthbotsClient) senderForBotsCreate(ctx context.Context, req *http.Request) (future BotsCreateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
