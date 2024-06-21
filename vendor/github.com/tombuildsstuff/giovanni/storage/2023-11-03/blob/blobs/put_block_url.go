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

	ContentMD5      *string
	LeaseID         *string
	Range           *string
	EncryptionScope *string
}

type PutBlockFromURLResponse struct {
	ContentMD5   string
	HttpResponse *http.Response
}

// PutBlockFromURL creates a new block to be committed as part of a blob where the contents are read from a URL
func (c Client) PutBlockFromURL(ctx context.Context, containerName, blobName string, input PutBlockFromURLInput) (result PutBlockFromURLResponse, err error) {
	if containerName == "" {
		err = fmt.Errorf("`containerName` cannot be an empty string")
		return
	}

	if strings.ToLower(containerName) != containerName {
		err = fmt.Errorf("`containerName` must be a lower-cased string")
		return
	}

	if blobName == "" {
		err = fmt.Errorf("`blobName` cannot be an empty string")
		return
	}

	if input.BlockID == "" {
		err = fmt.Errorf("`input.BlockID` cannot be an empty string")
		return
	}

	if input.CopySource == "" {
		err = fmt.Errorf("`input.CopySource` cannot be an empty string")
		return
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

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil && resp.Response != nil {
		result.HttpResponse = resp.Response

		if err == nil {
			if resp.Header != nil {
				result.ContentMD5 = resp.Header.Get("Content-MD5")
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
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
	if p.input.EncryptionScope != nil {
		headers.Append("x-ms-encryption-scope", *p.input.EncryptionScope)
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
