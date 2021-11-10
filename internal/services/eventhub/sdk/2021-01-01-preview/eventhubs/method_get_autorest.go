package eventhubs

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type GetResponse struct {
	HttpResponse *http.Response
	Model        *Eventhub
}

// Get ...
func (c EventHubsClient) Get(ctx context.Context, id EventhubId) (result GetResponse, err error) {
	req, err := c.preparerForGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubs.EventHubsClient", "Get", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubs.EventHubsClient", "Get", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubs.EventHubsClient", "Get", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGet prepares the Get request.
func (c EventHubsClient) preparerForGet(ctx context.Context, id EventhubId) (*http.Request, error) {
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

// responderForGet handles the response to the Get request. The method always
// closes the http.Response Body.
func (c EventHubsClient) responderForGet(resp *http.Response) (result GetResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
