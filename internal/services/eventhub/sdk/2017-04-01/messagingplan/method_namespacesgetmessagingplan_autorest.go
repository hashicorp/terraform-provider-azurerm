package messagingplan

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type NamespacesGetMessagingPlanResponse struct {
	HttpResponse *http.Response
	Model        *MessagingPlan
}

// NamespacesGetMessagingPlan ...
func (c MessagingPlanClient) NamespacesGetMessagingPlan(ctx context.Context, id NamespaceId) (result NamespacesGetMessagingPlanResponse, err error) {
	req, err := c.preparerForNamespacesGetMessagingPlan(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "messagingplan.MessagingPlanClient", "NamespacesGetMessagingPlan", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "messagingplan.MessagingPlanClient", "NamespacesGetMessagingPlan", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForNamespacesGetMessagingPlan(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "messagingplan.MessagingPlanClient", "NamespacesGetMessagingPlan", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForNamespacesGetMessagingPlan prepares the NamespacesGetMessagingPlan request.
func (c MessagingPlanClient) preparerForNamespacesGetMessagingPlan(ctx context.Context, id NamespaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/messagingplan", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForNamespacesGetMessagingPlan handles the response to the NamespacesGetMessagingPlan request. The method always
// closes the http.Response Body.
func (c MessagingPlanClient) responderForNamespacesGetMessagingPlan(resp *http.Response) (result NamespacesGetMessagingPlanResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
