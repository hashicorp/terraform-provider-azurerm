package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type SystemTopicEventSubscriptionsGetFullUrlResponse struct {
	HttpResponse *http.Response
	Model        *EventSubscriptionFullUrl
}

// SystemTopicEventSubscriptionsGetFullUrl ...
func (c EventSubscriptionsClient) SystemTopicEventSubscriptionsGetFullUrl(ctx context.Context, id EventSubscriptionId) (result SystemTopicEventSubscriptionsGetFullUrlResponse, err error) {
	req, err := c.preparerForSystemTopicEventSubscriptionsGetFullUrl(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "SystemTopicEventSubscriptionsGetFullUrl", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "SystemTopicEventSubscriptionsGetFullUrl", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSystemTopicEventSubscriptionsGetFullUrl(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "SystemTopicEventSubscriptionsGetFullUrl", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSystemTopicEventSubscriptionsGetFullUrl prepares the SystemTopicEventSubscriptionsGetFullUrl request.
func (c EventSubscriptionsClient) preparerForSystemTopicEventSubscriptionsGetFullUrl(ctx context.Context, id EventSubscriptionId) (*http.Request, error) {
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

// responderForSystemTopicEventSubscriptionsGetFullUrl handles the response to the SystemTopicEventSubscriptionsGetFullUrl request. The method always
// closes the http.Response Body.
func (c EventSubscriptionsClient) responderForSystemTopicEventSubscriptionsGetFullUrl(resp *http.Response) (result SystemTopicEventSubscriptionsGetFullUrlResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
