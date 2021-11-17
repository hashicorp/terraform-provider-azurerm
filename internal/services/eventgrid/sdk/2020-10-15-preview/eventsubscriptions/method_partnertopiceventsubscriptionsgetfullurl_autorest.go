package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type PartnerTopicEventSubscriptionsGetFullUrlResponse struct {
	HttpResponse *http.Response
	Model        *EventSubscriptionFullUrl
}

// PartnerTopicEventSubscriptionsGetFullUrl ...
func (c EventSubscriptionsClient) PartnerTopicEventSubscriptionsGetFullUrl(ctx context.Context, id PartnerTopicEventSubscriptionId) (result PartnerTopicEventSubscriptionsGetFullUrlResponse, err error) {
	req, err := c.preparerForPartnerTopicEventSubscriptionsGetFullUrl(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "PartnerTopicEventSubscriptionsGetFullUrl", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "PartnerTopicEventSubscriptionsGetFullUrl", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPartnerTopicEventSubscriptionsGetFullUrl(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "PartnerTopicEventSubscriptionsGetFullUrl", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPartnerTopicEventSubscriptionsGetFullUrl prepares the PartnerTopicEventSubscriptionsGetFullUrl request.
func (c EventSubscriptionsClient) preparerForPartnerTopicEventSubscriptionsGetFullUrl(ctx context.Context, id PartnerTopicEventSubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/getFullUrl", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForPartnerTopicEventSubscriptionsGetFullUrl handles the response to the PartnerTopicEventSubscriptionsGetFullUrl request. The method always
// closes the http.Response Body.
func (c EventSubscriptionsClient) responderForPartnerTopicEventSubscriptionsGetFullUrl(resp *http.Response) (result PartnerTopicEventSubscriptionsGetFullUrlResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
