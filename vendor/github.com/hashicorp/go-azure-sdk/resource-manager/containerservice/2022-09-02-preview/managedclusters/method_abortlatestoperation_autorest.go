package managedclusters

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

type AbortLatestOperationOperationResponse struct {
	HttpResponse *http.Response
}

// AbortLatestOperation ...
func (c ManagedClustersClient) AbortLatestOperation(ctx context.Context, id commonids.KubernetesClusterId) (result AbortLatestOperationOperationResponse, err error) {
	req, err := c.preparerForAbortLatestOperation(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "AbortLatestOperation", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "AbortLatestOperation", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAbortLatestOperation(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "AbortLatestOperation", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAbortLatestOperation prepares the AbortLatestOperation request.
func (c ManagedClustersClient) preparerForAbortLatestOperation(ctx context.Context, id commonids.KubernetesClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/abort", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAbortLatestOperation handles the response to the AbortLatestOperation request. The method always
// closes the http.Response Body.
func (c ManagedClustersClient) responderForAbortLatestOperation(resp *http.Response) (result AbortLatestOperationOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
