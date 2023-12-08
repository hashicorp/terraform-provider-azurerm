package containerapps

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticsGetDetectorOperationResponse struct {
	HttpResponse *http.Response
	Model        *Diagnostics
}

// DiagnosticsGetDetector ...
func (c ContainerAppsClient) DiagnosticsGetDetector(ctx context.Context, id DetectorId) (result DiagnosticsGetDetectorOperationResponse, err error) {
	req, err := c.preparerForDiagnosticsGetDetector(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerapps.ContainerAppsClient", "DiagnosticsGetDetector", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerapps.ContainerAppsClient", "DiagnosticsGetDetector", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDiagnosticsGetDetector(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerapps.ContainerAppsClient", "DiagnosticsGetDetector", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDiagnosticsGetDetector prepares the DiagnosticsGetDetector request.
func (c ContainerAppsClient) preparerForDiagnosticsGetDetector(ctx context.Context, id DetectorId) (*http.Request, error) {
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

// responderForDiagnosticsGetDetector handles the response to the DiagnosticsGetDetector request. The method always
// closes the http.Response Body.
func (c ContainerAppsClient) responderForDiagnosticsGetDetector(resp *http.Response) (result DiagnosticsGetDetectorOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
