package signalr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type RegenerateKeyResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// RegenerateKey ...
func (c SignalRClient) RegenerateKey(ctx context.Context, id SignalRId, input RegenerateKeyParameters) (result RegenerateKeyResponse, err error) {
	req, err := c.preparerForRegenerateKey(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "RegenerateKey", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForRegenerateKey(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "RegenerateKey", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// RegenerateKeyThenPoll performs RegenerateKey then polls until it's completed
func (c SignalRClient) RegenerateKeyThenPoll(ctx context.Context, id SignalRId, input RegenerateKeyParameters) error {
	result, err := c.RegenerateKey(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing RegenerateKey: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after RegenerateKey: %+v", err)
	}

	return nil
}

// preparerForRegenerateKey prepares the RegenerateKey request.
func (c SignalRClient) preparerForRegenerateKey(ctx context.Context, id SignalRId, input RegenerateKeyParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/regenerateKey", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForRegenerateKey sends the RegenerateKey request. The method will close the
// http.Response Body if it receives an error.
func (c SignalRClient) senderForRegenerateKey(ctx context.Context, req *http.Request) (future RegenerateKeyResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
