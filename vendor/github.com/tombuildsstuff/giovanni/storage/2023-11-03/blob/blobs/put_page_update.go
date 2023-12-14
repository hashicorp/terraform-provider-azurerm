package blobs

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

type PutPageUpdateInput struct {
	StartByte int64
	EndByte   int64
	Content   []byte

	IfSequenceNumberEQ *string
	IfSequenceNumberLE *string
	IfSequenceNumberLT *string
	IfModifiedSince    *string
	IfUnmodifiedSince  *string
	IfMatch            *string
	IfNoneMatch        *string
	LeaseID            *string
}

type PutPageUpdateResponse struct {
	HttpResponse *client.Response

	BlobSequenceNumber string
	ContentMD5         string
	LastModified       string
}

// PutPageUpdate writes a range of pages to a page blob.
func (c Client) PutPageUpdate(ctx context.Context, containerName, blobName string, input PutPageUpdateInput) (resp PutPageUpdateResponse, err error) {

	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return resp, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return resp, fmt.Errorf("`blobName` cannot be an empty string")
	}

	if input.StartByte < 0 {
		return resp, fmt.Errorf("`input.StartByte` must be greater than or equal to 0")
	}

	if input.EndByte <= 0 {
		return resp, fmt.Errorf("`input.EndByte` must be greater than 0")
	}

	expectedSize := (input.EndByte - input.StartByte) + 1
	actualSize := int64(len(input.Content))
	if expectedSize != actualSize {
		return resp, fmt.Errorf(fmt.Sprintf("Content Size was defined as %d but got %d.", expectedSize, actualSize))
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: putPageUpdateOptions{
			input: input,
		},
		Path: fmt.Sprintf("/%s/%s", containerName, blobName),
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

	if resp.HttpResponse != nil && resp.HttpResponse.Header != nil {
		resp.BlobSequenceNumber = resp.HttpResponse.Header.Get("x-ms-blob-sequence-number")
		resp.ContentMD5 = resp.HttpResponse.Header.Get("Content-MD5")
		resp.LastModified = resp.HttpResponse.Header.Get("Last-Modified")
	}

	return
}

type putPageUpdateOptions struct {
	input PutPageUpdateInput
}

func (p putPageUpdateOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	headers.Append("x-ms-page-write", "update")
	headers.Append("x-ms-range", fmt.Sprintf("bytes=%d-%d", p.input.StartByte, p.input.EndByte))
	headers.Append("Content-Length", strconv.Itoa(len(p.input.Content)))

	if p.input.LeaseID != nil {
		headers.Append("x-ms-lease-id", *p.input.LeaseID)
	}

	if p.input.IfSequenceNumberEQ != nil {
		headers.Append("x-ms-if-sequence-number-eq", *p.input.IfSequenceNumberEQ)
	}

	if p.input.IfSequenceNumberLE != nil {
		headers.Append("x-ms-if-sequence-number-le", *p.input.IfSequenceNumberLE)
	}

	if p.input.IfSequenceNumberLT != nil {
		headers.Append("x-ms-if-sequence-number-lt", *p.input.IfSequenceNumberLT)
	}

	if p.input.IfModifiedSince != nil {
		headers.Append("If-Modified-Since", *p.input.IfModifiedSince)
	}

	if p.input.IfUnmodifiedSince != nil {
		headers.Append("If-Unmodified-Since", *p.input.IfUnmodifiedSince)
	}

	if p.input.IfMatch != nil {
		headers.Append("If-Match", *p.input.IfMatch)
	}

	if p.input.IfNoneMatch != nil {
		headers.Append("If-None-Match", *p.input.IfNoneMatch)
	}

	return headers
}

func (p putPageUpdateOptions) ToOData() *odata.Query {
	return nil
}

func (p putPageUpdateOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "page")
	return out
}
