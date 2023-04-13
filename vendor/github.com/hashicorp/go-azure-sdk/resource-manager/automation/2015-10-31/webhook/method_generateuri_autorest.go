package webhook

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GenerateUriOperationResponse struct {
	HttpResponse *http.Response
	Model        *string
}

// GenerateUri ...
func (c WebhookClient) GenerateUri(ctx context.Context, id AutomationAccountId) (result GenerateUriOperationResponse, err error) {
	req, err := c.preparerForGenerateUri(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webhook.WebhookClient", "GenerateUri", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "webhook.WebhookClient", "GenerateUri", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGenerateUri(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webhook.WebhookClient", "GenerateUri", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGenerateUri prepares the GenerateUri request.
func (c WebhookClient) preparerForGenerateUri(ctx context.Context, id AutomationAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/webHooks/generateUri", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGenerateUri handles the response to the GenerateUri request. The method always
// closes the http.Response Body.
func (c WebhookClient) responderForGenerateUri(resp *http.Response) (result GenerateUriOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
