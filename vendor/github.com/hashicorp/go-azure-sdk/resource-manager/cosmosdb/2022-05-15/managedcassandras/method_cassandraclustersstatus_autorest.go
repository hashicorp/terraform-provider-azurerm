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

type CassandraClustersStatusOperationResponse struct {
	HttpResponse *http.Response
	Model        *CassandraClusterPublicStatus
}

// CassandraClustersStatus ...
func (c ManagedCassandrasClient) CassandraClustersStatus(ctx context.Context, id CassandraClusterId) (result CassandraClustersStatusOperationResponse, err error) {
	req, err := c.preparerForCassandraClustersStatus(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersStatus", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersStatus", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCassandraClustersStatus(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersStatus", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCassandraClustersStatus prepares the CassandraClustersStatus request.
func (c ManagedCassandrasClient) preparerForCassandraClustersStatus(ctx context.Context, id CassandraClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/status", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCassandraClustersStatus handles the response to the CassandraClustersStatus request. The method always
// closes the http.Response Body.
func (c ManagedCassandrasClient) responderForCassandraClustersStatus(resp *http.Response) (result CassandraClustersStatusOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
