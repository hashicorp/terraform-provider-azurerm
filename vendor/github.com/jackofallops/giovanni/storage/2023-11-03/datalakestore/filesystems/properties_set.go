package filesystems

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
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

type SetPropertiesResponse struct {
	HttpResponse *http.Response
}

// SetProperties sets the Properties for a Data Lake Store Gen2 FileSystem within a Storage Account
func (c Client) SetProperties(ctx context.Context, fileSystemName string, input SetPropertiesInput) (result SetPropertiesResponse, err error) {
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
		OptionsObject: setPropertiesOptions{
			properties:        input.Properties,
			ifUnmodifiedSince: input.IfUnmodifiedSince,
			ifModifiedSince:   input.IfModifiedSince,
		},

		Path: fmt.Sprintf("/%s", fileSystemName),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil && resp.Response != nil {
		result.HttpResponse = resp.Response
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

type setPropertiesOptions struct {
	properties        map[string]string
	ifModifiedSince   *string
	ifUnmodifiedSince *string
}

func (o setPropertiesOptions) ToHeaders() *client.Headers {

	headers := &client.Headers{}
	props := buildProperties(o.properties)
	if props != "" {
		headers.Append("x-ms-properties", props)
	}

	if o.ifModifiedSince != nil {
		headers.Append("If-Modified-Since", *o.ifModifiedSince)
	}
	if o.ifUnmodifiedSince != nil {
		headers.Append("If-Unmodified-Since", *o.ifUnmodifiedSince)
	}

	return headers
}

func (setPropertiesOptions) ToOData() *odata.Query {
	return nil
}

func (setPropertiesOptions) ToQuery() *client.QueryParams {
	return fileSystemOptions{}.ToQuery()
}
