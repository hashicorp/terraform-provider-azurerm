package files

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

type SetPropertiesInput struct {
	// Resizes a file to the specified size.
	// If the specified byte value is less than the current size of the file,
	// then all ranges above the specified byte value are cleared.
	ContentLength *int64

	// Modifies the cache control string for the file.
	// If this property is not specified on the request, then the property will be cleared for the file.
	// Subsequent calls to Get File Properties will not return this property,
	// unless it is explicitly set on the file again.
	ContentControl *string

	// Sets the file’s Content-Disposition header.
	// If this property is not specified on the request, then the property will be cleared for the file.
	// Subsequent calls to Get File Properties will not return this property,
	// unless it is explicitly set on the file again.
	ContentDisposition *string

	// Sets the file's content encoding.
	// If this property is not specified on the request, then the property will be cleared for the file.
	// Subsequent calls to Get File Properties will not return this property,
	// unless it is explicitly set on the file again.
	ContentEncoding *string

	// Sets the file's content language.
	// If this property is not specified on the request, then the property will be cleared for the file.
	// Subsequent calls to Get File Properties will not return this property,
	// unless it is explicitly set on the file again.
	ContentLanguage *string

	// Sets the file's MD5 hash.
	// If this property is not specified on the request, then the property will be cleared for the file.
	// Subsequent calls to Get File Properties will not return this property,
	// unless it is explicitly set on the file again.
	ContentMD5 *string

	// Sets the file's content type.
	// If this property is not specified on the request, then the property will be cleared for the file.
	// Subsequent calls to Get File Properties will not return this property,
	// unless it is explicitly set on the file again.
	ContentType *string

	// The time at which this file was created at - if omitted, this'll be set to "now"
	// This maps to the `x-ms-file-creation-time` field.
	CreatedAt *time.Time

	// The time at which this file was last modified - if omitted, this'll be set to "now"
	// This maps to the `x-ms-file-last-write-time` field.
	LastModified *time.Time
}

// SetProperties sets the specified properties on the specified File
func (client Client) SetProperties(ctx context.Context, accountName, shareName, path, fileName string, input SetPropertiesInput) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("files.Client", "SetProperties", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("files.Client", "SetProperties", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("files.Client", "SetProperties", "`shareName` must be a lower-cased string.")
	}
	if fileName == "" {
		return result, validation.NewError("files.Client", "SetProperties", "`fileName` cannot be an empty string.")
	}

	req, err := client.SetPropertiesPreparer(ctx, accountName, shareName, path, fileName, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "files.Client", "SetProperties", nil, "Failure preparing request")
		return
	}

	resp, err := client.SetPropertiesSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "files.Client", "SetProperties", resp, "Failure sending request")
		return
	}

	result, err = client.SetPropertiesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "files.Client", "SetProperties", resp, "Failure responding to request")
		return
	}

	return
}

// SetPropertiesPreparer prepares the SetProperties request.
func (client Client) SetPropertiesPreparer(ctx context.Context, accountName, shareName, path, fileName string, input SetPropertiesInput) (*http.Request, error) {
	if path != "" {
		path = fmt.Sprintf("%s/", path)
	}
	pathParameters := map[string]interface{}{
		"shareName": autorest.Encode("path", shareName),
		"directory": autorest.Encode("path", path),
		"fileName":  autorest.Encode("path", fileName),
	}

	var coalesceDate = func(input *time.Time, defaultVal string) string {
		if input == nil {
			return defaultVal
		}

		return input.Format(time.RFC1123)
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
		"x-ms-type":    "file",

		"x-ms-file-permission":      "inherit", // TODO: expose this in future
		"x-ms-file-attributes":      "None",    // TODO: expose this in future
		"x-ms-file-creation-time":   coalesceDate(input.CreatedAt, "now"),
		"x-ms-file-last-write-time": coalesceDate(input.LastModified, "now"),
	}

	if input.ContentControl != nil {
		headers["x-ms-cache-control"] = *input.ContentControl
	}
	if input.ContentDisposition != nil {
		headers["x-ms-content-disposition"] = *input.ContentDisposition
	}
	if input.ContentEncoding != nil {
		headers["x-ms-content-encoding"] = *input.ContentEncoding
	}
	if input.ContentLanguage != nil {
		headers["x-ms-content-language"] = *input.ContentLanguage
	}
	if input.ContentLength != nil {
		headers["x-ms-content-length"] = *input.ContentLength
	}
	if input.ContentMD5 != nil {
		headers["x-ms-content-md5"] = *input.ContentMD5
	}
	if input.ContentType != nil {
		headers["x-ms-content-type"] = *input.ContentType
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(endpoints.GetFileEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{shareName}/{directory}{fileName}", pathParameters),
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
		azure.WithErrorUnlessStatusCode(http.StatusCreated),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
