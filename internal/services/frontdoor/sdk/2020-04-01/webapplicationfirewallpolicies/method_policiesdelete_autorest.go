package webapplicationfirewallpolicies

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type PoliciesDeleteResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// PoliciesDelete ...
func (c WebApplicationFirewallPoliciesClient) PoliciesDelete(ctx context.Context, id FrontDoorWebApplicationFirewallPoliciesId) (result PoliciesDeleteResponse, err error) {
	req, err := c.preparerForPoliciesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPoliciesDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PoliciesDeleteThenPoll performs PoliciesDelete then polls until it's completed
func (c WebApplicationFirewallPoliciesClient) PoliciesDeleteThenPoll(ctx context.Context, id FrontDoorWebApplicationFirewallPoliciesId) error {
	result, err := c.PoliciesDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing PoliciesDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after PoliciesDelete: %+v", err)
	}

	return nil
}

// preparerForPoliciesDelete prepares the PoliciesDelete request.
func (c WebApplicationFirewallPoliciesClient) preparerForPoliciesDelete(ctx context.Context, id FrontDoorWebApplicationFirewallPoliciesId) (*http.Request, error) {
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

// senderForPoliciesDelete sends the PoliciesDelete request. The method will close the
// http.Response Body if it receives an error.
func (c WebApplicationFirewallPoliciesClient) senderForPoliciesDelete(ctx context.Context, req *http.Request) (future PoliciesDeleteResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
