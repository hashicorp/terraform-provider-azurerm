package managedcassandras

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CassandraDataCentersGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *DataCenterResource
}

// CassandraDataCentersGet ...
func (c ManagedCassandrasClient) CassandraDataCentersGet(ctx context.Context, id DataCenterId) (result CassandraDataCentersGetOperationResponse, err error) {
	req, err := c.preparerForCassandraDataCentersGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraDataCentersGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraDataCentersGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCassandraDataCentersGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraDataCentersGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCassandraDataCentersGet prepares the CassandraDataCentersGet request.
func (c ManagedCassandrasClient) preparerForCassandraDataCentersGet(ctx context.Context, id DataCenterId) (*http.Request, error) {
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

// responderForCassandraDataCentersGet handles the response to the CassandraDataCentersGet request. The method always
// closes the http.Response Body.
func (c ManagedCassandrasClient) responderForCassandraDataCentersGet(resp *http.Response) (result CassandraDataCentersGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
