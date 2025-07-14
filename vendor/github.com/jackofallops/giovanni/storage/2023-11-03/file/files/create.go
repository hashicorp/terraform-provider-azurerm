package files

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/jackofallops/giovanni/storage/internal/metadata"
)

type CreateInput struct {
	// This header specifies the maximum size for the file, up to 1 TiB.
	ContentLength int64

	// The MIME content type of the file
	// If not specified, the default type is application/octet-stream.
	ContentType *string

	// Specifies which content encodings have been applied to the file.
	// This value is returned to the client when the Get File operation is performed
	// on the file resource and can be used to decode file content.
	ContentEncoding *string

	// Specifies the natural languages used by this resource.
	ContentLanguage *string

	// The File service stores this value but does not use or modify it.
	CacheControl *string

	// Sets the file's MD5 hash.
	ContentMD5 *string

	// Sets the file’s Content-Disposition header.
	ContentDisposition *string

	// The time at which this file was created at - if omitted, this'll be set to "now"
	// This maps to the `x-ms-file-creation-time` field.
	CreatedAt *time.Time

	// The time at which this file was last modified - if omitted, this'll be set to "now"
	// This maps to the `x-ms-file-last-write-time` field.
	LastModified *time.Time

	// MetaData is a mapping of key value pairs which should be assigned to this file
	MetaData map[string]string
}

type CreateResponse struct {
	HttpResponse *http.Response
}

// Create creates a new file or replaces a file.
func (c Client) Create(ctx context.Context, shareName, path, fileName string, input CreateInput) (result CreateResponse, err error) {
	if shareName == "" {
		err = fmt.Errorf("`shareName` cannot be an empty string")
		return
	}

	if strings.ToLower(shareName) != shareName {
		err = fmt.Errorf("`shareName` must be a lower-cased string")
		return
	}

	if fileName == "" {
		err = fmt.Errorf("`fileName` cannot be an empty string")
		return
	}

	if err = metadata.Validate(input.MetaData); err != nil {
		err = fmt.Errorf("`input.MetaData` is not valid: %s", err)
		return
	}

	if path != "" {
		path = fmt.Sprintf("%s/", path)
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
		Path: fmt.Sprintf("/%s/%s%s", shareName, path, fileName),
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
	input CreateInput
}

func (c CreateOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

	var coalesceDate = func(input *time.Time, defaultVal string) string {
		if input == nil {
			return defaultVal
		}

		return input.Format(time.RFC1123)
	}

	if len(c.input.MetaData) > 0 {
		headers.Merge(metadata.SetMetaDataHeaders(c.input.MetaData))
	}

	headers.Append("x-ms-content-length", strconv.Itoa(int(c.input.ContentLength)))
	headers.Append("x-ms-type", "file")

	headers.Append("x-ms-file-permission", "inherit") // TODO: expose this in future
	headers.Append("x-ms-file-attributes", "None")    // TODO: expose this in future
	headers.Append("x-ms-file-creation-time", coalesceDate(c.input.CreatedAt, "now"))
	headers.Append("x-ms-file-last-write-time", coalesceDate(c.input.LastModified, "now"))

	if c.input.ContentDisposition != nil {
		headers.Append("x-ms-content-disposition", *c.input.ContentDisposition)
	}

	if c.input.ContentEncoding != nil {
		headers.Append("x-ms-content-encoding", *c.input.ContentEncoding)
	}

	if c.input.ContentMD5 != nil {
		headers.Append("x-ms-content-md5", *c.input.ContentMD5)
	}

	if c.input.ContentType != nil {
		headers.Append("x-ms-content-type", *c.input.ContentType)
	}

	return headers
}

func (c CreateOptions) ToOData() *odata.Query {
	return nil
}

func (c CreateOptions) ToQuery() *client.QueryParams {
	return nil
}
