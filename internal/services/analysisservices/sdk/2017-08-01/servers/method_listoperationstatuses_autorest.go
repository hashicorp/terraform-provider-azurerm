package servers

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListOperationStatusesResponse struct {
	HttpResponse *http.Response
	Model        *OperationStatus
}

// ListOperationStatuses ...
func (c ServersClient) ListOperationStatuses(ctx context.Context, id OperationstatuseId) (result ListOperationStatusesResponse, err error) {
	req, err := c.preparerForListOperationStatuses(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servers.ServersClient", "ListOperationStatuses", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "servers.ServersClient", "ListOperationStatuses", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListOperationStatuses(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servers.ServersClient", "ListOperationStatuses", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListOperationStatuses prepares the ListOperationStatuses request.
func (c ServersClient) preparerForListOperationStatuses(ctx context.Context, id OperationstatuseId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListOperationStatuses handles the response to the ListOperationStatuses request. The method always
// closes the http.Response Body.
func (c ServersClient) responderForListOperationStatuses(resp *http.Response) (result ListOperationStatusesResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusAccepted, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
