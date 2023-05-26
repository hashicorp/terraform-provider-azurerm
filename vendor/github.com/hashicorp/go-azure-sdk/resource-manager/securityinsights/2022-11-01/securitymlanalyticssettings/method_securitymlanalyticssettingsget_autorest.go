package securitymlanalyticssettings

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityMLAnalyticsSettingsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *SecurityMLAnalyticsSetting
}

// SecurityMLAnalyticsSettingsGet ...
func (c SecurityMLAnalyticsSettingsClient) SecurityMLAnalyticsSettingsGet(ctx context.Context, id SecurityMLAnalyticsSettingId) (result SecurityMLAnalyticsSettingsGetOperationResponse, err error) {
	req, err := c.preparerForSecurityMLAnalyticsSettingsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securitymlanalyticssettings.SecurityMLAnalyticsSettingsClient", "SecurityMLAnalyticsSettingsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "securitymlanalyticssettings.SecurityMLAnalyticsSettingsClient", "SecurityMLAnalyticsSettingsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSecurityMLAnalyticsSettingsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securitymlanalyticssettings.SecurityMLAnalyticsSettingsClient", "SecurityMLAnalyticsSettingsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSecurityMLAnalyticsSettingsGet prepares the SecurityMLAnalyticsSettingsGet request.
func (c SecurityMLAnalyticsSettingsClient) preparerForSecurityMLAnalyticsSettingsGet(ctx context.Context, id SecurityMLAnalyticsSettingId) (*http.Request, error) {
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

// responderForSecurityMLAnalyticsSettingsGet handles the response to the SecurityMLAnalyticsSettingsGet request. The method always
// closes the http.Response Body.
func (c SecurityMLAnalyticsSettingsClient) responderForSecurityMLAnalyticsSettingsGet(resp *http.Response) (result SecurityMLAnalyticsSettingsGetOperationResponse, err error) {
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
	model, err := unmarshalSecurityMLAnalyticsSettingImplementation(respObj)
	if err != nil {
		return
	}
	result.Model = &model
	return
}
