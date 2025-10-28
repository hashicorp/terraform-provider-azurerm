package blobs

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type GetInput struct {
	LeaseID   *string
	StartByte *int64
	EndByte   *int64
}

type GetResponse struct {
	HttpResponse *http.Response

	Contents *[]byte
}

// Get reads or downloads a blob from the system, including its metadata and properties.
func (c Client) Get(ctx context.Context, containerName, blobName string, input GetInput) (result GetResponse, err error) {
	if containerName == "" {
		return result, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return result, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return result, fmt.Errorf("`blobName` cannot be an empty string")
	}

	if input.LeaseID != nil && *input.LeaseID == "" {
		return result, fmt.Errorf("`input.LeaseID` should either be specified or nil, not an empty string")
	}

	if (input.StartByte != nil && input.EndByte == nil) || input.StartByte == nil && input.EndByte != nil {
		return result, fmt.Errorf("`input.StartByte` and `input.EndByte` must both be specified, or both be nil")
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusOK,
			http.StatusPartialContent,
		},
		HttpMethod: http.MethodGet,
		OptionsObject: getOptions{
			input: input,
		},
		Path: fmt.Sprintf("/%s/%s", containerName, blobName),
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
			if resp.Body != nil {
				defer resp.Body.Close()
				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					return result, fmt.Errorf("could not parse response body")
				}

				result.Contents = &respBody
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

type getOptions struct {
	input GetInput
}

func (g getOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	if g.input.StartByte != nil && g.input.EndByte != nil {
		headers.Append("x-ms-range", fmt.Sprintf("bytes=%d-%d", *g.input.StartByte, *g.input.EndByte))
	}
	return headers

}

func (g getOptions) ToOData() *odata.Query {
	return nil
}

func (g getOptions) ToQuery() *client.QueryParams {
	return nil
}
