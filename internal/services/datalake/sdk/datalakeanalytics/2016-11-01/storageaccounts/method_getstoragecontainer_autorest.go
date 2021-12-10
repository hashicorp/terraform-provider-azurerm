package storageaccounts

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type GetStorageContainerResponse struct {
	HttpResponse *http.Response
	Model        *StorageContainer
}

// GetStorageContainer ...
func (c StorageAccountsClient) GetStorageContainer(ctx context.Context, id ContainerId) (result GetStorageContainerResponse, err error) {
	req, err := c.preparerForGetStorageContainer(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "GetStorageContainer", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "GetStorageContainer", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetStorageContainer(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "GetStorageContainer", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetStorageContainer prepares the GetStorageContainer request.
func (c StorageAccountsClient) preparerForGetStorageContainer(ctx context.Context, id ContainerId) (*http.Request, error) {
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

// responderForGetStorageContainer handles the response to the GetStorageContainer request. The method always
// closes the http.Response Body.
func (c StorageAccountsClient) responderForGetStorageContainer(resp *http.Response) (result GetStorageContainerResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
