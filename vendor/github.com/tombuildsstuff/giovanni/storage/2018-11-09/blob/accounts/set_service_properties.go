package accounts

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
	"net/http"
)

// SetServicePropertiesSender sends the SetServiceProperties request. The method will close the
// http.Response Body if it receives an error.
func (client Client) SetServicePropertiesSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// SetServicePropertiesPreparer prepares the SetServiceProperties request.
func (client Client) SetServicePropertiesPreparer(ctx context.Context, accountName string, input StorageServiceProperties) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"restype": "service",
		"comp":    "properties",
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPut(),
		autorest.WithBaseURL(endpoints.GetBlobEndpoint(client.BaseURI, accountName)),
		autorest.WithHeaders(headers),
		autorest.WithXML(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// SetServicePropertiesResponder handles the response to the SetServiceProperties request. The method always
// closes the http.Response Body.
func (client Client) SetServicePropertiesResponder(resp *http.Response) (result SetServicePropertiesResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusAccepted),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

func (client Client) SetServiceProperties(ctx context.Context, accountName string, input StorageServiceProperties) (result SetServicePropertiesResult, err error) {
	if accountName == "" {
		return result, validation.NewError("accounts.Client", "SetServiceProperties", "`accountName` cannot be an empty string.")
	}

	req, err := client.SetServicePropertiesPreparer(ctx, accountName, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.Client", "SetServiceProperties", nil, "Failure preparing request")
		return
	}

	resp, err := client.SetServicePropertiesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "accounts.Client", "SetServiceProperties", resp, "Failure sending request")
		return
	}

	result, err = client.SetServicePropertiesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.Client", "SetServiceProperties", resp, "Failure responding to request")
		return
	}

	return
}
