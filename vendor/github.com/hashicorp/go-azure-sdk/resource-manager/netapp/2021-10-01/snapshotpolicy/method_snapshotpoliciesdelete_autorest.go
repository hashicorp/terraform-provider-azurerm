package snapshotpolicy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SnapshotPoliciesDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SnapshotPoliciesDelete ...
func (c SnapshotPolicyClient) SnapshotPoliciesDelete(ctx context.Context, id SnapshotPoliciesId) (result SnapshotPoliciesDeleteOperationResponse, err error) {
	req, err := c.preparerForSnapshotPoliciesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "snapshotpolicy.SnapshotPolicyClient", "SnapshotPoliciesDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSnapshotPoliciesDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "snapshotpolicy.SnapshotPolicyClient", "SnapshotPoliciesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SnapshotPoliciesDeleteThenPoll performs SnapshotPoliciesDelete then polls until it's completed
func (c SnapshotPolicyClient) SnapshotPoliciesDeleteThenPoll(ctx context.Context, id SnapshotPoliciesId) error {
	result, err := c.SnapshotPoliciesDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing SnapshotPoliciesDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SnapshotPoliciesDelete: %+v", err)
	}

	return nil
}

// preparerForSnapshotPoliciesDelete prepares the SnapshotPoliciesDelete request.
func (c SnapshotPolicyClient) preparerForSnapshotPoliciesDelete(ctx context.Context, id SnapshotPoliciesId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForSnapshotPoliciesDelete sends the SnapshotPoliciesDelete request. The method will close the
// http.Response Body if it receives an error.
func (c SnapshotPolicyClient) senderForSnapshotPoliciesDelete(ctx context.Context, req *http.Request) (future SnapshotPoliciesDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
