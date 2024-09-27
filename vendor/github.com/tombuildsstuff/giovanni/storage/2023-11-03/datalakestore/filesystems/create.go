package filesystems

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type CreateInput struct {
	// The encryption scope to set as the default on the filesystem.
	DefaultEncryptionScope string

	// A map of base64-encoded strings to store as user-defined properties with the File System
	// Note that items may only contain ASCII characters in the ISO-8859-1 character set.
	// This automatically gets converted to a comma-separated list of name and
	// value pairs before sending to the API
	Properties map[string]string
}

type CreateResponse struct {
	HttpResponse *http.Response
}

// Create creates a Data Lake Store Gen2 FileSystem within a Storage Account
func (c Client) Create(ctx context.Context, fileSystemName string, input CreateInput) (result CreateResponse, err error) {
	if fileSystemName == "" {
		err = fmt.Errorf("`fileSystemName` cannot be an empty string")
		return
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: createOptions{
			input: input,
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

type createOptions struct {
	input CreateInput
}

func (o createOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

	if o.input.DefaultEncryptionScope != "" {
		headers.Append("x-ms-default-encryption-scope", o.input.DefaultEncryptionScope)
	}

	props := buildProperties(o.input.Properties)
	if props != "" {
		headers.Append("x-ms-properties", props)
	}

	return headers
}

func (createOptions) ToOData() *odata.Query {
	return nil
}

func (createOptions) ToQuery() *client.QueryParams {
	return fileSystemOptions{}.ToQuery()
}
