package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type PutBlockFromURLInput struct {
	BlockID    string
	CopySource string

	ContentMD5 *string
	LeaseID    *string
	Range      *string
}

type PutBlockFromURLResponse struct {
	HttpResponse *client.Response
	ContentMD5   string
}

// PutBlockFromURL creates a new block to be committed as part of a blob where the contents are read from a URL
func (c Client) PutBlockFromURL(ctx context.Context, containerName, blobName string, input PutBlockFromURLInput) (resp PutBlockFromURLResponse, err error) {

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

	if input.CopySource == "" {
		return resp, fmt.Errorf("`input.CopySource` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: putBlockUrlOptions{
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
		if resp.HttpResponse.Header != nil {
			resp.ContentMD5 = resp.HttpResponse.Header.Get("Content-MD5")
		}
	}

	return
}

type putBlockUrlOptions struct {
	input PutBlockFromURLInput
}

func (p putBlockUrlOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

	headers.Append("x-ms-copy-source", p.input.CopySource)

	if p.input.ContentMD5 != nil {
		headers.Append("x-ms-source-content-md5", *p.input.ContentMD5)
	}
	if p.input.LeaseID != nil {
		headers.Append("x-ms-lease-id", *p.input.LeaseID)
	}
	if p.input.Range != nil {
		headers.Append("x-ms-source-range", *p.input.Range)
	}
	return headers
}

func (p putBlockUrlOptions) ToOData() *odata.Query {
	return nil
}

func (p putBlockUrlOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "block")
	out.Append("blockid", p.input.BlockID)
	return out
}
