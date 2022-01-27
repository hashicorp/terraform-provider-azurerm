package endpoints

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type LoadContentResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// LoadContent ...
func (c EndpointsClient) LoadContent(ctx context.Context, id EndpointId, input LoadParameters) (result LoadContentResponse, err error) {
	req, err := c.preparerForLoadContent(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "endpoints.EndpointsClient", "LoadContent", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForLoadContent(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "endpoints.EndpointsClient", "LoadContent", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// LoadContentThenPoll performs LoadContent then polls until it's completed
func (c EndpointsClient) LoadContentThenPoll(ctx context.Context, id EndpointId, input LoadParameters) error {
	result, err := c.LoadContent(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing LoadContent: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after LoadContent: %+v", err)
	}

	return nil
}

// preparerForLoadContent prepares the LoadContent request.
func (c EndpointsClient) preparerForLoadContent(ctx context.Context, id EndpointId, input LoadParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/load", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForLoadContent sends the LoadContent request. The method will close the
// http.Response Body if it receives an error.
func (c EndpointsClient) senderForLoadContent(ctx context.Context, req *http.Request) (future LoadContentResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
