package webapplicationfirewallpolicies

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type PoliciesUpdateResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// PoliciesUpdate ...
func (c WebApplicationFirewallPoliciesClient) PoliciesUpdate(ctx context.Context, id CdnWebApplicationFirewallPoliciesId, input CdnWebApplicationFirewallPolicyPatchParameters) (result PoliciesUpdateResponse, err error) {
	req, err := c.preparerForPoliciesUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPoliciesUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PoliciesUpdateThenPoll performs PoliciesUpdate then polls until it's completed
func (c WebApplicationFirewallPoliciesClient) PoliciesUpdateThenPoll(ctx context.Context, id CdnWebApplicationFirewallPoliciesId, input CdnWebApplicationFirewallPolicyPatchParameters) error {
	result, err := c.PoliciesUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing PoliciesUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after PoliciesUpdate: %+v", err)
	}

	return nil
}

// preparerForPoliciesUpdate prepares the PoliciesUpdate request.
func (c WebApplicationFirewallPoliciesClient) preparerForPoliciesUpdate(ctx context.Context, id CdnWebApplicationFirewallPoliciesId, input CdnWebApplicationFirewallPolicyPatchParameters) (*http.Request, error) {
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

// senderForPoliciesUpdate sends the PoliciesUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c WebApplicationFirewallPoliciesClient) senderForPoliciesUpdate(ctx context.Context, req *http.Request) (future PoliciesUpdateResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
