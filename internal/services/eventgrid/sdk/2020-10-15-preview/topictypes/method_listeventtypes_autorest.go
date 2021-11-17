package topictypes

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
func (c TopicTypesClient) ListEventTypes(ctx context.Context, id TopicTypeId) (result ListEventTypesResponse, err error) {
	req, err := c.preparerForListEventTypes(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "topictypes.TopicTypesClient", "ListEventTypes", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "topictypes.TopicTypesClient", "ListEventTypes", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListEventTypes(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "topictypes.TopicTypesClient", "ListEventTypes", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListEventTypes prepares the ListEventTypes request.
func (c TopicTypesClient) preparerForListEventTypes(ctx context.Context, id TopicTypeId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/eventTypes", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListEventTypes handles the response to the ListEventTypes request. The method always
// closes the http.Response Body.
func (c TopicTypesClient) responderForListEventTypes(resp *http.Response) (result ListEventTypesResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
