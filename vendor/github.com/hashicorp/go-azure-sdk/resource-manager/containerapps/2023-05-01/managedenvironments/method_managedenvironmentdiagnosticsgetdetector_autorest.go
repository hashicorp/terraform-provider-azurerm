package managedenvironments

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedEnvironmentDiagnosticsGetDetectorOperationResponse struct {
	HttpResponse *http.Response
	Model        *Diagnostics
}

// ManagedEnvironmentDiagnosticsGetDetector ...
func (c ManagedEnvironmentsClient) ManagedEnvironmentDiagnosticsGetDetector(ctx context.Context, id ManagedEnvironmentDetectorId) (result ManagedEnvironmentDiagnosticsGetDetectorOperationResponse, err error) {
	req, err := c.preparerForManagedEnvironmentDiagnosticsGetDetector(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedEnvironmentDiagnosticsGetDetector", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedEnvironmentDiagnosticsGetDetector", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForManagedEnvironmentDiagnosticsGetDetector(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedEnvironmentDiagnosticsGetDetector", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForManagedEnvironmentDiagnosticsGetDetector prepares the ManagedEnvironmentDiagnosticsGetDetector request.
func (c ManagedEnvironmentsClient) preparerForManagedEnvironmentDiagnosticsGetDetector(ctx context.Context, id ManagedEnvironmentDetectorId) (*http.Request, error) {
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

// responderForManagedEnvironmentDiagnosticsGetDetector handles the response to the ManagedEnvironmentDiagnosticsGetDetector request. The method always
// closes the http.Response Body.
func (c ManagedEnvironmentsClient) responderForManagedEnvironmentDiagnosticsGetDetector(resp *http.Response) (result ManagedEnvironmentDiagnosticsGetDetectorOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
