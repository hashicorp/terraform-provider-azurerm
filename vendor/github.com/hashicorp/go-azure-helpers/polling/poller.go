package polling

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type LongRunningPoller struct {
	// HttpResponse is the latest HTTP Response
	HttpResponse *http.Response

	future *azure.Future
	ctx    context.Context
	client autorest.Client
}

// NewLongRunningPollerFromResponse creates a new LongRunningPoller from the HTTP Response
func NewLongRunningPollerFromResponse(ctx context.Context, resp *http.Response, client autorest.Client) (LongRunningPoller, error) {
	poller := LongRunningPoller{
		ctx:    ctx,
		client: client,
	}
	future, err := azure.NewFutureFromResponse(resp)
	if err != nil {
		return poller, err
	}
	poller.future = &future
	return poller, nil
}

// PollUntilDone polls until this Long Running Poller is completed
func (fw *LongRunningPoller) PollUntilDone() error {
	err := fw.future.WaitForCompletionRef(fw.ctx, fw.client)
	fw.HttpResponse = fw.future.Response()
	return err
}
