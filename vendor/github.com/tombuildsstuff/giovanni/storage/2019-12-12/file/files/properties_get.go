package files

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
	"github.com/tombuildsstuff/giovanni/storage/internal/metadata"
)

type GetResult struct {
	autorest.Response

	CacheControl          string
	ContentDisposition    string
	ContentEncoding       string
	ContentLanguage       string
	ContentLength         *int64
	ContentMD5            string
	ContentType           string
	CopyID                string
	CopyStatus            string
	CopySource            string
	CopyProgress          string
	CopyStatusDescription string
	CopyCompletionTime    string
	Encrypted             bool

	MetaData map[string]string
}

// GetProperties returns the Properties for the specified file
func (client Client) GetProperties(ctx context.Context, accountName, shareName, path, fileName string) (result GetResult, err error) {
	if accountName == "" {
		return result, validation.NewError("files.Client", "GetProperties", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("files.Client", "GetProperties", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("files.Client", "GetProperties", "`shareName` must be a lower-cased string.")
	}
	if fileName == "" {
		return result, validation.NewError("files.Client", "GetProperties", "`fileName` cannot be an empty string.")
	}

	req, err := client.GetPropertiesPreparer(ctx, accountName, shareName, path, fileName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "files.Client", "GetProperties", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetPropertiesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "files.Client", "GetProperties", resp, "Failure sending request")
		return
	}

	result, err = client.GetPropertiesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "files.Client", "GetProperties", resp, "Failure responding to request")
		return
	}

	return
}

// GetPropertiesPreparer prepares the GetProperties request.
func (client Client) GetPropertiesPreparer(ctx context.Context, accountName, shareName, path, fileName string) (*http.Request, error) {
	if path != "" {
		path = fmt.Sprintf("%s/", path)
	}
	pathParameters := map[string]interface{}{
		"shareName": autorest.Encode("path", shareName),
		"directory": autorest.Encode("path", path),
		"fileName":  autorest.Encode("path", fileName),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsHead(),
		autorest.WithBaseURL(endpoints.GetFileEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{shareName}/{directory}{fileName}", pathParameters),
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
func (client Client) GetPropertiesResponder(resp *http.Response) (result GetResult, err error) {
	if resp != nil && resp.Header != nil {
		result.CacheControl = resp.Header.Get("Cache-Control")
		result.ContentDisposition = resp.Header.Get("Content-Disposition")
		result.ContentEncoding = resp.Header.Get("Content-Encoding")
		result.ContentLanguage = resp.Header.Get("Content-Language")
		result.ContentMD5 = resp.Header.Get("x-ms-content-md5")
		result.ContentType = resp.Header.Get("Content-Type")
		result.CopyID = resp.Header.Get("x-ms-copy-id")
		result.CopyProgress = resp.Header.Get("x-ms-copy-progress")
		result.CopySource = resp.Header.Get("x-ms-copy-source")
		result.CopyStatus = resp.Header.Get("x-ms-copy-status")
		result.CopyStatusDescription = resp.Header.Get("x-ms-copy-status-description")
		result.CopyCompletionTime = resp.Header.Get("x-ms-copy-completion-time")
		result.Encrypted = strings.EqualFold(resp.Header.Get("x-ms-server-encrypted"), "true")
		result.MetaData = metadata.ParseFromHeaders(resp.Header)

		contentLengthRaw := resp.Header.Get("Content-Length")
		if contentLengthRaw != "" {
			contentLength, err := strconv.Atoi(contentLengthRaw)
			if err != nil {
				return result, fmt.Errorf("Error parsing %q for Content-Length as an integer: %s", contentLengthRaw, err)
			}
			contentLengthI64 := int64(contentLength)
			result.ContentLength = &contentLengthI64
		}
	}

	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}

	return
}
