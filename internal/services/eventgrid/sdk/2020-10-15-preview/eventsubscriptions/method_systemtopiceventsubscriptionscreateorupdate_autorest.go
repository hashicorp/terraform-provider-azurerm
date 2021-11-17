package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type SystemTopicEventSubscriptionsCreateOrUpdateResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SystemTopicEventSubscriptionsCreateOrUpdate ...
func (c EventSubscriptionsClient) SystemTopicEventSubscriptionsCreateOrUpdate(ctx context.Context, id EventSubscriptionId, input EventSubscription) (result SystemTopicEventSubscriptionsCreateOrUpdateResponse, err error) {
	req, err := c.preparerForSystemTopicEventSubscriptionsCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "SystemTopicEventSubscriptionsCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSystemTopicEventSubscriptionsCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "SystemTopicEventSubscriptionsCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SystemTopicEventSubscriptionsCreateOrUpdateThenPoll performs SystemTopicEventSubscriptionsCreateOrUpdate then polls until it's completed
func (c EventSubscriptionsClient) SystemTopicEventSubscriptionsCreateOrUpdateThenPoll(ctx context.Context, id EventSubscriptionId, input EventSubscription) error {
	result, err := c.SystemTopicEventSubscriptionsCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing SystemTopicEventSubscriptionsCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SystemTopicEventSubscriptionsCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForSystemTopicEventSubscriptionsCreateOrUpdate prepares the SystemTopicEventSubscriptionsCreateOrUpdate request.
func (c EventSubscriptionsClient) preparerForSystemTopicEventSubscriptionsCreateOrUpdate(ctx context.Context, id EventSubscriptionId, input EventSubscription) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForSystemTopicEventSubscriptionsCreateOrUpdate sends the SystemTopicEventSubscriptionsCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c EventSubscriptionsClient) senderForSystemTopicEventSubscriptionsCreateOrUpdate(ctx context.Context, req *http.Request) (future SystemTopicEventSubscriptionsCreateOrUpdateResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
