package alerts

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type DismissResponse struct {
	HttpResponse *http.Response
	Model        *Alert
}

// Dismiss ...
func (c AlertsClient) Dismiss(ctx context.Context, id ScopedAlertId, input DismissAlertPayload) (result DismissResponse, err error) {
	req, err := c.preparerForDismiss(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alerts.AlertsClient", "Dismiss", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "alerts.AlertsClient", "Dismiss", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDismiss(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alerts.AlertsClient", "Dismiss", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDismiss prepares the Dismiss request.
func (c AlertsClient) preparerForDismiss(ctx context.Context, id ScopedAlertId, input DismissAlertPayload) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDismiss handles the response to the Dismiss request. The method always
// closes the http.Response Body.
func (c AlertsClient) responderForDismiss(resp *http.Response) (result DismissResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
