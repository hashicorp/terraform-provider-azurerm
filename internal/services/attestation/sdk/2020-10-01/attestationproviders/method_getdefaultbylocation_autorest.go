package attestationproviders

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type GetDefaultByLocationResponse struct {
	HttpResponse *http.Response
	Model        *AttestationProviders
}

// GetDefaultByLocation ...
func (c AttestationProvidersClient) GetDefaultByLocation(ctx context.Context, id LocationId) (result GetDefaultByLocationResponse, err error) {
	req, err := c.preparerForGetDefaultByLocation(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "attestationproviders.AttestationProvidersClient", "GetDefaultByLocation", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "attestationproviders.AttestationProvidersClient", "GetDefaultByLocation", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetDefaultByLocation(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "attestationproviders.AttestationProvidersClient", "GetDefaultByLocation", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetDefaultByLocation prepares the GetDefaultByLocation request.
func (c AttestationProvidersClient) preparerForGetDefaultByLocation(ctx context.Context, id LocationId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/defaultProvider", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetDefaultByLocation handles the response to the GetDefaultByLocation request. The method always
// closes the http.Response Body.
func (c AttestationProvidersClient) responderForGetDefaultByLocation(resp *http.Response) (result GetDefaultByLocationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
