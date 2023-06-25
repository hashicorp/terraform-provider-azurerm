package configurations

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetCoordinatorOperationResponse struct {
	HttpResponse *http.Response
	Model        *ServerConfiguration
}

// GetCoordinator ...
func (c ConfigurationsClient) GetCoordinator(ctx context.Context, id CoordinatorConfigurationId) (result GetCoordinatorOperationResponse, err error) {
	req, err := c.preparerForGetCoordinator(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "GetCoordinator", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "GetCoordinator", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetCoordinator(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "GetCoordinator", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetCoordinator prepares the GetCoordinator request.
func (c ConfigurationsClient) preparerForGetCoordinator(ctx context.Context, id CoordinatorConfigurationId) (*http.Request, error) {
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

// responderForGetCoordinator handles the response to the GetCoordinator request. The method always
// closes the http.Response Body.
func (c ConfigurationsClient) responderForGetCoordinator(resp *http.Response) (result GetCoordinatorOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
