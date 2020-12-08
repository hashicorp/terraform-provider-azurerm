package queues

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

// SetServiceProperties sets the properties for this queue
func (client Client) SetServiceProperties(ctx context.Context, accountName string, properties StorageServiceProperties) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("queues.Client", "SetServiceProperties", "`accountName` cannot be an empty string.")
	}

	req, err := client.SetServicePropertiesPreparer(ctx, accountName, properties)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queues.Client", "SetServiceProperties", nil, "Failure preparing request")
		return
	}

	resp, err := client.SetServicePropertiesSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "queues.Client", "SetServiceProperties", resp, "Failure sending request")
		return
	}

	result, err = client.SetServicePropertiesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queues.Client", "SetServiceProperties", resp, "Failure responding to request")
		return
	}

	return
}

// SetServicePropertiesPreparer prepares the SetServiceProperties request.
func (client Client) SetServicePropertiesPreparer(ctx context.Context, accountName string, properties StorageServiceProperties) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"comp":    autorest.Encode("path", "properties"),
		"restype": autorest.Encode("path", "service"),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(endpoints.GetQueueEndpoint(client.BaseURI, accountName)),
		autorest.WithPath("/"),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithXML(properties),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// SetServicePropertiesSender sends the SetServiceProperties request. The method will close the
// http.Response Body if it receives an error.
func (client Client) SetServicePropertiesSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// SetServicePropertiesResponder handles the response to the SetServiceProperties request. The method always
// closes the http.Response Body.
func (client Client) SetServicePropertiesResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusAccepted),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
