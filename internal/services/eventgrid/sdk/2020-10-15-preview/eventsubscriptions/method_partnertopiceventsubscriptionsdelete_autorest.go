package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type PartnerTopicEventSubscriptionsDeleteResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// PartnerTopicEventSubscriptionsDelete ...
func (c EventSubscriptionsClient) PartnerTopicEventSubscriptionsDelete(ctx context.Context, id PartnerTopicEventSubscriptionId) (result PartnerTopicEventSubscriptionsDeleteResponse, err error) {
	req, err := c.preparerForPartnerTopicEventSubscriptionsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "PartnerTopicEventSubscriptionsDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPartnerTopicEventSubscriptionsDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "PartnerTopicEventSubscriptionsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PartnerTopicEventSubscriptionsDeleteThenPoll performs PartnerTopicEventSubscriptionsDelete then polls until it's completed
func (c EventSubscriptionsClient) PartnerTopicEventSubscriptionsDeleteThenPoll(ctx context.Context, id PartnerTopicEventSubscriptionId) error {
	result, err := c.PartnerTopicEventSubscriptionsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing PartnerTopicEventSubscriptionsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after PartnerTopicEventSubscriptionsDelete: %+v", err)
	}

	return nil
}

// preparerForPartnerTopicEventSubscriptionsDelete prepares the PartnerTopicEventSubscriptionsDelete request.
func (c EventSubscriptionsClient) preparerForPartnerTopicEventSubscriptionsDelete(ctx context.Context, id PartnerTopicEventSubscriptionId) (*http.Request, error) {
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

// senderForPartnerTopicEventSubscriptionsDelete sends the PartnerTopicEventSubscriptionsDelete request. The method will close the
// http.Response Body if it receives an error.
func (c EventSubscriptionsClient) senderForPartnerTopicEventSubscriptionsDelete(ctx context.Context, req *http.Request) (future PartnerTopicEventSubscriptionsDeleteResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
