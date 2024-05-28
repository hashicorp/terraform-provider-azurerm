package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type PutPageClearInput struct {
	StartByte int64
	EndByte   int64

	LeaseID         *string
	EncryptionScope *string
}

type PutPageClearResponse struct {
	HttpResponse *http.Response
}

// PutPageClear clears a range of pages within a page blob.
func (c Client) PutPageClear(ctx context.Context, containerName, blobName string, input PutPageClearInput) (result PutPageClearResponse, err error) {
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

	if input.StartByte < 0 {
		err = fmt.Errorf("`input.StartByte` must be greater than or equal to 0")
		return
	}

	if input.EndByte <= 0 {
		err = fmt.Errorf("`input.EndByte` must be greater than 0")
		return
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: putPageClearOptions{
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
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

type putPageClearOptions struct {
	input PutPageClearInput
}

func (p putPageClearOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

	headers.Append("x-ms-page-write", "clear")
	headers.Append("x-ms-range", fmt.Sprintf("bytes=%d-%d", p.input.StartByte, p.input.EndByte))

	if p.input.LeaseID != nil {
		headers.Append("x-ms-lease-id", *p.input.LeaseID)
	}
	if p.input.EncryptionScope != nil {
		headers.Append("x-ms-encryption-scope", *p.input.EncryptionScope)
	}

	return headers
}

func (p putPageClearOptions) ToOData() *odata.Query {
	return nil
}

func (p putPageClearOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "page")
	return out
}
