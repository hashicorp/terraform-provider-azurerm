package webapplicationfirewallpolicies

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type PoliciesDeleteResponse struct {
	HttpResponse *http.Response
}

// PoliciesDelete ...
func (c WebApplicationFirewallPoliciesClient) PoliciesDelete(ctx context.Context, id CdnWebApplicationFirewallPoliciesId) (result PoliciesDeleteResponse, err error) {
	req, err := c.preparerForPoliciesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPoliciesDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPoliciesDelete prepares the PoliciesDelete request.
func (c WebApplicationFirewallPoliciesClient) preparerForPoliciesDelete(ctx context.Context, id CdnWebApplicationFirewallPoliciesId) (*http.Request, error) {
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

// responderForPoliciesDelete handles the response to the PoliciesDelete request. The method always
// closes the http.Response Body.
func (c WebApplicationFirewallPoliciesClient) responderForPoliciesDelete(resp *http.Response) (result PoliciesDeleteResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
