package tagrules

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubAccountTagRulesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *MonitoringTagRules
}

// SubAccountTagRulesGet ...
func (c TagRulesClient) SubAccountTagRulesGet(ctx context.Context, id AccountTagRuleId) (result SubAccountTagRulesGetOperationResponse, err error) {
	req, err := c.preparerForSubAccountTagRulesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tagrules.TagRulesClient", "SubAccountTagRulesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "tagrules.TagRulesClient", "SubAccountTagRulesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSubAccountTagRulesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tagrules.TagRulesClient", "SubAccountTagRulesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSubAccountTagRulesGet prepares the SubAccountTagRulesGet request.
func (c TagRulesClient) preparerForSubAccountTagRulesGet(ctx context.Context, id AccountTagRuleId) (*http.Request, error) {
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

// responderForSubAccountTagRulesGet handles the response to the SubAccountTagRulesGet request. The method always
// closes the http.Response Body.
func (c TagRulesClient) responderForSubAccountTagRulesGet(resp *http.Response) (result SubAccountTagRulesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
