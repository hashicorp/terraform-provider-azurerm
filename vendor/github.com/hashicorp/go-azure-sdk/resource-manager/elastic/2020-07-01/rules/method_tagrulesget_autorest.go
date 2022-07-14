package rules

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagRulesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *MonitoringTagRules
}

// TagRulesGet ...
func (c RulesClient) TagRulesGet(ctx context.Context, id TagRuleId) (result TagRulesGetOperationResponse, err error) {
	req, err := c.preparerForTagRulesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "rules.RulesClient", "TagRulesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "rules.RulesClient", "TagRulesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTagRulesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "rules.RulesClient", "TagRulesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTagRulesGet prepares the TagRulesGet request.
func (c RulesClient) preparerForTagRulesGet(ctx context.Context, id TagRuleId) (*http.Request, error) {
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

// responderForTagRulesGet handles the response to the TagRulesGet request. The method always
// closes the http.Response Body.
func (c RulesClient) responderForTagRulesGet(resp *http.Response) (result TagRulesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
