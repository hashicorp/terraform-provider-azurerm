package files

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
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
	HttpResponse *http.Response
	Contents     *[]byte
}

// GetByteRange returns the specified Byte Range from the specified File.
func (c Client) GetByteRange(ctx context.Context, shareName, path, fileName string, input GetByteRangeInput) (result GetByteRangeResponse, err error) {

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
	if expectedBytes < (4 * 1024) {
		err = fmt.Errorf("requested Byte Range must be at least 4KB")
		return
	}
	if expectedBytes > (4 * 1024 * 1024) {
		err = fmt.Errorf("requested Byte Range must be at most 4MB")
		return
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

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil && resp.Response != nil {
		result.HttpResponse = resp.Response

		if err == nil {
			result.Contents = &[]byte{}
			if resp.Body != nil {
				respBody, err := io.ReadAll(resp.Body)
				defer resp.Body.Close()
				if err != nil {
					return result, fmt.Errorf("could not parse response body")
				}

				if respBody != nil {
					result.Contents = pointer.To(respBody)
				}
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
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
