package endpoints

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type PurgeContentResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// PurgeContent ...
func (c EndpointsClient) PurgeContent(ctx context.Context, id EndpointId, input PurgeParameters) (result PurgeContentResponse, err error) {
	req, err := c.preparerForPurgeContent(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "endpoints.EndpointsClient", "PurgeContent", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPurgeContent(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "endpoints.EndpointsClient", "PurgeContent", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PurgeContentThenPoll performs PurgeContent then polls until it's completed
func (c EndpointsClient) PurgeContentThenPoll(ctx context.Context, id EndpointId, input PurgeParameters) error {
	result, err := c.PurgeContent(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing PurgeContent: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after PurgeContent: %+v", err)
	}

	return nil
}

// preparerForPurgeContent prepares the PurgeContent request.
func (c EndpointsClient) preparerForPurgeContent(ctx context.Context, id EndpointId, input PurgeParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/purge", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForPurgeContent sends the PurgeContent request. The method will close the
// http.Response Body if it receives an error.
func (c EndpointsClient) senderForPurgeContent(ctx context.Context, req *http.Request) (future PurgeContentResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
