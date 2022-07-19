package snapshotpolicy

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SnapshotPoliciesCreateOperationResponse struct {
	HttpResponse *http.Response
	Model        *SnapshotPolicy
}

// SnapshotPoliciesCreate ...
func (c SnapshotPolicyClient) SnapshotPoliciesCreate(ctx context.Context, id SnapshotPoliciesId, input SnapshotPolicy) (result SnapshotPoliciesCreateOperationResponse, err error) {
	req, err := c.preparerForSnapshotPoliciesCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "snapshotpolicy.SnapshotPolicyClient", "SnapshotPoliciesCreate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "snapshotpolicy.SnapshotPolicyClient", "SnapshotPoliciesCreate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSnapshotPoliciesCreate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "snapshotpolicy.SnapshotPolicyClient", "SnapshotPoliciesCreate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSnapshotPoliciesCreate prepares the SnapshotPoliciesCreate request.
func (c SnapshotPolicyClient) preparerForSnapshotPoliciesCreate(ctx context.Context, id SnapshotPoliciesId, input SnapshotPolicy) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSnapshotPoliciesCreate handles the response to the SnapshotPoliciesCreate request. The method always
// closes the http.Response Body.
func (c SnapshotPolicyClient) responderForSnapshotPoliciesCreate(resp *http.Response) (result SnapshotPoliciesCreateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
