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

type SecurityMLAnalyticsSettingsCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *SecurityMLAnalyticsSetting
}

// SecurityMLAnalyticsSettingsCreateOrUpdate ...
func (c SecurityMLAnalyticsSettingsClient) SecurityMLAnalyticsSettingsCreateOrUpdate(ctx context.Context, id SecurityMLAnalyticsSettingId, input SecurityMLAnalyticsSetting) (result SecurityMLAnalyticsSettingsCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForSecurityMLAnalyticsSettingsCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securitymlanalyticssettings.SecurityMLAnalyticsSettingsClient", "SecurityMLAnalyticsSettingsCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "securitymlanalyticssettings.SecurityMLAnalyticsSettingsClient", "SecurityMLAnalyticsSettingsCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSecurityMLAnalyticsSettingsCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securitymlanalyticssettings.SecurityMLAnalyticsSettingsClient", "SecurityMLAnalyticsSettingsCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSecurityMLAnalyticsSettingsCreateOrUpdate prepares the SecurityMLAnalyticsSettingsCreateOrUpdate request.
func (c SecurityMLAnalyticsSettingsClient) preparerForSecurityMLAnalyticsSettingsCreateOrUpdate(ctx context.Context, id SecurityMLAnalyticsSettingId, input SecurityMLAnalyticsSetting) (*http.Request, error) {
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

// responderForSecurityMLAnalyticsSettingsCreateOrUpdate handles the response to the SecurityMLAnalyticsSettingsCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c SecurityMLAnalyticsSettingsClient) responderForSecurityMLAnalyticsSettingsCreateOrUpdate(resp *http.Response) (result SecurityMLAnalyticsSettingsCreateOrUpdateOperationResponse, err error) {
	var respObj json.RawMessage
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
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
