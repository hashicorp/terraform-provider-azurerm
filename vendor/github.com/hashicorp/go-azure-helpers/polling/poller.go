package polling

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/lang/response"
)

type LongRunningPoller struct {
	// HttpResponse is the latest HTTP Response
	HttpResponse *http.Response

	future *azure.Future
	ctx    context.Context
	client autorest.Client
	method string
}

// NewLongRunningPollerFromResponse creates a new LongRunningPoller from the HTTP Response
// this is deprecated and replaced by NewPollerFromResponse. Can be removed once all the
// embedded SDKs have been removed.
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

// NewPollerFromResponse creates a new LongRunningPoller from the HTTP Response
func NewPollerFromResponse(ctx context.Context, resp *http.Response, client autorest.Client, method string) (LongRunningPoller, error) {
	poller := LongRunningPoller{
		ctx:    ctx,
		client: client,
		method: method,
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
	if fw.future == nil {
		return fmt.Errorf("internal error: cannot poll on a nil-future")
	}

	err := fw.future.WaitForCompletionRef(fw.ctx, fw.client)
	fw.HttpResponse = fw.future.Response()
	if strings.EqualFold(fw.method, "DELETE") {
		if response.WasNotFound(fw.HttpResponse) {
			err = nil
		}
	}
	return err
}
