package alerts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListExternalResponse struct {
	HttpResponse *http.Response
	Model        *AlertsResult
}

// ListExternal ...
func (c AlertsClient) ListExternal(ctx context.Context, id ExternalCloudProviderTypeId) (result ListExternalResponse, err error) {
	req, err := c.preparerForListExternal(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alerts.AlertsClient", "ListExternal", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "alerts.AlertsClient", "ListExternal", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListExternal(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alerts.AlertsClient", "ListExternal", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListExternal prepares the ListExternal request.
func (c AlertsClient) preparerForListExternal(ctx context.Context, id ExternalCloudProviderTypeId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/alerts", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListExternal handles the response to the ListExternal request. The method always
// closes the http.Response Body.
func (c AlertsClient) responderForListExternal(resp *http.Response) (result ListExternalResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
