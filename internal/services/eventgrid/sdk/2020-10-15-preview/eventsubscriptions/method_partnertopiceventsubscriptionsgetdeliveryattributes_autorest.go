package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type PartnerTopicEventSubscriptionsGetDeliveryAttributesResponse struct {
	HttpResponse *http.Response
	Model        *DeliveryAttributeListResult
}

// PartnerTopicEventSubscriptionsGetDeliveryAttributes ...
func (c EventSubscriptionsClient) PartnerTopicEventSubscriptionsGetDeliveryAttributes(ctx context.Context, id PartnerTopicEventSubscriptionId) (result PartnerTopicEventSubscriptionsGetDeliveryAttributesResponse, err error) {
	req, err := c.preparerForPartnerTopicEventSubscriptionsGetDeliveryAttributes(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "PartnerTopicEventSubscriptionsGetDeliveryAttributes", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "PartnerTopicEventSubscriptionsGetDeliveryAttributes", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPartnerTopicEventSubscriptionsGetDeliveryAttributes(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "PartnerTopicEventSubscriptionsGetDeliveryAttributes", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPartnerTopicEventSubscriptionsGetDeliveryAttributes prepares the PartnerTopicEventSubscriptionsGetDeliveryAttributes request.
func (c EventSubscriptionsClient) preparerForPartnerTopicEventSubscriptionsGetDeliveryAttributes(ctx context.Context, id PartnerTopicEventSubscriptionId) (*http.Request, error) {
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

// responderForPartnerTopicEventSubscriptionsGetDeliveryAttributes handles the response to the PartnerTopicEventSubscriptionsGetDeliveryAttributes request. The method always
// closes the http.Response Body.
func (c EventSubscriptionsClient) responderForPartnerTopicEventSubscriptionsGetDeliveryAttributes(resp *http.Response) (result PartnerTopicEventSubscriptionsGetDeliveryAttributesResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
