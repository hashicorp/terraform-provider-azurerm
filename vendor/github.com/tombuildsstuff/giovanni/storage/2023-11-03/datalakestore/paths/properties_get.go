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
	HttpResponse *http.Response

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
	Action GetPropertiesAction
}

type GetPropertiesAction string

const (
	GetPropertiesActionGetStatus        GetPropertiesAction = "getStatus"
	GetPropertiesActionGetAccessControl GetPropertiesAction = "getAccessControl"
)

// GetProperties gets the properties for a Data Lake Store Gen2 Path in a FileSystem within a Storage Account
func (c Client) GetProperties(ctx context.Context, fileSystemName string, path string, input GetPropertiesInput) (result GetPropertiesResponse, err error) {
	if fileSystemName == "" {
		err = fmt.Errorf("`fileSystemName` cannot be an empty string")
		return
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodHead,
		OptionsObject: getPropertyOptions{
			action: input.Action,
		},
		Path: fmt.Sprintf("/%s/%s", fileSystemName, path),
	}

	req, err := c.Client.NewRequest(ctx, opts)

	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return result, err
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil && resp.Response != nil {
		result.HttpResponse = resp.Response

		if err == nil {
			if resp.Header != nil {
				result.ResourceType = PathResource(resp.Header.Get("x-ms-resource-type"))
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
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return result, err
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
