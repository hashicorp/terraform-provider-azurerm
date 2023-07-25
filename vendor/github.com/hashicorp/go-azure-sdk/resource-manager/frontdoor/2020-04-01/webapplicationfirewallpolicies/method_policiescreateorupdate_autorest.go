package webapplicationfirewallpolicies

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

type PoliciesCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// PoliciesCreateOrUpdate ...
func (c WebApplicationFirewallPoliciesClient) PoliciesCreateOrUpdate(ctx context.Context, id FrontDoorWebApplicationFirewallPolicyId, input WebApplicationFirewallPolicy) (result PoliciesCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForPoliciesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPoliciesCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PoliciesCreateOrUpdateThenPoll performs PoliciesCreateOrUpdate then polls until it's completed
func (c WebApplicationFirewallPoliciesClient) PoliciesCreateOrUpdateThenPoll(ctx context.Context, id FrontDoorWebApplicationFirewallPolicyId, input WebApplicationFirewallPolicy) error {
	result, err := c.PoliciesCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing PoliciesCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after PoliciesCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForPoliciesCreateOrUpdate prepares the PoliciesCreateOrUpdate request.
func (c WebApplicationFirewallPoliciesClient) preparerForPoliciesCreateOrUpdate(ctx context.Context, id FrontDoorWebApplicationFirewallPolicyId, input WebApplicationFirewallPolicy) (*http.Request, error) {
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

// senderForPoliciesCreateOrUpdate sends the PoliciesCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c WebApplicationFirewallPoliciesClient) senderForPoliciesCreateOrUpdate(ctx context.Context, req *http.Request) (future PoliciesCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
