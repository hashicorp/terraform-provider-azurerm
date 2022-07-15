package monitorsresource

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorsUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *ElasticMonitorResource
}

// MonitorsUpdate ...
func (c MonitorsResourceClient) MonitorsUpdate(ctx context.Context, id MonitorId, input ElasticMonitorResourceUpdateParameters) (result MonitorsUpdateOperationResponse, err error) {
	req, err := c.preparerForMonitorsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMonitorsUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMonitorsUpdate prepares the MonitorsUpdate request.
func (c MonitorsResourceClient) preparerForMonitorsUpdate(ctx context.Context, id MonitorId, input ElasticMonitorResourceUpdateParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForMonitorsUpdate handles the response to the MonitorsUpdate request. The method always
// closes the http.Response Body.
func (c MonitorsResourceClient) responderForMonitorsUpdate(resp *http.Response) (result MonitorsUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
