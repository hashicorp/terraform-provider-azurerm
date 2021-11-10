package authorizationrulesdisasterrecoveryconfigs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type DisasterRecoveryConfigsListKeysResponse struct {
	HttpResponse *http.Response
	Model        *AccessKeys
}

// DisasterRecoveryConfigsListKeys ...
func (c AuthorizationRulesDisasterRecoveryConfigsClient) DisasterRecoveryConfigsListKeys(ctx context.Context, id DisasterRecoveryConfigAuthorizationRuleId) (result DisasterRecoveryConfigsListKeysResponse, err error) {
	req, err := c.preparerForDisasterRecoveryConfigsListKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesdisasterrecoveryconfigs.AuthorizationRulesDisasterRecoveryConfigsClient", "DisasterRecoveryConfigsListKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesdisasterrecoveryconfigs.AuthorizationRulesDisasterRecoveryConfigsClient", "DisasterRecoveryConfigsListKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDisasterRecoveryConfigsListKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesdisasterrecoveryconfigs.AuthorizationRulesDisasterRecoveryConfigsClient", "DisasterRecoveryConfigsListKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDisasterRecoveryConfigsListKeys prepares the DisasterRecoveryConfigsListKeys request.
func (c AuthorizationRulesDisasterRecoveryConfigsClient) preparerForDisasterRecoveryConfigsListKeys(ctx context.Context, id DisasterRecoveryConfigAuthorizationRuleId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listKeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDisasterRecoveryConfigsListKeys handles the response to the DisasterRecoveryConfigsListKeys request. The method always
// closes the http.Response Body.
func (c AuthorizationRulesDisasterRecoveryConfigsClient) responderForDisasterRecoveryConfigsListKeys(resp *http.Response) (result DisasterRecoveryConfigsListKeysResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
