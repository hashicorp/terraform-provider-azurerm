package managedcassandras

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CassandraClustersGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *ClusterResource
}

// CassandraClustersGet ...
func (c ManagedCassandrasClient) CassandraClustersGet(ctx context.Context, id CassandraClusterId) (result CassandraClustersGetOperationResponse, err error) {
	req, err := c.preparerForCassandraClustersGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCassandraClustersGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCassandraClustersGet prepares the CassandraClustersGet request.
func (c ManagedCassandrasClient) preparerForCassandraClustersGet(ctx context.Context, id CassandraClusterId) (*http.Request, error) {
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

// responderForCassandraClustersGet handles the response to the CassandraClustersGet request. The method always
// closes the http.Response Body.
func (c ManagedCassandrasClient) responderForCassandraClustersGet(resp *http.Response) (result CassandraClustersGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
