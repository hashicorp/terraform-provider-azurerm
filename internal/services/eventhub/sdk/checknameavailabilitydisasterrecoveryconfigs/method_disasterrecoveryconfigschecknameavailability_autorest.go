package checknameavailabilitydisasterrecoveryconfigs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type DisasterRecoveryConfigsCheckNameAvailabilityResponse struct {
	HttpResponse *http.Response
	Model        *CheckNameAvailabilityResult
}

// DisasterRecoveryConfigsCheckNameAvailability ...
func (c CheckNameAvailabilityDisasterRecoveryConfigsClient) DisasterRecoveryConfigsCheckNameAvailability(ctx context.Context, id NamespaceId, input CheckNameAvailabilityParameter) (result DisasterRecoveryConfigsCheckNameAvailabilityResponse, err error) {
	req, err := c.preparerForDisasterRecoveryConfigsCheckNameAvailability(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "checknameavailabilitydisasterrecoveryconfigs.CheckNameAvailabilityDisasterRecoveryConfigsClient", "DisasterRecoveryConfigsCheckNameAvailability", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "checknameavailabilitydisasterrecoveryconfigs.CheckNameAvailabilityDisasterRecoveryConfigsClient", "DisasterRecoveryConfigsCheckNameAvailability", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDisasterRecoveryConfigsCheckNameAvailability(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "checknameavailabilitydisasterrecoveryconfigs.CheckNameAvailabilityDisasterRecoveryConfigsClient", "DisasterRecoveryConfigsCheckNameAvailability", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDisasterRecoveryConfigsCheckNameAvailability prepares the DisasterRecoveryConfigsCheckNameAvailability request.
func (c CheckNameAvailabilityDisasterRecoveryConfigsClient) preparerForDisasterRecoveryConfigsCheckNameAvailability(ctx context.Context, id NamespaceId, input CheckNameAvailabilityParameter) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/disasterRecoveryConfigs/checkNameAvailability", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDisasterRecoveryConfigsCheckNameAvailability handles the response to the DisasterRecoveryConfigsCheckNameAvailability request. The method always
// closes the http.Response Body.
func (c CheckNameAvailabilityDisasterRecoveryConfigsClient) responderForDisasterRecoveryConfigsCheckNameAvailability(resp *http.Response) (result DisasterRecoveryConfigsCheckNameAvailabilityResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
