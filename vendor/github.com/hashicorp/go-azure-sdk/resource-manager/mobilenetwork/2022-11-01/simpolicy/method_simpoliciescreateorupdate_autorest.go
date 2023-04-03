package simpolicy

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

type SimPoliciesCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SimPoliciesCreateOrUpdate ...
func (c SIMPolicyClient) SimPoliciesCreateOrUpdate(ctx context.Context, id SimPolicyId, input SimPolicy) (result SimPoliciesCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForSimPoliciesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "simpolicy.SIMPolicyClient", "SimPoliciesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSimPoliciesCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "simpolicy.SIMPolicyClient", "SimPoliciesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SimPoliciesCreateOrUpdateThenPoll performs SimPoliciesCreateOrUpdate then polls until it's completed
func (c SIMPolicyClient) SimPoliciesCreateOrUpdateThenPoll(ctx context.Context, id SimPolicyId, input SimPolicy) error {
	result, err := c.SimPoliciesCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing SimPoliciesCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SimPoliciesCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForSimPoliciesCreateOrUpdate prepares the SimPoliciesCreateOrUpdate request.
func (c SIMPolicyClient) preparerForSimPoliciesCreateOrUpdate(ctx context.Context, id SimPolicyId, input SimPolicy) (*http.Request, error) {
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

// senderForSimPoliciesCreateOrUpdate sends the SimPoliciesCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c SIMPolicyClient) senderForSimPoliciesCreateOrUpdate(ctx context.Context, req *http.Request) (future SimPoliciesCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
