package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type GetFullUrlResponse struct {
	HttpResponse *http.Response
	Model        *EventSubscriptionFullUrl
}

// GetFullUrl ...
func (c EventSubscriptionsClient) GetFullUrl(ctx context.Context, id ScopedEventSubscriptionId) (result GetFullUrlResponse, err error) {
	req, err := c.preparerForGetFullUrl(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "GetFullUrl", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "GetFullUrl", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetFullUrl(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "GetFullUrl", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetFullUrl prepares the GetFullUrl request.
func (c EventSubscriptionsClient) preparerForGetFullUrl(ctx context.Context, id ScopedEventSubscriptionId) (*http.Request, error) {
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

// responderForGetFullUrl handles the response to the GetFullUrl request. The method always
// closes the http.Response Body.
func (c EventSubscriptionsClient) responderForGetFullUrl(resp *http.Response) (result GetFullUrlResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
