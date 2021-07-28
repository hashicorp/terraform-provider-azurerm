package servers

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type GetDetailsResponse struct {
	HttpResponse *http.Response
	Model        *AnalysisServicesServer
}

// GetDetails ...
func (c ServersClient) GetDetails(ctx context.Context, id ServerId) (result GetDetailsResponse, err error) {
	req, err := c.preparerForGetDetails(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servers.ServersClient", "GetDetails", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "servers.ServersClient", "GetDetails", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetDetails(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servers.ServersClient", "GetDetails", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetDetails prepares the GetDetails request.
func (c ServersClient) preparerForGetDetails(ctx context.Context, id ServerId) (*http.Request, error) {
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

// responderForGetDetails handles the response to the GetDetails request. The method always
// closes the http.Response Body.
func (c ServersClient) responderForGetDetails(resp *http.Response) (result GetDetailsResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
