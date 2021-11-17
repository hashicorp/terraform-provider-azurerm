package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type SystemTopicEventSubscriptionsUpdateResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SystemTopicEventSubscriptionsUpdate ...
func (c EventSubscriptionsClient) SystemTopicEventSubscriptionsUpdate(ctx context.Context, id EventSubscriptionId, input EventSubscriptionUpdateParameters) (result SystemTopicEventSubscriptionsUpdateResponse, err error) {
	req, err := c.preparerForSystemTopicEventSubscriptionsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "SystemTopicEventSubscriptionsUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSystemTopicEventSubscriptionsUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "SystemTopicEventSubscriptionsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SystemTopicEventSubscriptionsUpdateThenPoll performs SystemTopicEventSubscriptionsUpdate then polls until it's completed
func (c EventSubscriptionsClient) SystemTopicEventSubscriptionsUpdateThenPoll(ctx context.Context, id EventSubscriptionId, input EventSubscriptionUpdateParameters) error {
	result, err := c.SystemTopicEventSubscriptionsUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing SystemTopicEventSubscriptionsUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SystemTopicEventSubscriptionsUpdate: %+v", err)
	}

	return nil
}

// preparerForSystemTopicEventSubscriptionsUpdate prepares the SystemTopicEventSubscriptionsUpdate request.
func (c EventSubscriptionsClient) preparerForSystemTopicEventSubscriptionsUpdate(ctx context.Context, id EventSubscriptionId, input EventSubscriptionUpdateParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForSystemTopicEventSubscriptionsUpdate sends the SystemTopicEventSubscriptionsUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c EventSubscriptionsClient) senderForSystemTopicEventSubscriptionsUpdate(ctx context.Context, req *http.Request) (future SystemTopicEventSubscriptionsUpdateResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
