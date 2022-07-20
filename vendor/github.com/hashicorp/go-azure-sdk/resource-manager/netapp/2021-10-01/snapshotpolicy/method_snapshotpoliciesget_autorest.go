package snapshotpolicy

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SnapshotPoliciesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *SnapshotPolicy
}

// SnapshotPoliciesGet ...
func (c SnapshotPolicyClient) SnapshotPoliciesGet(ctx context.Context, id SnapshotPoliciesId) (result SnapshotPoliciesGetOperationResponse, err error) {
	req, err := c.preparerForSnapshotPoliciesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "snapshotpolicy.SnapshotPolicyClient", "SnapshotPoliciesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "snapshotpolicy.SnapshotPolicyClient", "SnapshotPoliciesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSnapshotPoliciesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "snapshotpolicy.SnapshotPolicyClient", "SnapshotPoliciesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSnapshotPoliciesGet prepares the SnapshotPoliciesGet request.
func (c SnapshotPolicyClient) preparerForSnapshotPoliciesGet(ctx context.Context, id SnapshotPoliciesId) (*http.Request, error) {
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

// responderForSnapshotPoliciesGet handles the response to the SnapshotPoliciesGet request. The method always
// closes the http.Response Body.
func (c SnapshotPolicyClient) responderForSnapshotPoliciesGet(resp *http.Response) (result SnapshotPoliciesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
