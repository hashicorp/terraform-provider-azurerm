package alertruletemplates

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertRuleTemplatesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *AlertRuleTemplate
}

// AlertRuleTemplatesGet ...
func (c AlertRuleTemplatesClient) AlertRuleTemplatesGet(ctx context.Context, id AlertRuleTemplateId) (result AlertRuleTemplatesGetOperationResponse, err error) {
	req, err := c.preparerForAlertRuleTemplatesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertruletemplates.AlertRuleTemplatesClient", "AlertRuleTemplatesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertruletemplates.AlertRuleTemplatesClient", "AlertRuleTemplatesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAlertRuleTemplatesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertruletemplates.AlertRuleTemplatesClient", "AlertRuleTemplatesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAlertRuleTemplatesGet prepares the AlertRuleTemplatesGet request.
func (c AlertRuleTemplatesClient) preparerForAlertRuleTemplatesGet(ctx context.Context, id AlertRuleTemplateId) (*http.Request, error) {
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

// responderForAlertRuleTemplatesGet handles the response to the AlertRuleTemplatesGet request. The method always
// closes the http.Response Body.
func (c AlertRuleTemplatesClient) responderForAlertRuleTemplatesGet(resp *http.Response) (result AlertRuleTemplatesGetOperationResponse, err error) {
	var respObj json.RawMessage
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&respObj),
		autorest.ByClosing())
	result.HttpResponse = resp
	if err != nil {
		return
	}
	model, err := unmarshalAlertRuleTemplateImplementation(respObj)
	if err != nil {
		return
	}
	result.Model = &model
	return
}
