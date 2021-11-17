package topics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListEventTypesResponse struct {
	HttpResponse *http.Response
	Model        *EventTypesListResult
}

// ListEventTypes ...
func (c TopicsClient) ListEventTypes(ctx context.Context, id ProviderId) (result ListEventTypesResponse, err error) {
	req, err := c.preparerForListEventTypes(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "topics.TopicsClient", "ListEventTypes", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "topics.TopicsClient", "ListEventTypes", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListEventTypes(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "topics.TopicsClient", "ListEventTypes", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListEventTypes prepares the ListEventTypes request.
func (c TopicsClient) preparerForListEventTypes(ctx context.Context, id ProviderId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.EventGrid/eventTypes", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListEventTypes handles the response to the ListEventTypes request. The method always
// closes the http.Response Body.
func (c TopicsClient) responderForListEventTypes(resp *http.Response) (result ListEventTypesResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
