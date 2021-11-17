package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type GetDeliveryAttributesResponse struct {
	HttpResponse *http.Response
	Model        *DeliveryAttributeListResult
}

// GetDeliveryAttributes ...
func (c EventSubscriptionsClient) GetDeliveryAttributes(ctx context.Context, id ScopedEventSubscriptionId) (result GetDeliveryAttributesResponse, err error) {
	req, err := c.preparerForGetDeliveryAttributes(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "GetDeliveryAttributes", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "GetDeliveryAttributes", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetDeliveryAttributes(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "GetDeliveryAttributes", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetDeliveryAttributes prepares the GetDeliveryAttributes request.
func (c EventSubscriptionsClient) preparerForGetDeliveryAttributes(ctx context.Context, id ScopedEventSubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/getDeliveryAttributes", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetDeliveryAttributes handles the response to the GetDeliveryAttributes request. The method always
// closes the http.Response Body.
func (c EventSubscriptionsClient) responderForGetDeliveryAttributes(resp *http.Response) (result GetDeliveryAttributesResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
