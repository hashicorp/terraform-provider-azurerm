package customdomains

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type EnableCustomHttpsResponse struct {
	HttpResponse *http.Response
	Model        *CustomDomain
}

// EnableCustomHttps ...
func (c CustomDomainsClient) EnableCustomHttps(ctx context.Context, id EndpointCustomDomainId, input CustomDomainHttpsParameters) (result EnableCustomHttpsResponse, err error) {
	req, err := c.preparerForEnableCustomHttps(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customdomains.CustomDomainsClient", "EnableCustomHttps", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "customdomains.CustomDomainsClient", "EnableCustomHttps", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForEnableCustomHttps(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customdomains.CustomDomainsClient", "EnableCustomHttps", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForEnableCustomHttps prepares the EnableCustomHttps request.
func (c CustomDomainsClient) preparerForEnableCustomHttps(ctx context.Context, id EndpointCustomDomainId, input CustomDomainHttpsParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/enableCustomHttps", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForEnableCustomHttps handles the response to the EnableCustomHttps request. The method always
// closes the http.Response Body.
func (c CustomDomainsClient) responderForEnableCustomHttps(resp *http.Response) (result EnableCustomHttpsResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusAccepted, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
