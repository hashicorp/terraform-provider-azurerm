package filesystems

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

type SetPropertiesInput struct {
	// A map of base64-encoded strings to store as user-defined properties with the File System
	// Note that items may only contain ASCII characters in the ISO-8859-1 character set.
	// This automatically gets converted to a comma-separated list of name and
	// value pairs before sending to the API
	Properties map[string]string

	// Optional - A date and time value.
	// Specify this header to perform the operation only if the resource has been modified since the specified date and time.
	IfModifiedSince *string

	// Optional - A date and time value.
	// Specify this header to perform the operation only if the resource has not been modified since the specified date and time.
	IfUnmodifiedSince *string
}

// SetProperties sets the Properties for a Data Lake Store Gen2 FileSystem within a Storage Account
func (client Client) SetProperties(ctx context.Context, accountName string, fileSystemName string, input SetPropertiesInput) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("datalakestore.Client", "SetProperties", "`accountName` cannot be an empty string.")
	}
	if fileSystemName == "" {
		return result, validation.NewError("datalakestore.Client", "SetProperties", "`fileSystemName` cannot be an empty string.")
	}

	req, err := client.SetPropertiesPreparer(ctx, accountName, fileSystemName, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datalakestore.Client", "SetProperties", nil, "Failure preparing request")
		return
	}

	resp, err := client.SetPropertiesSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "datalakestore.Client", "SetProperties", resp, "Failure sending request")
		return
	}

	result, err = client.SetPropertiesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datalakestore.Client", "SetProperties", resp, "Failure responding to request")
	}

	return
}

// SetPropertiesPreparer prepares the SetProperties request.
func (client Client) SetPropertiesPreparer(ctx context.Context, accountName string, fileSystemName string, input SetPropertiesInput) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"fileSystemName": autorest.Encode("path", fileSystemName),
	}

	queryParameters := map[string]interface{}{
		"resource": autorest.Encode("query", "filesystem"),
	}

	headers := map[string]interface{}{
		"x-ms-properties": buildProperties(input.Properties),
		"x-ms-version":    APIVersion,
	}

	if input.IfModifiedSince != nil {
		headers["If-Modified-Since"] = *input.IfModifiedSince
	}
	if input.IfUnmodifiedSince != nil {
		headers["If-Unmodified-Since"] = *input.IfUnmodifiedSince
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPatch(),
		autorest.WithBaseURL(endpoints.GetDataLakeStoreEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{fileSystemName}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))

	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// SetPropertiesSender sends the SetProperties request. The method will close the
// http.Response Body if it receives an error.
func (client Client) SetPropertiesSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// SetPropertiesResponder handles the response to the SetProperties request. The method always
// closes the http.Response Body.
func (client Client) SetPropertiesResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
