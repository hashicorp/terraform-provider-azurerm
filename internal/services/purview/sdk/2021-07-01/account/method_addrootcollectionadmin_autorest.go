package account

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type AddRootCollectionAdminResponse struct {
	HttpResponse *http.Response
}

// AddRootCollectionAdmin ...
func (c AccountClient) AddRootCollectionAdmin(ctx context.Context, id AccountId, input CollectionAdminUpdate) (result AddRootCollectionAdminResponse, err error) {
	req, err := c.preparerForAddRootCollectionAdmin(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.AccountClient", "AddRootCollectionAdmin", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.AccountClient", "AddRootCollectionAdmin", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAddRootCollectionAdmin(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.AccountClient", "AddRootCollectionAdmin", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAddRootCollectionAdmin prepares the AddRootCollectionAdmin request.
func (c AccountClient) preparerForAddRootCollectionAdmin(ctx context.Context, id AccountId, input CollectionAdminUpdate) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/addRootCollectionAdmin", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAddRootCollectionAdmin handles the response to the AddRootCollectionAdmin request. The method always
// closes the http.Response Body.
func (c AccountClient) responderForAddRootCollectionAdmin(resp *http.Response) (result AddRootCollectionAdminResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
