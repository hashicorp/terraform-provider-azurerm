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
	HttpResponse *client.Response
}

// PutByteRange puts the specified Byte Range in the specified File.
func (c Client) PutByteRange(ctx context.Context, shareName, path, fileName string, input PutByteRangeInput) (resp PutRangeResponse, err error) {

	if shareName == "" {
		return resp, fmt.Errorf("`shareName` cannot be an empty string")
	}

	if strings.ToLower(shareName) != shareName {
		return resp, fmt.Errorf("`shareName` must be a lower-cased string")
	}

	if fileName == "" {
		return resp, fmt.Errorf("`fileName` cannot be an empty string")
	}

	if input.StartBytes < 0 {
		return resp, fmt.Errorf("`input.StartBytes` must be greater or equal to 0")
	}
	if input.EndBytes <= 0 {
		return resp, fmt.Errorf("`input.EndBytes` must be greater than 0")
	}

	expectedBytes := input.EndBytes - input.StartBytes
	actualBytes := len(input.Content)
	if expectedBytes != int64(actualBytes) {
		return resp, fmt.Errorf(fmt.Sprintf("The specified byte-range (%d) didn't match the content size (%d).", expectedBytes, actualBytes))
	}

	if expectedBytes > (4 * 1024 * 1024) {
		return resp, fmt.Errorf("specified Byte Range must be at most 4MB")
	}

	if path != "" {
		path = fmt.Sprintf("/%s/", path)
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

	resp.HttpResponse, err = req.Execute(ctx)
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
