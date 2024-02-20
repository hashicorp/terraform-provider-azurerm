package blobs

import (
	"context"
	"fmt"
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
	HttpResponse *client.Response

	Contents []byte
}

// Get reads or downloads a blob from the system, including its metadata and properties.
func (c Client) Get(ctx context.Context, containerName, blobName string, input GetInput) (resp GetResponse, err error) {

	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return resp, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return resp, fmt.Errorf("`blobName` cannot be an empty string")
	}

	if input.LeaseID != nil && *input.LeaseID == "" {
		return resp, fmt.Errorf("`input.LeaseID` should either be specified or nil, not an empty string")
	}

	if (input.StartByte != nil && input.EndByte == nil) || input.StartByte == nil && input.EndByte != nil {
		return resp, fmt.Errorf("`input.StartByte` and `input.EndByte` must both be specified, or both be nil")
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

	resp.HttpResponse, err = req.Execute(ctx)
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	if resp.HttpResponse != nil {
		err = resp.HttpResponse.Unmarshal(&resp.Contents)
		if err != nil {
			return resp, fmt.Errorf("unmarshalling response: %v", err)
		}
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
