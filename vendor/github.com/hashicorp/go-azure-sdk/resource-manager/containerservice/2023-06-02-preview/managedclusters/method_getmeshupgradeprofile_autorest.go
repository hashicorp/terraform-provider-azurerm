package managedclusters

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetMeshUpgradeProfileOperationResponse struct {
	HttpResponse *http.Response
	Model        *MeshUpgradeProfile
}

// GetMeshUpgradeProfile ...
func (c ManagedClustersClient) GetMeshUpgradeProfile(ctx context.Context, id MeshUpgradeProfileId) (result GetMeshUpgradeProfileOperationResponse, err error) {
	req, err := c.preparerForGetMeshUpgradeProfile(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "GetMeshUpgradeProfile", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "GetMeshUpgradeProfile", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetMeshUpgradeProfile(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "GetMeshUpgradeProfile", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetMeshUpgradeProfile prepares the GetMeshUpgradeProfile request.
func (c ManagedClustersClient) preparerForGetMeshUpgradeProfile(ctx context.Context, id MeshUpgradeProfileId) (*http.Request, error) {
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

// responderForGetMeshUpgradeProfile handles the response to the GetMeshUpgradeProfile request. The method always
// closes the http.Response Body.
func (c ManagedClustersClient) responderForGetMeshUpgradeProfile(resp *http.Response) (result GetMeshUpgradeProfileOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
