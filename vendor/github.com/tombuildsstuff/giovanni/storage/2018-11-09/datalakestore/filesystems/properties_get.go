package filesystems

import (
	"context"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

type GetPropertiesResponse struct {
	autorest.Response

	// A map of base64-encoded strings to store as user-defined properties with the File System
	// Note that items may only contain ASCII characters in the ISO-8859-1 character set.
	// This automatically gets converted to a comma-separated list of name and
	// value pairs before sending to the API
	Properties map[string]string

	// Is Hierarchical Namespace Enabled?
	NamespaceEnabled bool
}

// GetProperties gets the properties for a Data Lake Store Gen2 FileSystem within a Storage Account
func (client Client) GetProperties(ctx context.Context, accountName string, fileSystemName string) (result GetPropertiesResponse, err error) {
	if accountName == "" {
		return result, validation.NewError("datalakestore.Client", "GetProperties", "`accountName` cannot be an empty string.")
	}
	if fileSystemName == "" {
		return result, validation.NewError("datalakestore.Client", "GetProperties", "`fileSystemName` cannot be an empty string.")
	}

	req, err := client.GetPropertiesPreparer(ctx, accountName, fileSystemName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datalakestore.Client", "GetProperties", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetPropertiesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "datalakestore.Client", "GetProperties", resp, "Failure sending request")
		return
	}

	result, err = client.GetPropertiesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datalakestore.Client", "GetProperties", resp, "Failure responding to request")
	}

	return
}

// GetPropertiesPreparer prepares the GetProperties request.
func (client Client) GetPropertiesPreparer(ctx context.Context, accountName string, fileSystemName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"fileSystemName": autorest.Encode("path", fileSystemName),
	}

	queryParameters := map[string]interface{}{
		"resource": autorest.Encode("query", "filesystem"),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsHead(),
		autorest.WithBaseURL(endpoints.GetDataLakeStoreEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{fileSystemName}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))

	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetPropertiesSender sends the GetProperties request. The method will close the
// http.Response Body if it receives an error.
func (client Client) GetPropertiesSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetPropertiesResponder handles the response to the GetProperties request. The method always
// closes the http.Response Body.
func (client Client) GetPropertiesResponder(resp *http.Response) (result GetPropertiesResponse, err error) {
	if resp != nil && resp.Header != nil {

		propertiesRaw := resp.Header.Get("x-ms-properties")
		var properties *map[string]string
		properties, err = parseProperties(propertiesRaw)
		if err != nil {
			return
		}

		result.Properties = *properties
		result.NamespaceEnabled = strings.EqualFold(resp.Header.Get("x-ms-namespace-enabled"), "tru")
	}

	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}

	return
}
