package tagrules

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubAccountTagRulesCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *MonitoringTagRules
}

// SubAccountTagRulesCreateOrUpdate ...
func (c TagRulesClient) SubAccountTagRulesCreateOrUpdate(ctx context.Context, id AccountTagRuleId, input MonitoringTagRules) (result SubAccountTagRulesCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForSubAccountTagRulesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tagrules.TagRulesClient", "SubAccountTagRulesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "tagrules.TagRulesClient", "SubAccountTagRulesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSubAccountTagRulesCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tagrules.TagRulesClient", "SubAccountTagRulesCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSubAccountTagRulesCreateOrUpdate prepares the SubAccountTagRulesCreateOrUpdate request.
func (c TagRulesClient) preparerForSubAccountTagRulesCreateOrUpdate(ctx context.Context, id AccountTagRuleId, input MonitoringTagRules) (*http.Request, error) {
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

// responderForSubAccountTagRulesCreateOrUpdate handles the response to the SubAccountTagRulesCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c TagRulesClient) responderForSubAccountTagRulesCreateOrUpdate(resp *http.Response) (result SubAccountTagRulesCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
