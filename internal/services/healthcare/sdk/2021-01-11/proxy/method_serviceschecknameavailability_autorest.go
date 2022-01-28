package proxy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

type ServicesCheckNameAvailabilityResponse struct {
	HttpResponse *http.Response
	Model        *ServicesNameAvailabilityInfo
}

// ServicesCheckNameAvailability ...
func (c ProxyClient) ServicesCheckNameAvailability(ctx context.Context, id commonids.SubscriptionId, input CheckNameAvailabilityParameters) (result ServicesCheckNameAvailabilityResponse, err error) {
	req, err := c.preparerForServicesCheckNameAvailability(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "proxy.ProxyClient", "ServicesCheckNameAvailability", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "proxy.ProxyClient", "ServicesCheckNameAvailability", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForServicesCheckNameAvailability(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "proxy.ProxyClient", "ServicesCheckNameAvailability", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForServicesCheckNameAvailability prepares the ServicesCheckNameAvailability request.
func (c ProxyClient) preparerForServicesCheckNameAvailability(ctx context.Context, id commonids.SubscriptionId, input CheckNameAvailabilityParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.HealthcareApis/checkNameAvailability", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForServicesCheckNameAvailability handles the response to the ServicesCheckNameAvailability request. The method always
// closes the http.Response Body.
func (c ProxyClient) responderForServicesCheckNameAvailability(resp *http.Response) (result ServicesCheckNameAvailabilityResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
