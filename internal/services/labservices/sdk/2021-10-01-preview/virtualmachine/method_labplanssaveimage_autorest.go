package virtualmachine

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type LabPlansSaveImageResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// LabPlansSaveImage ...
func (c VirtualMachineClient) LabPlansSaveImage(ctx context.Context, id LabPlanId, input SaveImageBody) (result LabPlansSaveImageResponse, err error) {
	req, err := c.preparerForLabPlansSaveImage(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachine.VirtualMachineClient", "LabPlansSaveImage", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForLabPlansSaveImage(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachine.VirtualMachineClient", "LabPlansSaveImage", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// LabPlansSaveImageThenPoll performs LabPlansSaveImage then polls until it's completed
func (c VirtualMachineClient) LabPlansSaveImageThenPoll(ctx context.Context, id LabPlanId, input SaveImageBody) error {
	result, err := c.LabPlansSaveImage(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing LabPlansSaveImage: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after LabPlansSaveImage: %+v", err)
	}

	return nil
}

// preparerForLabPlansSaveImage prepares the LabPlansSaveImage request.
func (c VirtualMachineClient) preparerForLabPlansSaveImage(ctx context.Context, id LabPlanId, input SaveImageBody) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/saveImage", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForLabPlansSaveImage sends the LabPlansSaveImage request. The method will close the
// http.Response Body if it receives an error.
func (c VirtualMachineClient) senderForLabPlansSaveImage(ctx context.Context, req *http.Request) (future LabPlansSaveImageResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
