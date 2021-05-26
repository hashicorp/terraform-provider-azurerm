package configurationstores

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type GetResponse struct {
	HttpResponse *http.Response
	Model        *ConfigurationStore
}

// Get ...
func (c ConfigurationStoresClient) Get(ctx context.Context, id ConfigurationStoreId) (result GetResponse, err error) {
	req, err := c.GetPreparer(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationstores.ConfigurationStoresClient", "Get", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationstores.ConfigurationStoresClient", "Get", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.GetResponder(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationstores.ConfigurationStoresClient", "Get", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (c ConfigurationStoresClient) GetPreparer(ctx context.Context, id ConfigurationStoreId) (*http.Request, error) {
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

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (c ConfigurationStoresClient) GetResponder(resp *http.Response) (result GetResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
