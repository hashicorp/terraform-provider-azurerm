package directories

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/jackofallops/giovanni/storage/internal/metadata"
)

type CreateDirectoryInput struct {
	// The time at which this file was created at - if omitted, this'll be set to "now"
	// This maps to the `x-ms-file-creation-time` field.
	// ... Yes I know it says File not Directory, I didn't design the API.
	CreatedAt *time.Time

	// The time at which this file was last modified - if omitted, this'll be set to "now"
	// This maps to the `x-ms-file-last-write-time` field.
	// ... Yes I know it says File not Directory, I didn't design the API.
	LastModified *time.Time

	// MetaData is a mapping of key value pairs which should be assigned to this directory
	MetaData map[string]string
}

type CreateDirectoryResponse struct {
	HttpResponse *http.Response
}

// Create creates a new directory under the specified share or parent directory.
func (c Client) Create(ctx context.Context, shareName, path string, input CreateDirectoryInput) (result CreateDirectoryResponse, err error) {

	if shareName == "" {
		err = fmt.Errorf("`shareName` cannot be an empty string")
		return
	}

	if strings.ToLower(shareName) != shareName {
		err = fmt.Errorf("`shareName` must be a lower-cased string")
		return
	}

	if err = metadata.Validate(input.MetaData); err != nil {
		err = fmt.Errorf("`input.MetaData` is not valid: %s", err)
		return
	}

	if path == "" {
		err = fmt.Errorf("`path` cannot be an empty string")
		return
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: CreateOptions{
			input: input,
		},
		Path: fmt.Sprintf("/%s/%s", shareName, path),
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

type CreateOptions struct {
	input CreateDirectoryInput
}

func (c CreateOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

	if len(c.input.MetaData) > 0 {
		headers.Merge(metadata.SetMetaDataHeaders(c.input.MetaData))
	}

	var coalesceDate = func(input *time.Time, defaultVal string) string {
		if input == nil {
			return defaultVal
		}
		return input.Format(time.RFC1123)
	}

	// ... Yes I know these say File not Directory, I didn't design the API.
	headers.Append("x-ms-file-permission", "inherit") // TODO: expose this in future
	headers.Append("x-ms-file-attributes", "None")    // TODO: expose this in future
	headers.Append("x-ms-file-creation-time", coalesceDate(c.input.CreatedAt, "now"))
	headers.Append("x-ms-file-last-write-time", coalesceDate(c.input.LastModified, "now"))

	return headers
}

func (c CreateOptions) ToOData() *odata.Query {
	return nil
}

func (c CreateOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("restype", "directory")
	return out
}
