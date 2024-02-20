package paths

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type GetPropertiesResponse struct {
	HttpResponse *client.Response

	ETag         string
	LastModified time.Time
	// ResourceType is only returned for GetPropertiesActionGetStatus requests
	ResourceType PathResource
	Owner        string
	Group        string
	// ACL is only returned for GetPropertiesActionGetAccessControl requests
	ACL string
}

type GetPropertiesInput struct {
	action GetPropertiesAction
}

type GetPropertiesAction string

const (
	GetPropertiesActionGetStatus        GetPropertiesAction = "getStatus"
	GetPropertiesActionGetAccessControl GetPropertiesAction = "getAccessControl"
)

// GetProperties gets the properties for a Data Lake Store Gen2 Path in a FileSystem within a Storage Account
func (c Client) GetProperties(ctx context.Context, fileSystemName string, path string, input GetPropertiesInput) (resp GetPropertiesResponse, err error) {
	if fileSystemName == "" {
		return resp, fmt.Errorf("`fileSystemName` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodHead,
		OptionsObject: getPropertyOptions{
			action: input.action,
		},
		Path: fmt.Sprintf("/%s/%s", fileSystemName, path),
	}

	req, err := c.Client.NewRequest(ctx, opts)

	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return resp, err
	}

	resp.HttpResponse, err = req.Execute(ctx)
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return resp, err
	}

	if resp.HttpResponse != nil {
		resp.ResourceType = PathResource(resp.HttpResponse.Header.Get("x-ms-resource-type"))
		resp.ETag = resp.HttpResponse.Header.Get("ETag")

		if lastModifiedRaw := resp.HttpResponse.Header.Get("Last-Modified"); lastModifiedRaw != "" {
			lastModified, err := time.Parse(time.RFC1123, lastModifiedRaw)
			if err != nil {
				return GetPropertiesResponse{}, err
			}
			resp.LastModified = lastModified
		}

		resp.Owner = resp.HttpResponse.Header.Get("x-ms-owner")
		resp.Group = resp.HttpResponse.Header.Get("x-ms-group")
		resp.ACL = resp.HttpResponse.Header.Get("x-ms-acl")
	}
	return
}

type getPropertyOptions struct {
	action GetPropertiesAction
}

func (g getPropertyOptions) ToHeaders() *client.Headers {
	return nil
}

func (g getPropertyOptions) ToOData() *odata.Query {
	return nil
}

func (g getPropertyOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("action", string(g.action))
	return out
}
