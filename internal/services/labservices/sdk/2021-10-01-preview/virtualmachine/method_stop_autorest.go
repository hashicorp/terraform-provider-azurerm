package virtualmachine

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type StopResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Stop ...
func (c VirtualMachineClient) Stop(ctx context.Context, id VirtualMachineId) (result StopResponse, err error) {
	req, err := c.preparerForStop(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachine.VirtualMachineClient", "Stop", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForStop(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachine.VirtualMachineClient", "Stop", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// StopThenPoll performs Stop then polls until it's completed
func (c VirtualMachineClient) StopThenPoll(ctx context.Context, id VirtualMachineId) error {
	result, err := c.Stop(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Stop: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Stop: %+v", err)
	}

	return nil
}

// preparerForStop prepares the Stop request.
func (c VirtualMachineClient) preparerForStop(ctx context.Context, id VirtualMachineId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/stop", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForStop sends the Stop request. The method will close the
// http.Response Body if it receives an error.
func (c VirtualMachineClient) senderForStop(ctx context.Context, req *http.Request) (future StopResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
