package managedcassandras

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CassandraDataCentersListOperationResponse struct {
	HttpResponse *http.Response
	Model        *ListDataCenters
}

// CassandraDataCentersList ...
func (c ManagedCassandrasClient) CassandraDataCentersList(ctx context.Context, id CassandraClusterId) (result CassandraDataCentersListOperationResponse, err error) {
	req, err := c.preparerForCassandraDataCentersList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraDataCentersList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraDataCentersList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCassandraDataCentersList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraDataCentersList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCassandraDataCentersList prepares the CassandraDataCentersList request.
func (c ManagedCassandrasClient) preparerForCassandraDataCentersList(ctx context.Context, id CassandraClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/dataCenters", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCassandraDataCentersList handles the response to the CassandraDataCentersList request. The method always
// closes the http.Response Body.
func (c ManagedCassandrasClient) responderForCassandraDataCentersList(resp *http.Response) (result CassandraDataCentersListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
