package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type PutBlockInput struct {
	BlockID    string
	Content    []byte
	ContentMD5 *string
	LeaseID    *string
}

type PutBlockResponse struct {
	HttpResponse *client.Response

	ContentMD5 string
}

// PutBlock creates a new block to be committed as part of a blob.
func (c Client) PutBlock(ctx context.Context, containerName, blobName string, input PutBlockInput) (resp PutBlockResponse, err error) {

	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return resp, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return resp, fmt.Errorf("`blobName` cannot be an empty string")
	}

	if input.BlockID == "" {
		return resp, fmt.Errorf("`input.BlockID` cannot be an empty string")
	}

	if len(input.Content) == 0 {
		return resp, fmt.Errorf("`input.Content` cannot be empty")
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: putBlockOptions{
			input: input,
		},
		Path: fmt.Sprintf("/%s/%s", containerName, blobName),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	err = req.Marshal(&input.Content)
	if err != nil {
		return resp, fmt.Errorf("marshalling request: %v", err)
	}

	resp.HttpResponse, err = req.Execute(ctx)
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

type putBlockOptions struct {
	input PutBlockInput
}

func (p putBlockOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	headers.Append("Content-Length", strconv.Itoa(len(p.input.Content)))

	if p.input.ContentMD5 != nil {
		headers.Append("x-ms-blob-content-md5", *p.input.ContentMD5)
	}
	if p.input.LeaseID != nil {
		headers.Append("x-ms-lease-id", *p.input.LeaseID)
	}

	return headers
}

func (p putBlockOptions) ToOData() *odata.Query {
	return nil
}

func (p putBlockOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "block")
	out.Append("blockid", p.input.BlockID)
	return out
}
