package managedclusters

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetMeshRevisionProfileOperationResponse struct {
	HttpResponse *http.Response
	Model        *MeshRevisionProfile
}

// GetMeshRevisionProfile ...
func (c ManagedClustersClient) GetMeshRevisionProfile(ctx context.Context, id MeshRevisionProfileId) (result GetMeshRevisionProfileOperationResponse, err error) {
	req, err := c.preparerForGetMeshRevisionProfile(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "GetMeshRevisionProfile", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "GetMeshRevisionProfile", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetMeshRevisionProfile(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "GetMeshRevisionProfile", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetMeshRevisionProfile prepares the GetMeshRevisionProfile request.
func (c ManagedClustersClient) preparerForGetMeshRevisionProfile(ctx context.Context, id MeshRevisionProfileId) (*http.Request, error) {
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

// responderForGetMeshRevisionProfile handles the response to the GetMeshRevisionProfile request. The method always
// closes the http.Response Body.
func (c ManagedClustersClient) responderForGetMeshRevisionProfile(resp *http.Response) (result GetMeshRevisionProfileOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
