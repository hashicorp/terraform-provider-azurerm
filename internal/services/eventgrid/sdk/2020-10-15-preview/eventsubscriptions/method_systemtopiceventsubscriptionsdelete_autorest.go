package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type SystemTopicEventSubscriptionsDeleteResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SystemTopicEventSubscriptionsDelete ...
func (c EventSubscriptionsClient) SystemTopicEventSubscriptionsDelete(ctx context.Context, id EventSubscriptionId) (result SystemTopicEventSubscriptionsDeleteResponse, err error) {
	req, err := c.preparerForSystemTopicEventSubscriptionsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "SystemTopicEventSubscriptionsDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSystemTopicEventSubscriptionsDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "SystemTopicEventSubscriptionsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SystemTopicEventSubscriptionsDeleteThenPoll performs SystemTopicEventSubscriptionsDelete then polls until it's completed
func (c EventSubscriptionsClient) SystemTopicEventSubscriptionsDeleteThenPoll(ctx context.Context, id EventSubscriptionId) error {
	result, err := c.SystemTopicEventSubscriptionsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing SystemTopicEventSubscriptionsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SystemTopicEventSubscriptionsDelete: %+v", err)
	}

	return nil
}

// preparerForSystemTopicEventSubscriptionsDelete prepares the SystemTopicEventSubscriptionsDelete request.
func (c EventSubscriptionsClient) preparerForSystemTopicEventSubscriptionsDelete(ctx context.Context, id EventSubscriptionId) (*http.Request, error) {
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

// senderForSystemTopicEventSubscriptionsDelete sends the SystemTopicEventSubscriptionsDelete request. The method will close the
// http.Response Body if it receives an error.
func (c EventSubscriptionsClient) senderForSystemTopicEventSubscriptionsDelete(ctx context.Context, req *http.Request) (future SystemTopicEventSubscriptionsDeleteResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
