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

type MonitorsSetDefaultKeyOperationResponse struct {
	HttpResponse *http.Response
}

// MonitorsSetDefaultKey ...
func (c ApiKeyClient) MonitorsSetDefaultKey(ctx context.Context, id MonitorId, input DatadogApiKey) (result MonitorsSetDefaultKeyOperationResponse, err error) {
	req, err := c.preparerForMonitorsSetDefaultKey(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apikey.ApiKeyClient", "MonitorsSetDefaultKey", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "apikey.ApiKeyClient", "MonitorsSetDefaultKey", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMonitorsSetDefaultKey(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apikey.ApiKeyClient", "MonitorsSetDefaultKey", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMonitorsSetDefaultKey prepares the MonitorsSetDefaultKey request.
func (c ApiKeyClient) preparerForMonitorsSetDefaultKey(ctx context.Context, id MonitorId, input DatadogApiKey) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/setDefaultKey", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForMonitorsSetDefaultKey handles the response to the MonitorsSetDefaultKey request. The method always
// closes the http.Response Body.
func (c ApiKeyClient) responderForMonitorsSetDefaultKey(resp *http.Response) (result MonitorsSetDefaultKeyOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
