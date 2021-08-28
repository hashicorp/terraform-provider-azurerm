package managedidentity

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type UserAssignedIdentitiesUpdateResponse struct {
	HttpResponse *http.Response
	Model        *Identity
}

// UserAssignedIdentitiesUpdate ...
func (c ManagedIdentityClient) UserAssignedIdentitiesUpdate(ctx context.Context, id UserAssignedIdentitiesId, input IdentityUpdate) (result UserAssignedIdentitiesUpdateResponse, err error) {
	req, err := c.preparerForUserAssignedIdentitiesUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "UserAssignedIdentitiesUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "UserAssignedIdentitiesUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUserAssignedIdentitiesUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "UserAssignedIdentitiesUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUserAssignedIdentitiesUpdate prepares the UserAssignedIdentitiesUpdate request.
func (c ManagedIdentityClient) preparerForUserAssignedIdentitiesUpdate(ctx context.Context, id UserAssignedIdentitiesId, input IdentityUpdate) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForUserAssignedIdentitiesUpdate handles the response to the UserAssignedIdentitiesUpdate request. The method always
// closes the http.Response Body.
func (c ManagedIdentityClient) responderForUserAssignedIdentitiesUpdate(resp *http.Response) (result UserAssignedIdentitiesUpdateResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
