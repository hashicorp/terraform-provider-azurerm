package configurationstores

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type CreateResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Create ...
func (c ConfigurationStoresClient) Create(ctx context.Context, id ConfigurationStoreId, input ConfigurationStore) (result CreateResponse, err error) {
	req, err := c.preparerForCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationstores.ConfigurationStoresClient", "Create", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCreate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationstores.ConfigurationStoresClient", "Create", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CreateThenPoll performs Create then polls until it's completed
func (c ConfigurationStoresClient) CreateThenPoll(ctx context.Context, id ConfigurationStoreId, input ConfigurationStore) error {
	result, err := c.Create(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Create: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Create: %+v", err)
	}

	return nil
}

// preparerForCreate prepares the Create request.
func (c ConfigurationStoresClient) preparerForCreate(ctx context.Context, id ConfigurationStoreId, input ConfigurationStore) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForCreate sends the Create request. The method will close the
// http.Response Body if it receives an error.
func (c ConfigurationStoresClient) senderForCreate(ctx context.Context, req *http.Request) (future CreateResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
