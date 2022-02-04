package diskpools

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type StartResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Start ...
func (c DiskPoolsClient) Start(ctx context.Context, id DiskPoolId) (result StartResponse, err error) {
	req, err := c.preparerForStart(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "diskpools.DiskPoolsClient", "Start", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForStart(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "diskpools.DiskPoolsClient", "Start", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// StartThenPoll performs Start then polls until it's completed
func (c DiskPoolsClient) StartThenPoll(ctx context.Context, id DiskPoolId) error {
	result, err := c.Start(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Start: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Start: %+v", err)
	}

	return nil
}

// preparerForStart prepares the Start request.
func (c DiskPoolsClient) preparerForStart(ctx context.Context, id DiskPoolId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/start", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForStart sends the Start request. The method will close the
// http.Response Body if it receives an error.
func (c DiskPoolsClient) senderForStart(ctx context.Context, req *http.Request) (future StartResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
