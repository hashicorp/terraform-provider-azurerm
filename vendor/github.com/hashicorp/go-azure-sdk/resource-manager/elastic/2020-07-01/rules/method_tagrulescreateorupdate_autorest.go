package rules

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagRulesCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *MonitoringTagRules
}

// TagRulesCreateOrUpdate ...
func (c RulesClient) TagRulesCreateOrUpdate(ctx context.Context, id TagRuleId, input MonitoringTagRules) (result TagRulesCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForTagRulesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "rules.RulesClient", "TagRulesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "rules.RulesClient", "TagRulesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTagRulesCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "rules.RulesClient", "TagRulesCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTagRulesCreateOrUpdate prepares the TagRulesCreateOrUpdate request.
func (c RulesClient) preparerForTagRulesCreateOrUpdate(ctx context.Context, id TagRuleId, input MonitoringTagRules) (*http.Request, error) {
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

// responderForTagRulesCreateOrUpdate handles the response to the TagRulesCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c RulesClient) responderForTagRulesCreateOrUpdate(resp *http.Response) (result TagRulesCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
