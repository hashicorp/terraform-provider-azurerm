package checknameavailabilitynamespaces

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type NamespacesCheckNameAvailabilityResponse struct {
	HttpResponse *http.Response
	Model        *CheckNameAvailabilityResult
}

// NamespacesCheckNameAvailability ...
func (c CheckNameAvailabilityNamespacesClient) NamespacesCheckNameAvailability(ctx context.Context, id SubscriptionId, input CheckNameAvailabilityParameter) (result NamespacesCheckNameAvailabilityResponse, err error) {
	req, err := c.preparerForNamespacesCheckNameAvailability(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "checknameavailabilitynamespaces.CheckNameAvailabilityNamespacesClient", "NamespacesCheckNameAvailability", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "checknameavailabilitynamespaces.CheckNameAvailabilityNamespacesClient", "NamespacesCheckNameAvailability", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForNamespacesCheckNameAvailability(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "checknameavailabilitynamespaces.CheckNameAvailabilityNamespacesClient", "NamespacesCheckNameAvailability", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForNamespacesCheckNameAvailability prepares the NamespacesCheckNameAvailability request.
func (c CheckNameAvailabilityNamespacesClient) preparerForNamespacesCheckNameAvailability(ctx context.Context, id SubscriptionId, input CheckNameAvailabilityParameter) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.EventHub/checkNameAvailability", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForNamespacesCheckNameAvailability handles the response to the NamespacesCheckNameAvailability request. The method always
// closes the http.Response Body.
func (c CheckNameAvailabilityNamespacesClient) responderForNamespacesCheckNameAvailability(resp *http.Response) (result NamespacesCheckNameAvailabilityResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
