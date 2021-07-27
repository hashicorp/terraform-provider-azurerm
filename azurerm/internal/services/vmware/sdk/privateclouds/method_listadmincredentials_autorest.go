package privateclouds

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListAdminCredentialsResponse struct {
	HttpResponse *http.Response
	Model        *AdminCredentials
}

// ListAdminCredentials ...
func (c PrivateCloudsClient) ListAdminCredentials(ctx context.Context, id PrivateCloudId) (result ListAdminCredentialsResponse, err error) {
	req, err := c.preparerForListAdminCredentials(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privateclouds.PrivateCloudsClient", "ListAdminCredentials", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "privateclouds.PrivateCloudsClient", "ListAdminCredentials", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListAdminCredentials(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privateclouds.PrivateCloudsClient", "ListAdminCredentials", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListAdminCredentials prepares the ListAdminCredentials request.
func (c PrivateCloudsClient) preparerForListAdminCredentials(ctx context.Context, id PrivateCloudId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listAdminCredentials", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListAdminCredentials handles the response to the ListAdminCredentials request. The method always
// closes the http.Response Body.
func (c PrivateCloudsClient) responderForListAdminCredentials(resp *http.Response) (result ListAdminCredentialsResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
