package diagnosticsettingscategories

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticSettingsCategoryGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *DiagnosticSettingsCategoryResource
}

// DiagnosticSettingsCategoryGet ...
func (c DiagnosticSettingsCategoriesClient) DiagnosticSettingsCategoryGet(ctx context.Context, id ScopedDiagnosticSettingsCategoryId) (result DiagnosticSettingsCategoryGetOperationResponse, err error) {
	req, err := c.preparerForDiagnosticSettingsCategoryGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "diagnosticsettingscategories.DiagnosticSettingsCategoriesClient", "DiagnosticSettingsCategoryGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "diagnosticsettingscategories.DiagnosticSettingsCategoriesClient", "DiagnosticSettingsCategoryGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDiagnosticSettingsCategoryGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "diagnosticsettingscategories.DiagnosticSettingsCategoriesClient", "DiagnosticSettingsCategoryGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDiagnosticSettingsCategoryGet prepares the DiagnosticSettingsCategoryGet request.
func (c DiagnosticSettingsCategoriesClient) preparerForDiagnosticSettingsCategoryGet(ctx context.Context, id ScopedDiagnosticSettingsCategoryId) (*http.Request, error) {
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

// responderForDiagnosticSettingsCategoryGet handles the response to the DiagnosticSettingsCategoryGet request. The method always
// closes the http.Response Body.
func (c DiagnosticSettingsCategoriesClient) responderForDiagnosticSettingsCategoryGet(resp *http.Response) (result DiagnosticSettingsCategoryGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
