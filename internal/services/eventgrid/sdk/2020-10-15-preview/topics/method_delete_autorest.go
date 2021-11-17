package topics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type DeleteResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Delete ...
func (c TopicsClient) Delete(ctx context.Context, id TopicId) (result DeleteResponse, err error) {
	req, err := c.preparerForDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "topics.TopicsClient", "Delete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "topics.TopicsClient", "Delete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DeleteThenPoll performs Delete then polls until it's completed
func (c TopicsClient) DeleteThenPoll(ctx context.Context, id TopicId) error {
	result, err := c.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Delete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Delete: %+v", err)
	}

	return nil
}

// preparerForDelete prepares the Delete request.
func (c TopicsClient) preparerForDelete(ctx context.Context, id TopicId) (*http.Request, error) {
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

// senderForDelete sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (c TopicsClient) senderForDelete(ctx context.Context, req *http.Request) (future DeleteResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
