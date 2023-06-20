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

type BotsUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// BotsUpdate ...
func (c HealthbotsClient) BotsUpdate(ctx context.Context, id HealthBotId, input HealthBotUpdateParameters) (result BotsUpdateOperationResponse, err error) {
	req, err := c.preparerForBotsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForBotsUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// BotsUpdateThenPoll performs BotsUpdate then polls until it's completed
func (c HealthbotsClient) BotsUpdateThenPoll(ctx context.Context, id HealthBotId, input HealthBotUpdateParameters) error {
	result, err := c.BotsUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing BotsUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after BotsUpdate: %+v", err)
	}

	return nil
}

// preparerForBotsUpdate prepares the BotsUpdate request.
func (c HealthbotsClient) preparerForBotsUpdate(ctx context.Context, id HealthBotId, input HealthBotUpdateParameters) (*http.Request, error) {
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

// senderForBotsUpdate sends the BotsUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c HealthbotsClient) senderForBotsUpdate(ctx context.Context, req *http.Request) (future BotsUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
