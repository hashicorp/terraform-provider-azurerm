package queues

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

type StorageServicePropertiesResponse struct {
	StorageServiceProperties
	autorest.Response
}

// SetServiceProperties gets the properties for this queue
func (client Client) GetServiceProperties(ctx context.Context, accountName string) (result StorageServicePropertiesResponse, err error) {
	if accountName == "" {
		return result, validation.NewError("queues.Client", "GetServiceProperties", "`accountName` cannot be an empty string.")
	}

	req, err := client.GetServicePropertiesPreparer(ctx, accountName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queues.Client", "GetServiceProperties", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetServicePropertiesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "queues.Client", "GetServiceProperties", resp, "Failure sending request")
		return
	}

	result, err = client.GetServicePropertiesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queues.Client", "GetServiceProperties", resp, "Failure responding to request")
		return
	}

	return
}

// GetServicePropertiesPreparer prepares the GetServiceProperties request.
func (client Client) GetServicePropertiesPreparer(ctx context.Context, accountName string) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"comp":    autorest.Encode("path", "properties"),
		"restype": autorest.Encode("path", "service"),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(endpoints.GetQueueEndpoint(client.BaseURI, accountName)),
		autorest.WithPath("/"),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetServicePropertiesSender sends the GetServiceProperties request. The method will close the
// http.Response Body if it receives an error.
func (client Client) GetServicePropertiesSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetServicePropertiesResponder handles the response to the GetServiceProperties request. The method always
// closes the http.Response Body.
func (client Client) GetServicePropertiesResponder(resp *http.Response) (result StorageServicePropertiesResponse, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingXML(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}

	return
}
