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

type SetPropertiesInput struct {
	// Resizes a file to the specified size.
	// If the specified byte value is less than the current size of the file,
	// then all ranges above the specified byte value are cleared.
	ContentLength int64

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

	// MetaData is a mapping of key value pairs which should be assigned to this file
	MetaData map[string]string
}

type SetPropertiesResponse struct {
	HttpResponse *http.Response
}

// SetProperties sets the specified properties on the specified File
func (c Client) SetProperties(ctx context.Context, shareName, path, fileName string, input SetPropertiesInput) (result SetPropertiesResponse, err error) {
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

	if path != "" {
		path = fmt.Sprintf("%s/", path)
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: SetPropertiesOptions{
			input: input,
		},
		Path: fmt.Sprintf("%s/%s%s", shareName, path, fileName),
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

type SetPropertiesOptions struct {
	input SetPropertiesInput
}

func (s SetPropertiesOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	var coalesceDate = func(input *time.Time, defaultVal string) string {
		if input == nil {
			return defaultVal
		}

		return input.Format(time.RFC1123)
	}

	headers.Append("x-ms-type", "file")

	headers.Append("x-ms-content-length", strconv.Itoa(int(s.input.ContentLength)))
	headers.Append("x-ms-file-permission", "inherit") // TODO: expose this in future
	headers.Append("x-ms-file-attributes", "None")    // TODO: expose this in future
	headers.Append("x-ms-file-creation-time", coalesceDate(s.input.CreatedAt, "now"))
	headers.Append("x-ms-file-last-write-time", coalesceDate(s.input.LastModified, "now"))

	if s.input.ContentControl != nil {
		headers.Append("x-ms-cache-control", *s.input.ContentControl)
	}
	if s.input.ContentDisposition != nil {
		headers.Append("x-ms-content-disposition", *s.input.ContentDisposition)
	}
	if s.input.ContentEncoding != nil {
		headers.Append("x-ms-content-encoding", *s.input.ContentEncoding)
	}
	if s.input.ContentLanguage != nil {
		headers.Append("x-ms-content-language", *s.input.ContentLanguage)
	}
	if s.input.ContentMD5 != nil {
		headers.Append("x-ms-content-md5", *s.input.ContentMD5)
	}
	if s.input.ContentType != nil {
		headers.Append("x-ms-content-type", *s.input.ContentType)
	}

	if len(s.input.MetaData) > 0 {
		headers.Merge(metadata.SetMetaDataHeaders(s.input.MetaData))
	}

	return headers
}

func (s SetPropertiesOptions) ToOData() *odata.Query {
	return nil
}

func (s SetPropertiesOptions) ToQuery() *client.QueryParams {
	return nil
}
