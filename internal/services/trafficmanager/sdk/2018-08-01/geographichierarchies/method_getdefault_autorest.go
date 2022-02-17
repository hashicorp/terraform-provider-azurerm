package geographichierarchies

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type GetDefaultResponse struct {
	HttpResponse *http.Response
	Model        *TrafficManagerGeographicHierarchy
}

// GetDefault ...
func (c GeographicHierarchiesClient) GetDefault(ctx context.Context) (result GetDefaultResponse, err error) {
	req, err := c.preparerForGetDefault(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "geographichierarchies.GeographicHierarchiesClient", "GetDefault", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "geographichierarchies.GeographicHierarchiesClient", "GetDefault", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetDefault(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "geographichierarchies.GeographicHierarchiesClient", "GetDefault", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetDefault prepares the GetDefault request.
func (c GeographicHierarchiesClient) preparerForGetDefault(ctx context.Context) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath("/providers/Microsoft.Network/trafficManagerGeographicHierarchies/default"),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetDefault handles the response to the GetDefault request. The method always
// closes the http.Response Body.
func (c GeographicHierarchiesClient) responderForGetDefault(resp *http.Response) (result GetDefaultResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
