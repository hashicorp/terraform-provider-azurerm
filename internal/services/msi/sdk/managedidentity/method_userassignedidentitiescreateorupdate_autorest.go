package managedidentity

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type UserAssignedIdentitiesCreateOrUpdateResponse struct {
	HttpResponse *http.Response
	Model        *Identity
}

// UserAssignedIdentitiesCreateOrUpdate ...
func (c ManagedIdentityClient) UserAssignedIdentitiesCreateOrUpdate(ctx context.Context, id UserAssignedIdentitiesId, input Identity) (result UserAssignedIdentitiesCreateOrUpdateResponse, err error) {
	req, err := c.preparerForUserAssignedIdentitiesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "UserAssignedIdentitiesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "UserAssignedIdentitiesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUserAssignedIdentitiesCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "UserAssignedIdentitiesCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUserAssignedIdentitiesCreateOrUpdate prepares the UserAssignedIdentitiesCreateOrUpdate request.
func (c ManagedIdentityClient) preparerForUserAssignedIdentitiesCreateOrUpdate(ctx context.Context, id UserAssignedIdentitiesId, input Identity) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForUserAssignedIdentitiesCreateOrUpdate handles the response to the UserAssignedIdentitiesCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c ManagedIdentityClient) responderForUserAssignedIdentitiesCreateOrUpdate(resp *http.Response) (result UserAssignedIdentitiesCreateOrUpdateResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
