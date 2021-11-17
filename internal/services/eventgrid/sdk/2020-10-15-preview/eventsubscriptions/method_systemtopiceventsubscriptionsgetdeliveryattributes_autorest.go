package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type SystemTopicEventSubscriptionsGetDeliveryAttributesResponse struct {
	HttpResponse *http.Response
	Model        *DeliveryAttributeListResult
}

// SystemTopicEventSubscriptionsGetDeliveryAttributes ...
func (c EventSubscriptionsClient) SystemTopicEventSubscriptionsGetDeliveryAttributes(ctx context.Context, id EventSubscriptionId) (result SystemTopicEventSubscriptionsGetDeliveryAttributesResponse, err error) {
	req, err := c.preparerForSystemTopicEventSubscriptionsGetDeliveryAttributes(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "SystemTopicEventSubscriptionsGetDeliveryAttributes", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "SystemTopicEventSubscriptionsGetDeliveryAttributes", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSystemTopicEventSubscriptionsGetDeliveryAttributes(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "SystemTopicEventSubscriptionsGetDeliveryAttributes", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSystemTopicEventSubscriptionsGetDeliveryAttributes prepares the SystemTopicEventSubscriptionsGetDeliveryAttributes request.
func (c EventSubscriptionsClient) preparerForSystemTopicEventSubscriptionsGetDeliveryAttributes(ctx context.Context, id EventSubscriptionId) (*http.Request, error) {
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

// responderForSystemTopicEventSubscriptionsGetDeliveryAttributes handles the response to the SystemTopicEventSubscriptionsGetDeliveryAttributes request. The method always
// closes the http.Response Body.
func (c EventSubscriptionsClient) responderForSystemTopicEventSubscriptionsGetDeliveryAttributes(resp *http.Response) (result SystemTopicEventSubscriptionsGetDeliveryAttributesResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
