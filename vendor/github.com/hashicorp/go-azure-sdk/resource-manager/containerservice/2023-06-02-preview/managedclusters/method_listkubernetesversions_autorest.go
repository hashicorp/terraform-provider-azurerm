package managedclusters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListKubernetesVersionsOperationResponse struct {
	HttpResponse *http.Response
	Model        *KubernetesVersionListResult
}

// ListKubernetesVersions ...
func (c ManagedClustersClient) ListKubernetesVersions(ctx context.Context, id LocationId) (result ListKubernetesVersionsOperationResponse, err error) {
	req, err := c.preparerForListKubernetesVersions(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ListKubernetesVersions", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ListKubernetesVersions", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListKubernetesVersions(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ListKubernetesVersions", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListKubernetesVersions prepares the ListKubernetesVersions request.
func (c ManagedClustersClient) preparerForListKubernetesVersions(ctx context.Context, id LocationId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/kubernetesVersions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListKubernetesVersions handles the response to the ListKubernetesVersions request. The method always
// closes the http.Response Body.
func (c ManagedClustersClient) responderForListKubernetesVersions(resp *http.Response) (result ListKubernetesVersionsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
