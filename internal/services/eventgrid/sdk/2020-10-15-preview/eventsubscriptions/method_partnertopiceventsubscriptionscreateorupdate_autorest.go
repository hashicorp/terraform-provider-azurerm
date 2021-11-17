package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type PartnerTopicEventSubscriptionsCreateOrUpdateResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// PartnerTopicEventSubscriptionsCreateOrUpdate ...
func (c EventSubscriptionsClient) PartnerTopicEventSubscriptionsCreateOrUpdate(ctx context.Context, id PartnerTopicEventSubscriptionId, input EventSubscription) (result PartnerTopicEventSubscriptionsCreateOrUpdateResponse, err error) {
	req, err := c.preparerForPartnerTopicEventSubscriptionsCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "PartnerTopicEventSubscriptionsCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPartnerTopicEventSubscriptionsCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "PartnerTopicEventSubscriptionsCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PartnerTopicEventSubscriptionsCreateOrUpdateThenPoll performs PartnerTopicEventSubscriptionsCreateOrUpdate then polls until it's completed
func (c EventSubscriptionsClient) PartnerTopicEventSubscriptionsCreateOrUpdateThenPoll(ctx context.Context, id PartnerTopicEventSubscriptionId, input EventSubscription) error {
	result, err := c.PartnerTopicEventSubscriptionsCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing PartnerTopicEventSubscriptionsCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after PartnerTopicEventSubscriptionsCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForPartnerTopicEventSubscriptionsCreateOrUpdate prepares the PartnerTopicEventSubscriptionsCreateOrUpdate request.
func (c EventSubscriptionsClient) preparerForPartnerTopicEventSubscriptionsCreateOrUpdate(ctx context.Context, id PartnerTopicEventSubscriptionId, input EventSubscription) (*http.Request, error) {
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

// senderForPartnerTopicEventSubscriptionsCreateOrUpdate sends the PartnerTopicEventSubscriptionsCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c EventSubscriptionsClient) senderForPartnerTopicEventSubscriptionsCreateOrUpdate(ctx context.Context, req *http.Request) (future PartnerTopicEventSubscriptionsCreateOrUpdateResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
