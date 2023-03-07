package azuresdkhacks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotecteditems"
)

type DeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Delete ...
func (c ReplicationProtectedItemsClient) Delete(ctx context.Context, id replicationprotecteditems.ReplicationProtectedItemId, input DisableProtectionInput) (result DeleteOperationResponse, err error) {
	req, err := c.preparerForDelete(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotecteditems.ReplicationProtectedItemsClient", "Delete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotecteditems.ReplicationProtectedItemsClient", "Delete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DeleteThenPoll performs Delete then polls until it's completed
func (c ReplicationProtectedItemsClient) DeleteThenPoll(ctx context.Context, id replicationprotecteditems.ReplicationProtectedItemId, input DisableProtectionInput) error {
	result, err := c.Delete(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Delete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Delete: %+v", err)
	}

	return nil
}

// preparerForDelete prepares the Delete request.
func (c ReplicationProtectedItemsClient) preparerForDelete(ctx context.Context, id replicationprotecteditems.ReplicationProtectedItemId, input DisableProtectionInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.BaseUri),
		autorest.WithPath(fmt.Sprintf("%s/remove", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForDelete sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationProtectedItemsClient) senderForDelete(ctx context.Context, req *http.Request) (future DeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
