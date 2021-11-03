package paths

import (
	"context"
	"net/http"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

type GetPropertiesResponse struct {
	autorest.Response

	ETag         string
	LastModified time.Time
	// ResourceType is only returned for GetPropertiesActionGetStatus requests
	ResourceType PathResource
	Owner        string
	Group        string
	// ACL is only returned for GetPropertiesActionGetAccessControl requests
	ACL string
}

type GetPropertiesAction string

const (
	GetPropertiesActionGetStatus        GetPropertiesAction = "getStatus"
	GetPropertiesActionGetAccessControl GetPropertiesAction = "getAccessControl"
)

// GetProperties gets the properties for a Data Lake Store Gen2 Path in a FileSystem within a Storage Account
func (client Client) GetProperties(ctx context.Context, accountName string, fileSystemName string, path string, action GetPropertiesAction) (result GetPropertiesResponse, err error) {
	if accountName == "" {
		return result, validation.NewError("datalakestore.Client", "GetProperties", "`accountName` cannot be an empty string.")
	}
	if fileSystemName == "" {
		return result, validation.NewError("datalakestore.Client", "GetProperties", "`fileSystemName` cannot be an empty string.")
	}

	req, err := client.GetPropertiesPreparer(ctx, accountName, fileSystemName, path, action)
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
func (client Client) GetPropertiesPreparer(ctx context.Context, accountName string, fileSystemName string, path string, action GetPropertiesAction) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"fileSystemName": autorest.Encode("path", fileSystemName),
		"path":           autorest.Encode("path", path),
	}

	queryParameters := map[string]interface{}{
		"action": autorest.Encode("query", string(action)),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsHead(),
		autorest.WithBaseURL(endpoints.GetDataLakeStoreEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{fileSystemName}/{path}", pathParameters),
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
	result = GetPropertiesResponse{}
	if resp != nil && resp.Header != nil {

		resourceTypeRaw := resp.Header.Get("x-ms-resource-type")
		var resourceType PathResource
		if resourceTypeRaw != "" {
			resourceType, err = parsePathResource(resourceTypeRaw)
			if err != nil {
				return GetPropertiesResponse{}, err
			}
			result.ResourceType = resourceType
		}
		result.ETag = resp.Header.Get("ETag")

		if lastModifiedRaw := resp.Header.Get("Last-Modified"); lastModifiedRaw != "" {
			lastModified, err := time.Parse(time.RFC1123, lastModifiedRaw)
			if err != nil {
				return GetPropertiesResponse{}, err
			}
			result.LastModified = lastModified
		}

		result.Owner = resp.Header.Get("x-ms-owner")
		result.Group = resp.Header.Get("x-ms-group")
		result.ACL = resp.Header.Get("x-ms-acl")
	}
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),

		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}

	return result, err
}
