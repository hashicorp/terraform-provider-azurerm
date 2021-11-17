package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type PartnerTopicEventSubscriptionsUpdateResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// PartnerTopicEventSubscriptionsUpdate ...
func (c EventSubscriptionsClient) PartnerTopicEventSubscriptionsUpdate(ctx context.Context, id PartnerTopicEventSubscriptionId, input EventSubscriptionUpdateParameters) (result PartnerTopicEventSubscriptionsUpdateResponse, err error) {
	req, err := c.preparerForPartnerTopicEventSubscriptionsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "PartnerTopicEventSubscriptionsUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPartnerTopicEventSubscriptionsUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "PartnerTopicEventSubscriptionsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PartnerTopicEventSubscriptionsUpdateThenPoll performs PartnerTopicEventSubscriptionsUpdate then polls until it's completed
func (c EventSubscriptionsClient) PartnerTopicEventSubscriptionsUpdateThenPoll(ctx context.Context, id PartnerTopicEventSubscriptionId, input EventSubscriptionUpdateParameters) error {
	result, err := c.PartnerTopicEventSubscriptionsUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing PartnerTopicEventSubscriptionsUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after PartnerTopicEventSubscriptionsUpdate: %+v", err)
	}

	return nil
}

// preparerForPartnerTopicEventSubscriptionsUpdate prepares the PartnerTopicEventSubscriptionsUpdate request.
func (c EventSubscriptionsClient) preparerForPartnerTopicEventSubscriptionsUpdate(ctx context.Context, id PartnerTopicEventSubscriptionId, input EventSubscriptionUpdateParameters) (*http.Request, error) {
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

// senderForPartnerTopicEventSubscriptionsUpdate sends the PartnerTopicEventSubscriptionsUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c EventSubscriptionsClient) senderForPartnerTopicEventSubscriptionsUpdate(ctx context.Context, req *http.Request) (future PartnerTopicEventSubscriptionsUpdateResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
