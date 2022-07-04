package monitorsresource

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *ElasticMonitorResource
}

// MonitorsGet ...
func (c MonitorsResourceClient) MonitorsGet(ctx context.Context, id MonitorId) (result MonitorsGetOperationResponse, err error) {
	req, err := c.preparerForMonitorsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMonitorsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMonitorsGet prepares the MonitorsGet request.
func (c MonitorsResourceClient) preparerForMonitorsGet(ctx context.Context, id MonitorId) (*http.Request, error) {
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

// responderForMonitorsGet handles the response to the MonitorsGet request. The method always
// closes the http.Response Body.
func (c MonitorsResourceClient) responderForMonitorsGet(resp *http.Response) (result MonitorsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
