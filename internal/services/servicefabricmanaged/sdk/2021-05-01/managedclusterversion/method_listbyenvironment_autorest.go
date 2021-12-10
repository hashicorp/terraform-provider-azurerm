package managedclusterversion

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByEnvironmentResponse struct {
	HttpResponse *http.Response
	Model        *[]ManagedClusterCodeVersionResult
}

// ListByEnvironment ...
func (c ManagedClusterVersionClient) ListByEnvironment(ctx context.Context, id EnvironmentId) (result ListByEnvironmentResponse, err error) {
	req, err := c.preparerForListByEnvironment(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusterversion.ManagedClusterVersionClient", "ListByEnvironment", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusterversion.ManagedClusterVersionClient", "ListByEnvironment", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListByEnvironment(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusterversion.ManagedClusterVersionClient", "ListByEnvironment", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListByEnvironment prepares the ListByEnvironment request.
func (c ManagedClusterVersionClient) preparerForListByEnvironment(ctx context.Context, id EnvironmentId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/managedClusterVersions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListByEnvironment handles the response to the ListByEnvironment request. The method always
// closes the http.Response Body.
func (c ManagedClusterVersionClient) responderForListByEnvironment(resp *http.Response) (result ListByEnvironmentResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
