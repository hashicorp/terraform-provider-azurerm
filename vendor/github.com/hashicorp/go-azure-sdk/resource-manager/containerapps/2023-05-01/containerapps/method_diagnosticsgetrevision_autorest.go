package containerapps

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticsGetRevisionOperationResponse struct {
	HttpResponse *http.Response
	Model        *Revision
}

// DiagnosticsGetRevision ...
func (c ContainerAppsClient) DiagnosticsGetRevision(ctx context.Context, id RevisionsApiRevisionId) (result DiagnosticsGetRevisionOperationResponse, err error) {
	req, err := c.preparerForDiagnosticsGetRevision(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerapps.ContainerAppsClient", "DiagnosticsGetRevision", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerapps.ContainerAppsClient", "DiagnosticsGetRevision", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDiagnosticsGetRevision(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerapps.ContainerAppsClient", "DiagnosticsGetRevision", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDiagnosticsGetRevision prepares the DiagnosticsGetRevision request.
func (c ContainerAppsClient) preparerForDiagnosticsGetRevision(ctx context.Context, id RevisionsApiRevisionId) (*http.Request, error) {
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

// responderForDiagnosticsGetRevision handles the response to the DiagnosticsGetRevision request. The method always
// closes the http.Response Body.
func (c ContainerAppsClient) responderForDiagnosticsGetRevision(resp *http.Response) (result DiagnosticsGetRevisionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
