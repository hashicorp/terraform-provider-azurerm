package managedidentity

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type UserAssignedIdentitiesDeleteResponse struct {
	HttpResponse *http.Response
}

// UserAssignedIdentitiesDelete ...
func (c ManagedIdentityClient) UserAssignedIdentitiesDelete(ctx context.Context, id UserAssignedIdentitiesId) (result UserAssignedIdentitiesDeleteResponse, err error) {
	req, err := c.preparerForUserAssignedIdentitiesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "UserAssignedIdentitiesDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "UserAssignedIdentitiesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUserAssignedIdentitiesDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "UserAssignedIdentitiesDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUserAssignedIdentitiesDelete prepares the UserAssignedIdentitiesDelete request.
func (c ManagedIdentityClient) preparerForUserAssignedIdentitiesDelete(ctx context.Context, id UserAssignedIdentitiesId) (*http.Request, error) {
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

// responderForUserAssignedIdentitiesDelete handles the response to the UserAssignedIdentitiesDelete request. The method always
// closes the http.Response Body.
func (c ManagedIdentityClient) responderForUserAssignedIdentitiesDelete(resp *http.Response) (result UserAssignedIdentitiesDeleteResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
