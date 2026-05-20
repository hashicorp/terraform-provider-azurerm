package paths

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type SetAccessControlInput struct {
	Owner *string
	Group *string
	ACL   *string

	// Optional - A date and time value.
	// Specify this header to perform the operation only if the resource has been modified since the specified date and time.
	IfModifiedSince *string

	// Optional - A date and time value.
	// Specify this header to perform the operation only if the resource has not been modified since the specified date and time.
	IfUnmodifiedSince *string
}

type SetPropertiesResponse struct {
	HttpResponse *http.Response
}

// SetAccessControl sets the access control properties for a Data Lake Store Gen2 Path within a Storage Account File System
func (c Client) SetAccessControl(ctx context.Context, fileSystemName string, path string, input SetAccessControlInput) (result SetPropertiesResponse, err error) {
	if fileSystemName == "" {
		err = fmt.Errorf("`fileSystemName` cannot be an empty string")
		return
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPatch,
		OptionsObject: setPropertyOptions{
			input: input,
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
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return result, err
	}

	return
}

type setPropertyOptions struct {
	input SetAccessControlInput
}

func (s setPropertyOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

	if s.input.ACL != nil {
		headers.Append("x-ms-acl", *s.input.ACL)
	}

	if s.input.Owner != nil {
		headers.Append("x-ms-owner", *s.input.Owner)
	}

	if s.input.Group != nil {
		headers.Append("x-ms-group", *s.input.Group)
	}

	if s.input.IfModifiedSince != nil {
		headers.Append("If-Modified-Since", *s.input.IfModifiedSince)
	}

	if s.input.IfUnmodifiedSince != nil {
		headers.Append("If-Unmodified-Since", *s.input.IfUnmodifiedSince)
	}

	return headers
}

func (s setPropertyOptions) ToOData() *odata.Query {
	return nil
}

func (s setPropertyOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("action", "setAccessControl")
	return out
}
