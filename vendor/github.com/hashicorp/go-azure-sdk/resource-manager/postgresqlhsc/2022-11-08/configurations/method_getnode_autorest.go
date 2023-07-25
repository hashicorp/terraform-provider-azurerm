package configurations

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetNodeOperationResponse struct {
	HttpResponse *http.Response
	Model        *ServerConfiguration
}

// GetNode ...
func (c ConfigurationsClient) GetNode(ctx context.Context, id NodeConfigurationId) (result GetNodeOperationResponse, err error) {
	req, err := c.preparerForGetNode(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "GetNode", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "GetNode", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetNode(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "GetNode", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetNode prepares the GetNode request.
func (c ConfigurationsClient) preparerForGetNode(ctx context.Context, id NodeConfigurationId) (*http.Request, error) {
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

// responderForGetNode handles the response to the GetNode request. The method always
// closes the http.Response Body.
func (c ConfigurationsClient) responderForGetNode(resp *http.Response) (result GetNodeOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
