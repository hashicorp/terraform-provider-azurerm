package grafanaresource

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GrafanaGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *ManagedGrafana
}

// GrafanaGet ...
func (c GrafanaResourceClient) GrafanaGet(ctx context.Context, id GrafanaId) (result GrafanaGetOperationResponse, err error) {
	req, err := c.preparerForGrafanaGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "grafanaresource.GrafanaResourceClient", "GrafanaGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "grafanaresource.GrafanaResourceClient", "GrafanaGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGrafanaGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "grafanaresource.GrafanaResourceClient", "GrafanaGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGrafanaGet prepares the GrafanaGet request.
func (c GrafanaResourceClient) preparerForGrafanaGet(ctx context.Context, id GrafanaId) (*http.Request, error) {
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

// responderForGrafanaGet handles the response to the GrafanaGet request. The method always
// closes the http.Response Body.
func (c GrafanaResourceClient) responderForGrafanaGet(resp *http.Response) (result GrafanaGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
