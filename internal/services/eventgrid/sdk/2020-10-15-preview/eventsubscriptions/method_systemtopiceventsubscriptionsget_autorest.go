package eventsubscriptions

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type SystemTopicEventSubscriptionsGetResponse struct {
	HttpResponse *http.Response
	Model        *EventSubscription
}

// SystemTopicEventSubscriptionsGet ...
func (c EventSubscriptionsClient) SystemTopicEventSubscriptionsGet(ctx context.Context, id EventSubscriptionId) (result SystemTopicEventSubscriptionsGetResponse, err error) {
	req, err := c.preparerForSystemTopicEventSubscriptionsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "SystemTopicEventSubscriptionsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "SystemTopicEventSubscriptionsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSystemTopicEventSubscriptionsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "SystemTopicEventSubscriptionsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSystemTopicEventSubscriptionsGet prepares the SystemTopicEventSubscriptionsGet request.
func (c EventSubscriptionsClient) preparerForSystemTopicEventSubscriptionsGet(ctx context.Context, id EventSubscriptionId) (*http.Request, error) {
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

// responderForSystemTopicEventSubscriptionsGet handles the response to the SystemTopicEventSubscriptionsGet request. The method always
// closes the http.Response Body.
func (c EventSubscriptionsClient) responderForSystemTopicEventSubscriptionsGet(resp *http.Response) (result SystemTopicEventSubscriptionsGetResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
