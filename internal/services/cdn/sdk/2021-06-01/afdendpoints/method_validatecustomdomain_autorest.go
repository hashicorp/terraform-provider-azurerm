package afdendpoints

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ValidateCustomDomainResponse struct {
	HttpResponse *http.Response
	Model        *ValidateCustomDomainOutput
}

// ValidateCustomDomain ...
func (c AFDEndpointsClient) ValidateCustomDomain(ctx context.Context, id AfdEndpointId, input ValidateCustomDomainInput) (result ValidateCustomDomainResponse, err error) {
	req, err := c.preparerForValidateCustomDomain(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "afdendpoints.AFDEndpointsClient", "ValidateCustomDomain", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "afdendpoints.AFDEndpointsClient", "ValidateCustomDomain", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForValidateCustomDomain(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "afdendpoints.AFDEndpointsClient", "ValidateCustomDomain", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForValidateCustomDomain prepares the ValidateCustomDomain request.
func (c AFDEndpointsClient) preparerForValidateCustomDomain(ctx context.Context, id AfdEndpointId, input ValidateCustomDomainInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/validateCustomDomain", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForValidateCustomDomain handles the response to the ValidateCustomDomain request. The method always
// closes the http.Response Body.
func (c AFDEndpointsClient) responderForValidateCustomDomain(resp *http.Response) (result ValidateCustomDomainResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
