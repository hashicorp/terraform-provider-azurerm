package customdomains

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type DisableCustomHttpsResponse struct {
	HttpResponse *http.Response
	Model        *CustomDomain
}

// DisableCustomHttps ...
func (c CustomDomainsClient) DisableCustomHttps(ctx context.Context, id EndpointCustomDomainId) (result DisableCustomHttpsResponse, err error) {
	req, err := c.preparerForDisableCustomHttps(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customdomains.CustomDomainsClient", "DisableCustomHttps", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "customdomains.CustomDomainsClient", "DisableCustomHttps", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDisableCustomHttps(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customdomains.CustomDomainsClient", "DisableCustomHttps", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDisableCustomHttps prepares the DisableCustomHttps request.
func (c CustomDomainsClient) preparerForDisableCustomHttps(ctx context.Context, id EndpointCustomDomainId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/disableCustomHttps", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDisableCustomHttps handles the response to the DisableCustomHttps request. The method always
// closes the http.Response Body.
func (c CustomDomainsClient) responderForDisableCustomHttps(resp *http.Response) (result DisableCustomHttpsResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusAccepted, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
