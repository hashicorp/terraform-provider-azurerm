package virtualmachine

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type ResetPasswordResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ResetPassword ...
func (c VirtualMachineClient) ResetPassword(ctx context.Context, id VirtualMachineId, input ResetPasswordBody) (result ResetPasswordResponse, err error) {
	req, err := c.preparerForResetPassword(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachine.VirtualMachineClient", "ResetPassword", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForResetPassword(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachine.VirtualMachineClient", "ResetPassword", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ResetPasswordThenPoll performs ResetPassword then polls until it's completed
func (c VirtualMachineClient) ResetPasswordThenPoll(ctx context.Context, id VirtualMachineId, input ResetPasswordBody) error {
	result, err := c.ResetPassword(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ResetPassword: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ResetPassword: %+v", err)
	}

	return nil
}

// preparerForResetPassword prepares the ResetPassword request.
func (c VirtualMachineClient) preparerForResetPassword(ctx context.Context, id VirtualMachineId, input ResetPasswordBody) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/resetPassword", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForResetPassword sends the ResetPassword request. The method will close the
// http.Response Body if it receives an error.
func (c VirtualMachineClient) senderForResetPassword(ctx context.Context, req *http.Request) (future ResetPasswordResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
