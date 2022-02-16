package diskpools

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type UpgradeResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Upgrade ...
func (c DiskPoolsClient) Upgrade(ctx context.Context, id DiskPoolId) (result UpgradeResponse, err error) {
	req, err := c.preparerForUpgrade(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "diskpools.DiskPoolsClient", "Upgrade", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForUpgrade(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "diskpools.DiskPoolsClient", "Upgrade", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// UpgradeThenPoll performs Upgrade then polls until it's completed
func (c DiskPoolsClient) UpgradeThenPoll(ctx context.Context, id DiskPoolId) error {
	result, err := c.Upgrade(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Upgrade: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Upgrade: %+v", err)
	}

	return nil
}

// preparerForUpgrade prepares the Upgrade request.
func (c DiskPoolsClient) preparerForUpgrade(ctx context.Context, id DiskPoolId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/upgrade", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForUpgrade sends the Upgrade request. The method will close the
// http.Response Body if it receives an error.
func (c DiskPoolsClient) senderForUpgrade(ctx context.Context, req *http.Request) (future UpgradeResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
