package subscriptions

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RulesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *Rule
}

// RulesGet ...
func (c SubscriptionsClient) RulesGet(ctx context.Context, id RuleId) (result RulesGetOperationResponse, err error) {
	req, err := c.preparerForRulesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "RulesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "RulesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRulesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "RulesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRulesGet prepares the RulesGet request.
func (c SubscriptionsClient) preparerForRulesGet(ctx context.Context, id RuleId) (*http.Request, error) {
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

// responderForRulesGet handles the response to the RulesGet request. The method always
// closes the http.Response Body.
func (c SubscriptionsClient) responderForRulesGet(resp *http.Response) (result RulesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
