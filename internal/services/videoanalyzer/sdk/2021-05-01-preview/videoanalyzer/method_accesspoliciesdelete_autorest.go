package videoanalyzer

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type AccessPoliciesDeleteResponse struct {
	HttpResponse *http.Response
}

// AccessPoliciesDelete ...
func (c VideoAnalyzerClient) AccessPoliciesDelete(ctx context.Context, id AccessPoliciesId) (result AccessPoliciesDeleteResponse, err error) {
	req, err := c.preparerForAccessPoliciesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "AccessPoliciesDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "AccessPoliciesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAccessPoliciesDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "AccessPoliciesDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAccessPoliciesDelete prepares the AccessPoliciesDelete request.
func (c VideoAnalyzerClient) preparerForAccessPoliciesDelete(ctx context.Context, id AccessPoliciesId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAccessPoliciesDelete handles the response to the AccessPoliciesDelete request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForAccessPoliciesDelete(resp *http.Response) (result AccessPoliciesDeleteResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
