package authorizationruleseventhubs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type EventHubsListKeysResponse struct {
	HttpResponse *http.Response
	Model        *AccessKeys
}

// EventHubsListKeys ...
func (c AuthorizationRulesEventHubsClient) EventHubsListKeys(ctx context.Context, id AuthorizationRuleId) (result EventHubsListKeysResponse, err error) {
	req, err := c.preparerForEventHubsListKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationruleseventhubs.AuthorizationRulesEventHubsClient", "EventHubsListKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationruleseventhubs.AuthorizationRulesEventHubsClient", "EventHubsListKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForEventHubsListKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationruleseventhubs.AuthorizationRulesEventHubsClient", "EventHubsListKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForEventHubsListKeys prepares the EventHubsListKeys request.
func (c AuthorizationRulesEventHubsClient) preparerForEventHubsListKeys(ctx context.Context, id AuthorizationRuleId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listKeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForEventHubsListKeys handles the response to the EventHubsListKeys request. The method always
// closes the http.Response Body.
func (c AuthorizationRulesEventHubsClient) responderForEventHubsListKeys(resp *http.Response) (result EventHubsListKeysResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
