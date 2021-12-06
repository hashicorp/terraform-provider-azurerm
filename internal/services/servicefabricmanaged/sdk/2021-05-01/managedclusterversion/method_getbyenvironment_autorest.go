package managedclusterversion

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type GetByEnvironmentResponse struct {
	HttpResponse *http.Response
	Model        *ManagedClusterCodeVersionResult
}

// GetByEnvironment ...
func (c ManagedClusterVersionClient) GetByEnvironment(ctx context.Context, id EnvironmentManagedClusterVersionId) (result GetByEnvironmentResponse, err error) {
	req, err := c.preparerForGetByEnvironment(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusterversion.ManagedClusterVersionClient", "GetByEnvironment", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusterversion.ManagedClusterVersionClient", "GetByEnvironment", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetByEnvironment(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusterversion.ManagedClusterVersionClient", "GetByEnvironment", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetByEnvironment prepares the GetByEnvironment request.
func (c ManagedClusterVersionClient) preparerForGetByEnvironment(ctx context.Context, id EnvironmentManagedClusterVersionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetByEnvironment handles the response to the GetByEnvironment request. The method always
// closes the http.Response Body.
func (c ManagedClusterVersionClient) responderForGetByEnvironment(resp *http.Response) (result GetByEnvironmentResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
