package attestationproviders

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListResponse struct {
	HttpResponse *http.Response
	Model        *AttestationProviderListResult
}

// List ...
func (c AttestationProvidersClient) List(ctx context.Context, id SubscriptionId) (result ListResponse, err error) {
	req, err := c.preparerForList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "attestationproviders.AttestationProvidersClient", "List", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "attestationproviders.AttestationProvidersClient", "List", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "attestationproviders.AttestationProvidersClient", "List", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForList prepares the List request.
func (c AttestationProvidersClient) preparerForList(ctx context.Context, id SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Attestation/attestationProviders", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForList handles the response to the List request. The method always
// closes the http.Response Body.
func (c AttestationProvidersClient) responderForList(resp *http.Response) (result ListResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
