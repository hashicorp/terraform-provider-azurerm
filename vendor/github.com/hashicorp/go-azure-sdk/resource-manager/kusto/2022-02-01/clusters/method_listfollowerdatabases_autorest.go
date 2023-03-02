package clusters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListFollowerDatabasesOperationResponse struct {
	HttpResponse *http.Response
	Model        *FollowerDatabaseListResult
}

// ListFollowerDatabases ...
func (c ClustersClient) ListFollowerDatabases(ctx context.Context, id ClusterId) (result ListFollowerDatabasesOperationResponse, err error) {
	req, err := c.preparerForListFollowerDatabases(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "ListFollowerDatabases", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "ListFollowerDatabases", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListFollowerDatabases(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "ListFollowerDatabases", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListFollowerDatabases prepares the ListFollowerDatabases request.
func (c ClustersClient) preparerForListFollowerDatabases(ctx context.Context, id ClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listFollowerDatabases", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListFollowerDatabases handles the response to the ListFollowerDatabases request. The method always
// closes the http.Response Body.
func (c ClustersClient) responderForListFollowerDatabases(resp *http.Response) (result ListFollowerDatabasesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
