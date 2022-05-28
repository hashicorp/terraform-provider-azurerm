package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

type AccountsCheckNameAvailabilityResponse struct {
	HttpResponse *http.Response
	Model        *CheckNameAvailabilityResult
}

// AccountsCheckNameAvailability ...
func (c ProviderClient) AccountsCheckNameAvailability(ctx context.Context, id commonids.SubscriptionId, input CheckNameAvailabilityRequest) (result AccountsCheckNameAvailabilityResponse, err error) {
	req, err := c.preparerForAccountsCheckNameAvailability(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "provider.ProviderClient", "AccountsCheckNameAvailability", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "provider.ProviderClient", "AccountsCheckNameAvailability", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAccountsCheckNameAvailability(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "provider.ProviderClient", "AccountsCheckNameAvailability", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAccountsCheckNameAvailability prepares the AccountsCheckNameAvailability request.
func (c ProviderClient) preparerForAccountsCheckNameAvailability(ctx context.Context, id commonids.SubscriptionId, input CheckNameAvailabilityRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Purview/checkNameAvailability", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAccountsCheckNameAvailability handles the response to the AccountsCheckNameAvailability request. The method always
// closes the http.Response Body.
func (c ProviderClient) responderForAccountsCheckNameAvailability(resp *http.Response) (result AccountsCheckNameAvailabilityResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
