package objectreplicationpolicies

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type CreateOrUpdateResponse struct {
	HttpResponse *http.Response
	Model        *ObjectReplicationPolicy
}

// CreateOrUpdate ...
func (c ObjectReplicationPoliciesClient) CreateOrUpdate(ctx context.Context, id ObjectReplicationPoliciesId, input ObjectReplicationPolicy) (result CreateOrUpdateResponse, err error) {
	req, err := c.preparerForCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "objectreplicationpolicies.ObjectReplicationPoliciesClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "objectreplicationpolicies.ObjectReplicationPoliciesClient", "CreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "objectreplicationpolicies.ObjectReplicationPoliciesClient", "CreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCreateOrUpdate prepares the CreateOrUpdate request.
func (c ObjectReplicationPoliciesClient) preparerForCreateOrUpdate(ctx context.Context, id ObjectReplicationPoliciesId, input ObjectReplicationPolicy) (*http.Request, error) {
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

// responderForCreateOrUpdate handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (c ObjectReplicationPoliciesClient) responderForCreateOrUpdate(resp *http.Response) (result CreateOrUpdateResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
