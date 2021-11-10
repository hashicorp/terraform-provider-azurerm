package cognitiveservicesaccounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type CheckSkuAvailabilityResponse struct {
	HttpResponse *http.Response
	Model        *SkuAvailabilityListResult
}

// CheckSkuAvailability ...
func (c CognitiveServicesAccountsClient) CheckSkuAvailability(ctx context.Context, id LocationId, input CheckSkuAvailabilityParameter) (result CheckSkuAvailabilityResponse, err error) {
	req, err := c.preparerForCheckSkuAvailability(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "CheckSkuAvailability", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "CheckSkuAvailability", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCheckSkuAvailability(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "CheckSkuAvailability", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCheckSkuAvailability prepares the CheckSkuAvailability request.
func (c CognitiveServicesAccountsClient) preparerForCheckSkuAvailability(ctx context.Context, id LocationId, input CheckSkuAvailabilityParameter) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/checkSkuAvailability", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCheckSkuAvailability handles the response to the CheckSkuAvailability request. The method always
// closes the http.Response Body.
func (c CognitiveServicesAccountsClient) responderForCheckSkuAvailability(resp *http.Response) (result CheckSkuAvailabilityResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
