package disasterrecoveryconfigs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type BreakPairingResponse struct {
	HttpResponse *http.Response
}

// BreakPairing ...
func (c DisasterRecoveryConfigsClient) BreakPairing(ctx context.Context, id DisasterRecoveryConfigId) (result BreakPairingResponse, err error) {
	req, err := c.preparerForBreakPairing(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "disasterrecoveryconfigs.DisasterRecoveryConfigsClient", "BreakPairing", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "disasterrecoveryconfigs.DisasterRecoveryConfigsClient", "BreakPairing", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForBreakPairing(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "disasterrecoveryconfigs.DisasterRecoveryConfigsClient", "BreakPairing", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForBreakPairing prepares the BreakPairing request.
func (c DisasterRecoveryConfigsClient) preparerForBreakPairing(ctx context.Context, id DisasterRecoveryConfigId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/breakPairing", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForBreakPairing handles the response to the BreakPairing request. The method always
// closes the http.Response Body.
func (c DisasterRecoveryConfigsClient) responderForBreakPairing(resp *http.Response) (result BreakPairingResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
