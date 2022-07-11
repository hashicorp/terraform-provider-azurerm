package snapshotpolicy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SnapshotPoliciesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *SnapshotPoliciesList
}

// SnapshotPoliciesList ...
func (c SnapshotPolicyClient) SnapshotPoliciesList(ctx context.Context, id NetAppAccountId) (result SnapshotPoliciesListOperationResponse, err error) {
	req, err := c.preparerForSnapshotPoliciesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "snapshotpolicy.SnapshotPolicyClient", "SnapshotPoliciesList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "snapshotpolicy.SnapshotPolicyClient", "SnapshotPoliciesList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSnapshotPoliciesList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "snapshotpolicy.SnapshotPolicyClient", "SnapshotPoliciesList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSnapshotPoliciesList prepares the SnapshotPoliciesList request.
func (c SnapshotPolicyClient) preparerForSnapshotPoliciesList(ctx context.Context, id NetAppAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/snapshotPolicies", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSnapshotPoliciesList handles the response to the SnapshotPoliciesList request. The method always
// closes the http.Response Body.
func (c SnapshotPolicyClient) responderForSnapshotPoliciesList(resp *http.Response) (result SnapshotPoliciesListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
