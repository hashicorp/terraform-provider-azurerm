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

type BotsListSecretsOperationResponse struct {
	HttpResponse *http.Response
	Model        *HealthBotKeysResponse
}

// BotsListSecrets ...
func (c HealthbotsClient) BotsListSecrets(ctx context.Context, id HealthBotId) (result BotsListSecretsOperationResponse, err error) {
	req, err := c.preparerForBotsListSecrets(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsListSecrets", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsListSecrets", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForBotsListSecrets(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsListSecrets", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForBotsListSecrets prepares the BotsListSecrets request.
func (c HealthbotsClient) preparerForBotsListSecrets(ctx context.Context, id HealthBotId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listSecrets", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForBotsListSecrets handles the response to the BotsListSecrets request. The method always
// closes the http.Response Body.
func (c HealthbotsClient) responderForBotsListSecrets(resp *http.Response) (result BotsListSecretsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
