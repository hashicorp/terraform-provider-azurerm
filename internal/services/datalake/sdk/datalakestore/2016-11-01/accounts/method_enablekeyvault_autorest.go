package accounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type EnableKeyVaultResponse struct {
	HttpResponse *http.Response
}

// EnableKeyVault ...
func (c AccountsClient) EnableKeyVault(ctx context.Context, id AccountId) (result EnableKeyVaultResponse, err error) {
	req, err := c.preparerForEnableKeyVault(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "EnableKeyVault", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "EnableKeyVault", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForEnableKeyVault(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "EnableKeyVault", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForEnableKeyVault prepares the EnableKeyVault request.
func (c AccountsClient) preparerForEnableKeyVault(ctx context.Context, id AccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/enableKeyVault", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForEnableKeyVault handles the response to the EnableKeyVault request. The method always
// closes the http.Response Body.
func (c AccountsClient) responderForEnableKeyVault(resp *http.Response) (result EnableKeyVaultResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
