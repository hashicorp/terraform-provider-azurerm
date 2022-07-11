package diskpools

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type DeallocateResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Deallocate ...
func (c DiskPoolsClient) Deallocate(ctx context.Context, id DiskPoolId) (result DeallocateResponse, err error) {
	req, err := c.preparerForDeallocate(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "diskpools.DiskPoolsClient", "Deallocate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDeallocate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "diskpools.DiskPoolsClient", "Deallocate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DeallocateThenPoll performs Deallocate then polls until it's completed
func (c DiskPoolsClient) DeallocateThenPoll(ctx context.Context, id DiskPoolId) error {
	result, err := c.Deallocate(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Deallocate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Deallocate: %+v", err)
	}

	return nil
}

// preparerForDeallocate prepares the Deallocate request.
func (c DiskPoolsClient) preparerForDeallocate(ctx context.Context, id DiskPoolId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/deallocate", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForDeallocate sends the Deallocate request. The method will close the
// http.Response Body if it receives an error.
func (c DiskPoolsClient) senderForDeallocate(ctx context.Context, req *http.Request) (future DeallocateResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
