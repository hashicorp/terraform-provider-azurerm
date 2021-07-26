package eventhubs

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type DeleteAuthorizationRuleResponse struct {
	HttpResponse *http.Response
}

// DeleteAuthorizationRule ...
func (c EventHubsClient) DeleteAuthorizationRule(ctx context.Context, id AuthorizationRuleId) (result DeleteAuthorizationRuleResponse, err error) {
	req, err := c.preparerForDeleteAuthorizationRule(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubs.EventHubsClient", "DeleteAuthorizationRule", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubs.EventHubsClient", "DeleteAuthorizationRule", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDeleteAuthorizationRule(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubs.EventHubsClient", "DeleteAuthorizationRule", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDeleteAuthorizationRule prepares the DeleteAuthorizationRule request.
func (c EventHubsClient) preparerForDeleteAuthorizationRule(ctx context.Context, id AuthorizationRuleId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDeleteAuthorizationRule handles the response to the DeleteAuthorizationRule request. The method always
// closes the http.Response Body.
func (c EventHubsClient) responderForDeleteAuthorizationRule(resp *http.Response) (result DeleteAuthorizationRuleResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
