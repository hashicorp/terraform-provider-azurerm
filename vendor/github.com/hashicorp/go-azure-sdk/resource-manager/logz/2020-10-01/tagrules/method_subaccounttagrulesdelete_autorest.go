package tagrules

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubAccountTagRulesDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// SubAccountTagRulesDelete ...
func (c TagRulesClient) SubAccountTagRulesDelete(ctx context.Context, id AccountTagRuleId) (result SubAccountTagRulesDeleteOperationResponse, err error) {
	req, err := c.preparerForSubAccountTagRulesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tagrules.TagRulesClient", "SubAccountTagRulesDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "tagrules.TagRulesClient", "SubAccountTagRulesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSubAccountTagRulesDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tagrules.TagRulesClient", "SubAccountTagRulesDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSubAccountTagRulesDelete prepares the SubAccountTagRulesDelete request.
func (c TagRulesClient) preparerForSubAccountTagRulesDelete(ctx context.Context, id AccountTagRuleId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSubAccountTagRulesDelete handles the response to the SubAccountTagRulesDelete request. The method always
// closes the http.Response Body.
func (c TagRulesClient) responderForSubAccountTagRulesDelete(resp *http.Response) (result SubAccountTagRulesDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusAccepted, http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
