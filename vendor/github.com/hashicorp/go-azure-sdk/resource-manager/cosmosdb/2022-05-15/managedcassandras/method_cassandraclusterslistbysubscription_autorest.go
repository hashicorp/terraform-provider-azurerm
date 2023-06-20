package managedcassandras

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CassandraClustersListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *ListClusters
}

// CassandraClustersListBySubscription ...
func (c ManagedCassandrasClient) CassandraClustersListBySubscription(ctx context.Context, id commonids.SubscriptionId) (result CassandraClustersListBySubscriptionOperationResponse, err error) {
	req, err := c.preparerForCassandraClustersListBySubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersListBySubscription", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersListBySubscription", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCassandraClustersListBySubscription(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersListBySubscription", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCassandraClustersListBySubscription prepares the CassandraClustersListBySubscription request.
func (c ManagedCassandrasClient) preparerForCassandraClustersListBySubscription(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.DocumentDB/cassandraClusters", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCassandraClustersListBySubscription handles the response to the CassandraClustersListBySubscription request. The method always
// closes the http.Response Body.
func (c ManagedCassandrasClient) responderForCassandraClustersListBySubscription(resp *http.Response) (result CassandraClustersListBySubscriptionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
