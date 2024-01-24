package managedenvironments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedEnvironmentDiagnosticsListDetectorsOperationResponse struct {
	HttpResponse *http.Response
	Model        *DiagnosticsCollection
}

// ManagedEnvironmentDiagnosticsListDetectors ...
func (c ManagedEnvironmentsClient) ManagedEnvironmentDiagnosticsListDetectors(ctx context.Context, id ManagedEnvironmentId) (result ManagedEnvironmentDiagnosticsListDetectorsOperationResponse, err error) {
	req, err := c.preparerForManagedEnvironmentDiagnosticsListDetectors(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedEnvironmentDiagnosticsListDetectors", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedEnvironmentDiagnosticsListDetectors", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForManagedEnvironmentDiagnosticsListDetectors(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedEnvironmentDiagnosticsListDetectors", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForManagedEnvironmentDiagnosticsListDetectors prepares the ManagedEnvironmentDiagnosticsListDetectors request.
func (c ManagedEnvironmentsClient) preparerForManagedEnvironmentDiagnosticsListDetectors(ctx context.Context, id ManagedEnvironmentId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/detectors", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForManagedEnvironmentDiagnosticsListDetectors handles the response to the ManagedEnvironmentDiagnosticsListDetectors request. The method always
// closes the http.Response Body.
func (c ManagedEnvironmentsClient) responderForManagedEnvironmentDiagnosticsListDetectors(resp *http.Response) (result ManagedEnvironmentDiagnosticsListDetectorsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
