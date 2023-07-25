package healthbots

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BotsRegenerateApiJwtSecretOperationResponse struct {
	HttpResponse *http.Response
	Model        *HealthBotKey
}

// BotsRegenerateApiJwtSecret ...
func (c HealthbotsClient) BotsRegenerateApiJwtSecret(ctx context.Context, id HealthBotId) (result BotsRegenerateApiJwtSecretOperationResponse, err error) {
	req, err := c.preparerForBotsRegenerateApiJwtSecret(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsRegenerateApiJwtSecret", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsRegenerateApiJwtSecret", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForBotsRegenerateApiJwtSecret(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsRegenerateApiJwtSecret", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForBotsRegenerateApiJwtSecret prepares the BotsRegenerateApiJwtSecret request.
func (c HealthbotsClient) preparerForBotsRegenerateApiJwtSecret(ctx context.Context, id HealthBotId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/regenerateApiJwtSecret", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForBotsRegenerateApiJwtSecret handles the response to the BotsRegenerateApiJwtSecret request. The method always
// closes the http.Response Body.
func (c HealthbotsClient) responderForBotsRegenerateApiJwtSecret(resp *http.Response) (result BotsRegenerateApiJwtSecretOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
