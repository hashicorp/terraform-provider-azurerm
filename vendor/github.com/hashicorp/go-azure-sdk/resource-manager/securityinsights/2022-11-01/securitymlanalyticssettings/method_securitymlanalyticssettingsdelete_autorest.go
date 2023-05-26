package securitymlanalyticssettings

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityMLAnalyticsSettingsDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// SecurityMLAnalyticsSettingsDelete ...
func (c SecurityMLAnalyticsSettingsClient) SecurityMLAnalyticsSettingsDelete(ctx context.Context, id SecurityMLAnalyticsSettingId) (result SecurityMLAnalyticsSettingsDeleteOperationResponse, err error) {
	req, err := c.preparerForSecurityMLAnalyticsSettingsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securitymlanalyticssettings.SecurityMLAnalyticsSettingsClient", "SecurityMLAnalyticsSettingsDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "securitymlanalyticssettings.SecurityMLAnalyticsSettingsClient", "SecurityMLAnalyticsSettingsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSecurityMLAnalyticsSettingsDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securitymlanalyticssettings.SecurityMLAnalyticsSettingsClient", "SecurityMLAnalyticsSettingsDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSecurityMLAnalyticsSettingsDelete prepares the SecurityMLAnalyticsSettingsDelete request.
func (c SecurityMLAnalyticsSettingsClient) preparerForSecurityMLAnalyticsSettingsDelete(ctx context.Context, id SecurityMLAnalyticsSettingId) (*http.Request, error) {
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

// responderForSecurityMLAnalyticsSettingsDelete handles the response to the SecurityMLAnalyticsSettingsDelete request. The method always
// closes the http.Response Body.
func (c SecurityMLAnalyticsSettingsClient) responderForSecurityMLAnalyticsSettingsDelete(resp *http.Response) (result SecurityMLAnalyticsSettingsDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
