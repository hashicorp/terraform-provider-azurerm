package powerbidedicated

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type CapacitiesListSkusResponse struct {
	HttpResponse *http.Response
	Model        *SkuEnumerationForNewResourceResult
}

// CapacitiesListSkus ...
func (c PowerBIDedicatedClient) CapacitiesListSkus(ctx context.Context, id SubscriptionId) (result CapacitiesListSkusResponse, err error) {
	req, err := c.preparerForCapacitiesListSkus(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "powerbidedicated.PowerBIDedicatedClient", "CapacitiesListSkus", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "powerbidedicated.PowerBIDedicatedClient", "CapacitiesListSkus", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCapacitiesListSkus(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "powerbidedicated.PowerBIDedicatedClient", "CapacitiesListSkus", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCapacitiesListSkus prepares the CapacitiesListSkus request.
func (c PowerBIDedicatedClient) preparerForCapacitiesListSkus(ctx context.Context, id SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.PowerBIDedicated/skus", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCapacitiesListSkus handles the response to the CapacitiesListSkus request. The method always
// closes the http.Response Body.
func (c PowerBIDedicatedClient) responderForCapacitiesListSkus(resp *http.Response) (result CapacitiesListSkusResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
