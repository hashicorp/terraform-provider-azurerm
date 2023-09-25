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

type ListSkusByResourceOperationResponse struct {
	HttpResponse *http.Response
	Model        *ListResourceSkusResult
}

// ListSkusByResource ...
func (c ClustersClient) ListSkusByResource(ctx context.Context, id ClusterId) (result ListSkusByResourceOperationResponse, err error) {
	req, err := c.preparerForListSkusByResource(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "ListSkusByResource", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "ListSkusByResource", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListSkusByResource(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "ListSkusByResource", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListSkusByResource prepares the ListSkusByResource request.
func (c ClustersClient) preparerForListSkusByResource(ctx context.Context, id ClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/skus", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListSkusByResource handles the response to the ListSkusByResource request. The method always
// closes the http.Response Body.
func (c ClustersClient) responderForListSkusByResource(resp *http.Response) (result ListSkusByResourceOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
