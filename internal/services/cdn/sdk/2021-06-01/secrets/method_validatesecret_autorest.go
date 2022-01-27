package secrets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ValidateSecretResponse struct {
	HttpResponse *http.Response
	Model        *ValidateSecretOutput
}

// ValidateSecret ...
func (c SecretsClient) ValidateSecret(ctx context.Context, id SubscriptionId, input ValidateSecretInput) (result ValidateSecretResponse, err error) {
	req, err := c.preparerForValidateSecret(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "secrets.SecretsClient", "ValidateSecret", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "secrets.SecretsClient", "ValidateSecret", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForValidateSecret(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "secrets.SecretsClient", "ValidateSecret", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForValidateSecret prepares the ValidateSecret request.
func (c SecretsClient) preparerForValidateSecret(ctx context.Context, id SubscriptionId, input ValidateSecretInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.CDN/validateSecret", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForValidateSecret handles the response to the ValidateSecret request. The method always
// closes the http.Response Body.
func (c SecretsClient) responderForValidateSecret(resp *http.Response) (result ValidateSecretResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
