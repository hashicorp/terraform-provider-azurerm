package resource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type ServicesCreateOrUpdateResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ServicesCreateOrUpdate ...
func (c ResourceClient) ServicesCreateOrUpdate(ctx context.Context, id ServiceId, input ServicesDescription) (result ServicesCreateOrUpdateResponse, err error) {
	req, err := c.preparerForServicesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "ServicesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForServicesCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "ServicesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ServicesCreateOrUpdateThenPoll performs ServicesCreateOrUpdate then polls until it's completed
func (c ResourceClient) ServicesCreateOrUpdateThenPoll(ctx context.Context, id ServiceId, input ServicesDescription) error {
	result, err := c.ServicesCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ServicesCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ServicesCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForServicesCreateOrUpdate prepares the ServicesCreateOrUpdate request.
func (c ResourceClient) preparerForServicesCreateOrUpdate(ctx context.Context, id ServiceId, input ServicesDescription) (*http.Request, error) {
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

// senderForServicesCreateOrUpdate sends the ServicesCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c ResourceClient) senderForServicesCreateOrUpdate(ctx context.Context, req *http.Request) (future ServicesCreateOrUpdateResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
