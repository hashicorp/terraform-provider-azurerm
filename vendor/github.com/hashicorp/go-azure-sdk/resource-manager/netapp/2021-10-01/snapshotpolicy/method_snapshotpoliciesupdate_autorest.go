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

type SnapshotPoliciesUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SnapshotPoliciesUpdate ...
func (c SnapshotPolicyClient) SnapshotPoliciesUpdate(ctx context.Context, id SnapshotPoliciesId, input SnapshotPolicyPatch) (result SnapshotPoliciesUpdateOperationResponse, err error) {
	req, err := c.preparerForSnapshotPoliciesUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "snapshotpolicy.SnapshotPolicyClient", "SnapshotPoliciesUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSnapshotPoliciesUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "snapshotpolicy.SnapshotPolicyClient", "SnapshotPoliciesUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SnapshotPoliciesUpdateThenPoll performs SnapshotPoliciesUpdate then polls until it's completed
func (c SnapshotPolicyClient) SnapshotPoliciesUpdateThenPoll(ctx context.Context, id SnapshotPoliciesId, input SnapshotPolicyPatch) error {
	result, err := c.SnapshotPoliciesUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing SnapshotPoliciesUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SnapshotPoliciesUpdate: %+v", err)
	}

	return nil
}

// preparerForSnapshotPoliciesUpdate prepares the SnapshotPoliciesUpdate request.
func (c SnapshotPolicyClient) preparerForSnapshotPoliciesUpdate(ctx context.Context, id SnapshotPoliciesId, input SnapshotPolicyPatch) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForSnapshotPoliciesUpdate sends the SnapshotPoliciesUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c SnapshotPolicyClient) senderForSnapshotPoliciesUpdate(ctx context.Context, req *http.Request) (future SnapshotPoliciesUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
