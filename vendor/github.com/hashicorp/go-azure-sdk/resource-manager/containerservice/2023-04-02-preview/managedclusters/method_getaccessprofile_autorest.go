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

type GetAccessProfileOperationResponse struct {
	HttpResponse *http.Response
	Model        *ManagedClusterAccessProfile
}

// GetAccessProfile ...
func (c ManagedClustersClient) GetAccessProfile(ctx context.Context, id AccessProfileId) (result GetAccessProfileOperationResponse, err error) {
	req, err := c.preparerForGetAccessProfile(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "GetAccessProfile", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "GetAccessProfile", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetAccessProfile(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "GetAccessProfile", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetAccessProfile prepares the GetAccessProfile request.
func (c ManagedClustersClient) preparerForGetAccessProfile(ctx context.Context, id AccessProfileId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listCredential", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetAccessProfile handles the response to the GetAccessProfile request. The method always
// closes the http.Response Body.
func (c ManagedClustersClient) responderForGetAccessProfile(resp *http.Response) (result GetAccessProfileOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
