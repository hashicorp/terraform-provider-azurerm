package validateprobe

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ValidateProbeResponse struct {
	HttpResponse *http.Response
	Model        *ValidateProbeOutput
}

// ValidateProbe ...
func (c ValidateProbeClient) ValidateProbe(ctx context.Context, id SubscriptionId, input ValidateProbeInput) (result ValidateProbeResponse, err error) {
	req, err := c.preparerForValidateProbe(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "validateprobe.ValidateProbeClient", "ValidateProbe", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "validateprobe.ValidateProbeClient", "ValidateProbe", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForValidateProbe(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "validateprobe.ValidateProbeClient", "ValidateProbe", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForValidateProbe prepares the ValidateProbe request.
func (c ValidateProbeClient) preparerForValidateProbe(ctx context.Context, id SubscriptionId, input ValidateProbeInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.CDN/validateProbe", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForValidateProbe handles the response to the ValidateProbe request. The method always
// closes the http.Response Body.
func (c ValidateProbeClient) responderForValidateProbe(resp *http.Response) (result ValidateProbeResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
