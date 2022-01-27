package securitypolicies

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type PatchResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Patch ...
func (c SecurityPoliciesClient) Patch(ctx context.Context, id SecurityPoliciesId, input SecurityPolicyUpdateParameters) (result PatchResponse, err error) {
	req, err := c.preparerForPatch(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securitypolicies.SecurityPoliciesClient", "Patch", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPatch(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securitypolicies.SecurityPoliciesClient", "Patch", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PatchThenPoll performs Patch then polls until it's completed
func (c SecurityPoliciesClient) PatchThenPoll(ctx context.Context, id SecurityPoliciesId, input SecurityPolicyUpdateParameters) error {
	result, err := c.Patch(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Patch: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Patch: %+v", err)
	}

	return nil
}

// preparerForPatch prepares the Patch request.
func (c SecurityPoliciesClient) preparerForPatch(ctx context.Context, id SecurityPoliciesId, input SecurityPolicyUpdateParameters) (*http.Request, error) {
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

// senderForPatch sends the Patch request. The method will close the
// http.Response Body if it receives an error.
func (c SecurityPoliciesClient) senderForPatch(ctx context.Context, req *http.Request) (future PatchResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
