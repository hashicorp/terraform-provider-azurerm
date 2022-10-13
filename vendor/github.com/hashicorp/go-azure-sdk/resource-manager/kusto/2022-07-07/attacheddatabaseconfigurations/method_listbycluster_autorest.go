package attacheddatabaseconfigurations

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByClusterOperationResponse struct {
	HttpResponse *http.Response
	Model        *AttachedDatabaseConfigurationListResult
}

// ListByCluster ...
func (c AttachedDatabaseConfigurationsClient) ListByCluster(ctx context.Context, id ClusterId) (result ListByClusterOperationResponse, err error) {
	req, err := c.preparerForListByCluster(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "attacheddatabaseconfigurations.AttachedDatabaseConfigurationsClient", "ListByCluster", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "attacheddatabaseconfigurations.AttachedDatabaseConfigurationsClient", "ListByCluster", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListByCluster(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "attacheddatabaseconfigurations.AttachedDatabaseConfigurationsClient", "ListByCluster", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListByCluster prepares the ListByCluster request.
func (c AttachedDatabaseConfigurationsClient) preparerForListByCluster(ctx context.Context, id ClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/attachedDatabaseConfigurations", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListByCluster handles the response to the ListByCluster request. The method always
// closes the http.Response Body.
func (c AttachedDatabaseConfigurationsClient) responderForListByCluster(resp *http.Response) (result ListByClusterOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
