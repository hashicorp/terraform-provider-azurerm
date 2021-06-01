package disasterrecoveryconfigs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type FailOverResponse struct {
	HttpResponse *http.Response
}

// FailOver ...
func (c DisasterRecoveryConfigsClient) FailOver(ctx context.Context, id DisasterRecoveryConfigId) (result FailOverResponse, err error) {
	req, err := c.preparerForFailOver(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "disasterrecoveryconfigs.DisasterRecoveryConfigsClient", "FailOver", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "disasterrecoveryconfigs.DisasterRecoveryConfigsClient", "FailOver", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForFailOver(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "disasterrecoveryconfigs.DisasterRecoveryConfigsClient", "FailOver", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForFailOver prepares the FailOver request.
func (c DisasterRecoveryConfigsClient) preparerForFailOver(ctx context.Context, id DisasterRecoveryConfigId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/failover", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForFailOver handles the response to the FailOver request. The method always
// closes the http.Response Body.
func (c DisasterRecoveryConfigsClient) responderForFailOver(resp *http.Response) (result FailOverResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
