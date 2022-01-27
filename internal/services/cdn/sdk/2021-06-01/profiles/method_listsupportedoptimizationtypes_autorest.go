package profiles

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListSupportedOptimizationTypesResponse struct {
	HttpResponse *http.Response
	Model        *SupportedOptimizationTypesListResult
}

// ListSupportedOptimizationTypes ...
func (c ProfilesClient) ListSupportedOptimizationTypes(ctx context.Context, id ProfileId) (result ListSupportedOptimizationTypesResponse, err error) {
	req, err := c.preparerForListSupportedOptimizationTypes(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "profiles.ProfilesClient", "ListSupportedOptimizationTypes", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "profiles.ProfilesClient", "ListSupportedOptimizationTypes", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListSupportedOptimizationTypes(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "profiles.ProfilesClient", "ListSupportedOptimizationTypes", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListSupportedOptimizationTypes prepares the ListSupportedOptimizationTypes request.
func (c ProfilesClient) preparerForListSupportedOptimizationTypes(ctx context.Context, id ProfileId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/getSupportedOptimizationTypes", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListSupportedOptimizationTypes handles the response to the ListSupportedOptimizationTypes request. The method always
// closes the http.Response Body.
func (c ProfilesClient) responderForListSupportedOptimizationTypes(resp *http.Response) (result ListSupportedOptimizationTypesResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
