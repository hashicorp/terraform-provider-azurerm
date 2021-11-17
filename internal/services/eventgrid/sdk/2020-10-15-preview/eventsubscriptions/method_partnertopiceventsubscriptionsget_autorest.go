package eventsubscriptions

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type PartnerTopicEventSubscriptionsGetResponse struct {
	HttpResponse *http.Response
	Model        *EventSubscription
}

// PartnerTopicEventSubscriptionsGet ...
func (c EventSubscriptionsClient) PartnerTopicEventSubscriptionsGet(ctx context.Context, id PartnerTopicEventSubscriptionId) (result PartnerTopicEventSubscriptionsGetResponse, err error) {
	req, err := c.preparerForPartnerTopicEventSubscriptionsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "PartnerTopicEventSubscriptionsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "PartnerTopicEventSubscriptionsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPartnerTopicEventSubscriptionsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "PartnerTopicEventSubscriptionsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPartnerTopicEventSubscriptionsGet prepares the PartnerTopicEventSubscriptionsGet request.
func (c EventSubscriptionsClient) preparerForPartnerTopicEventSubscriptionsGet(ctx context.Context, id PartnerTopicEventSubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForPartnerTopicEventSubscriptionsGet handles the response to the PartnerTopicEventSubscriptionsGet request. The method always
// closes the http.Response Body.
func (c EventSubscriptionsClient) responderForPartnerTopicEventSubscriptionsGet(resp *http.Response) (result PartnerTopicEventSubscriptionsGetResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
