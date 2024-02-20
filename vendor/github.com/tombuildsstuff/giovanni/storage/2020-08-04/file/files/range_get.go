package files

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type GetByteRangeInput struct {
	StartBytes int64
	EndBytes   int64
}

type GetByteRangeResponse struct {
	HttpResponse *client.Response

	Contents []byte
}

// GetByteRange returns the specified Byte Range from the specified File.
func (c Client) GetByteRange(ctx context.Context, shareName, path, fileName string, input GetByteRangeInput) (resp GetByteRangeResponse, err error) {

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
	if expectedBytes < (4 * 1024) {
		return resp, fmt.Errorf("requested Byte Range must be at least 4KB")
	}
	if expectedBytes > (4 * 1024 * 1024) {
		return resp, fmt.Errorf("requested Byte Range must be at most 4MB")
	}

	if path != "" {
		path = fmt.Sprintf("%s/", path)
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusOK,
			http.StatusPartialContent,
		},
		HttpMethod: http.MethodGet,
		OptionsObject: GetByteRangeOptions{
			input: input,
		},
		Path: fmt.Sprintf("/%s/%s%s", shareName, path, fileName),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	resp.HttpResponse, err = req.Execute(ctx)
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	if resp.HttpResponse != nil {
		bytes, err := io.ReadAll(resp.HttpResponse.Body)
		if err != nil {
			return resp, fmt.Errorf("reading response body: %v", err)
		}
		resp.Contents = bytes
	}

	return
}

type GetByteRangeOptions struct {
	input GetByteRangeInput
}

func (g GetByteRangeOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	headers.Append("x-ms-range", fmt.Sprintf("bytes=%d-%d", g.input.StartBytes, g.input.EndBytes-1))
	return headers
}

func (g GetByteRangeOptions) ToOData() *odata.Query {
	return nil
}

func (g GetByteRangeOptions) ToQuery() *client.QueryParams {
	return nil
}
