package files

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type PutByteRangeInput struct {
	StartBytes int64
	EndBytes   int64

	// Content is the File Contents for the specified range
	// which can be at most 4MB
	Content []byte
}

type PutRangeResponse struct {
	HttpResponse *http.Response
}

// PutByteRange puts the specified Byte Range in the specified File.
func (c Client) PutByteRange(ctx context.Context, shareName, path, fileName string, input PutByteRangeInput) (result PutRangeResponse, err error) {
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

	if input.StartBytes < 0 {
		err = fmt.Errorf("`input.StartBytes` must be greater or equal to 0")
		return
	}
	if input.EndBytes <= 0 {
		err = fmt.Errorf("`input.EndBytes` must be greater than 0")
		return
	}

	expectedBytes := input.EndBytes - input.StartBytes
	actualBytes := len(input.Content)
	if expectedBytes != int64(actualBytes) {
		err = fmt.Errorf(fmt.Sprintf("The specified byte-range (%d) didn't match the content size (%d).", expectedBytes, actualBytes))
		return
	}

	if expectedBytes > (4 * 1024 * 1024) {
		err = fmt.Errorf("specified Byte Range must be at most 4MB")
		return
	}

	if path != "" {
		path = fmt.Sprintf("%s/", path)
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: PutRangeOptions{
			input: input,
		},
		Path: fmt.Sprintf("%s/%s%s", shareName, path, fileName),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	req.Body = io.NopCloser(bytes.NewReader(input.Content))
	req.ContentLength = int64(len(input.Content))

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

type PutRangeOptions struct {
	input PutByteRangeInput
}

func (p PutRangeOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	headers.Append("x-ms-write", "update")
	headers.Append("x-ms-range", fmt.Sprintf("bytes=%d-%d", p.input.StartBytes, p.input.EndBytes-1))
	headers.Append("Content-Length", strconv.Itoa(len(p.input.Content)))
	return headers
}

func (p PutRangeOptions) ToOData() *odata.Query {
	return nil
}

func (p PutRangeOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "range")
	return out
}
