package diagnosticsettingscategories

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticSettingsCategoryListOperationResponse struct {
	HttpResponse *http.Response
	Model        *DiagnosticSettingsCategoryResourceCollection
}

// DiagnosticSettingsCategoryList ...
func (c DiagnosticSettingsCategoriesClient) DiagnosticSettingsCategoryList(ctx context.Context, id commonids.ScopeId) (result DiagnosticSettingsCategoryListOperationResponse, err error) {
	req, err := c.preparerForDiagnosticSettingsCategoryList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "diagnosticsettingscategories.DiagnosticSettingsCategoriesClient", "DiagnosticSettingsCategoryList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "diagnosticsettingscategories.DiagnosticSettingsCategoriesClient", "DiagnosticSettingsCategoryList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDiagnosticSettingsCategoryList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "diagnosticsettingscategories.DiagnosticSettingsCategoriesClient", "DiagnosticSettingsCategoryList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDiagnosticSettingsCategoryList prepares the DiagnosticSettingsCategoryList request.
func (c DiagnosticSettingsCategoriesClient) preparerForDiagnosticSettingsCategoryList(ctx context.Context, id commonids.ScopeId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Insights/diagnosticSettingsCategories", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDiagnosticSettingsCategoryList handles the response to the DiagnosticSettingsCategoryList request. The method always
// closes the http.Response Body.
func (c DiagnosticSettingsCategoriesClient) responderForDiagnosticSettingsCategoryList(resp *http.Response) (result DiagnosticSettingsCategoryListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
