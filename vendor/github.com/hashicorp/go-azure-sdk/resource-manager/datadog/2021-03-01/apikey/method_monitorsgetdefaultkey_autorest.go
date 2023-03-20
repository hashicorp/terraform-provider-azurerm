package apikey

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorsGetDefaultKeyOperationResponse struct {
	HttpResponse *http.Response
	Model        *DatadogApiKey
}

// MonitorsGetDefaultKey ...
func (c ApiKeyClient) MonitorsGetDefaultKey(ctx context.Context, id MonitorId) (result MonitorsGetDefaultKeyOperationResponse, err error) {
	req, err := c.preparerForMonitorsGetDefaultKey(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apikey.ApiKeyClient", "MonitorsGetDefaultKey", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "apikey.ApiKeyClient", "MonitorsGetDefaultKey", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMonitorsGetDefaultKey(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apikey.ApiKeyClient", "MonitorsGetDefaultKey", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMonitorsGetDefaultKey prepares the MonitorsGetDefaultKey request.
func (c ApiKeyClient) preparerForMonitorsGetDefaultKey(ctx context.Context, id MonitorId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/getDefaultKey", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForMonitorsGetDefaultKey handles the response to the MonitorsGetDefaultKey request. The method always
// closes the http.Response Body.
func (c ApiKeyClient) responderForMonitorsGetDefaultKey(resp *http.Response) (result MonitorsGetDefaultKeyOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
