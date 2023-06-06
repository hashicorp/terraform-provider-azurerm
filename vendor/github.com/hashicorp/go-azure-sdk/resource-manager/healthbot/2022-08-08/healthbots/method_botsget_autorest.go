package healthbots

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BotsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *HealthBot
}

// BotsGet ...
func (c HealthbotsClient) BotsGet(ctx context.Context, id HealthBotId) (result BotsGetOperationResponse, err error) {
	req, err := c.preparerForBotsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForBotsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForBotsGet prepares the BotsGet request.
func (c HealthbotsClient) preparerForBotsGet(ctx context.Context, id HealthBotId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForBotsGet handles the response to the BotsGet request. The method always
// closes the http.Response Body.
func (c HealthbotsClient) responderForBotsGet(resp *http.Response) (result BotsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
