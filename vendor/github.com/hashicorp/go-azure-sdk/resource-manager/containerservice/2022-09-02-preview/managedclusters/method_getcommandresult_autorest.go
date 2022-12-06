package managedclusters

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetCommandResultOperationResponse struct {
	HttpResponse *http.Response
	Model        *RunCommandResult
}

// GetCommandResult ...
func (c ManagedClustersClient) GetCommandResult(ctx context.Context, id CommandResultId) (result GetCommandResultOperationResponse, err error) {
	req, err := c.preparerForGetCommandResult(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "GetCommandResult", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "GetCommandResult", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetCommandResult(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "GetCommandResult", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetCommandResult prepares the GetCommandResult request.
func (c ManagedClustersClient) preparerForGetCommandResult(ctx context.Context, id CommandResultId) (*http.Request, error) {
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

// responderForGetCommandResult handles the response to the GetCommandResult request. The method always
// closes the http.Response Body.
func (c ManagedClustersClient) responderForGetCommandResult(resp *http.Response) (result GetCommandResultOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusAccepted, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
