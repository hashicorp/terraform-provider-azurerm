package databases

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type ForceUnlinkResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ForceUnlink ...
func (c DatabasesClient) ForceUnlink(ctx context.Context, id DatabaseId, input ForceUnlinkParameters) (result ForceUnlinkResponse, err error) {
	req, err := c.preparerForForceUnlink(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "ForceUnlink", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForForceUnlink(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "ForceUnlink", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ForceUnlinkThenPoll performs ForceUnlink then polls until it's completed
func (c DatabasesClient) ForceUnlinkThenPoll(ctx context.Context, id DatabaseId, input ForceUnlinkParameters) error {
	result, err := c.ForceUnlink(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ForceUnlink: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ForceUnlink: %+v", err)
	}

	return nil
}

// preparerForForceUnlink prepares the ForceUnlink request.
func (c DatabasesClient) preparerForForceUnlink(ctx context.Context, id DatabaseId, input ForceUnlinkParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/forceUnlink", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForForceUnlink sends the ForceUnlink request. The method will close the
// http.Response Body if it receives an error.
func (c DatabasesClient) senderForForceUnlink(ctx context.Context, req *http.Request) (future ForceUnlinkResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
