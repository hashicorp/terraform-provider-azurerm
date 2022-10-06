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

type CassandraClustersListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *ListClusters
}

// CassandraClustersListByResourceGroup ...
func (c ManagedCassandrasClient) CassandraClustersListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result CassandraClustersListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForCassandraClustersListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersListByResourceGroup", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersListByResourceGroup", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCassandraClustersListByResourceGroup(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersListByResourceGroup", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCassandraClustersListByResourceGroup prepares the CassandraClustersListByResourceGroup request.
func (c ManagedCassandrasClient) preparerForCassandraClustersListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
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

// responderForCassandraClustersListByResourceGroup handles the response to the CassandraClustersListByResourceGroup request. The method always
// closes the http.Response Body.
func (c ManagedCassandrasClient) responderForCassandraClustersListByResourceGroup(resp *http.Response) (result CassandraClustersListByResourceGroupOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
