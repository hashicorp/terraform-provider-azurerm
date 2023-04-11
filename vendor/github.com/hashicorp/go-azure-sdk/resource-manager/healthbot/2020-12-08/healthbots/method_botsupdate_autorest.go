package healthbots

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BotsUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *HealthBot
}

// BotsUpdate ...
func (c HealthbotsClient) BotsUpdate(ctx context.Context, id HealthBotId, input HealthBotUpdateParameters) (result BotsUpdateOperationResponse, err error) {
	req, err := c.preparerForBotsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForBotsUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
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

// responderForBotsUpdate handles the response to the BotsUpdate request. The method always
// closes the http.Response Body.
func (c HealthbotsClient) responderForBotsUpdate(resp *http.Response) (result BotsUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
