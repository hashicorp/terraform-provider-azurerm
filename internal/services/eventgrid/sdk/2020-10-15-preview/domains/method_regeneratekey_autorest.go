package domains

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type RegenerateKeyResponse struct {
	HttpResponse *http.Response
	Model        *DomainSharedAccessKeys
}

// RegenerateKey ...
func (c DomainsClient) RegenerateKey(ctx context.Context, id DomainId, input DomainRegenerateKeyRequest) (result RegenerateKeyResponse, err error) {
	req, err := c.preparerForRegenerateKey(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "domains.DomainsClient", "RegenerateKey", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "domains.DomainsClient", "RegenerateKey", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRegenerateKey(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "domains.DomainsClient", "RegenerateKey", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRegenerateKey prepares the RegenerateKey request.
func (c DomainsClient) preparerForRegenerateKey(ctx context.Context, id DomainId, input DomainRegenerateKeyRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/regenerateKey", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForRegenerateKey handles the response to the RegenerateKey request. The method always
// closes the http.Response Body.
func (c DomainsClient) responderForRegenerateKey(resp *http.Response) (result RegenerateKeyResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
