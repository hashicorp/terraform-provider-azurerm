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
	"github.com/jackofallops/giovanni/storage/internal/metadata"
)

type PutBlockBlobInput struct {
	CacheControl       *string
	Content            *[]byte
	ContentDisposition *string
	ContentEncoding    *string
	ContentLanguage    *string
	ContentMD5         *string
	ContentType        *string
	LeaseID            *string
	EncryptionScope    *string
	MetaData           map[string]string
}

type PutBlockBlobResponse struct {
	HttpResponse *http.Response
}

// PutBlockBlob is a wrapper around the Put API call (with a stricter input object)
// which creates a new block append blob, or updates the content of an existing block blob.
func (c Client) PutBlockBlob(ctx context.Context, containerName, blobName string, input PutBlockBlobInput) (result PutBlockBlobResponse, err error) {
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

	if input.Content != nil && len(*input.Content) == 0 {
		err = fmt.Errorf("`input.Content` must either be nil or not empty")
		return
	}

	if err = metadata.Validate(input.MetaData); err != nil {
		err = fmt.Errorf(fmt.Sprintf("`input.MetaData` is not valid: %s.", err))
		return
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: putBlockBlobOptions{
			input: input,
		},
		Path: fmt.Sprintf("/%s/%s", containerName, blobName),
	}

	if input.ContentType != nil {
		opts.ContentType = *input.ContentType
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	if input.Content != nil {
		req.ContentLength = int64(len(*input.Content))
		req.Body = io.NopCloser(bytes.NewReader(*input.Content))
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

type putBlockBlobOptions struct {
	input PutBlockBlobInput
}

func (p putBlockBlobOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	headers.Append("x-ms-blob-type", string(BlockBlob))

	if p.input.CacheControl != nil {
		headers.Append("x-ms-blob-cache-control", *p.input.CacheControl)
	}
	if p.input.ContentDisposition != nil {
		headers.Append("x-ms-blob-content-disposition", *p.input.ContentDisposition)
	}
	if p.input.ContentEncoding != nil {
		headers.Append("x-ms-blob-content-encoding", *p.input.ContentEncoding)
	}
	if p.input.ContentLanguage != nil {
		headers.Append("x-ms-blob-content-language", *p.input.ContentLanguage)
	}
	if p.input.ContentMD5 != nil {
		headers.Append("x-ms-blob-content-md5", *p.input.ContentMD5)
	}
	if p.input.ContentType != nil {
		headers.Append("x-ms-blob-content-type", *p.input.ContentType)
	}
	if p.input.LeaseID != nil {
		headers.Append("x-ms-lease-id", *p.input.LeaseID)
	}
	if p.input.EncryptionScope != nil {
		headers.Append("x-ms-encryption-scope", *p.input.EncryptionScope)
	}
	if p.input.Content != nil {
		headers.Append("Content-Length", strconv.Itoa(len(*p.input.Content)))
	}

	headers.Merge(metadata.SetMetaDataHeaders(p.input.MetaData))

	return headers
}

func (p putBlockBlobOptions) ToOData() *odata.Query {
	return nil
}

func (p putBlockBlobOptions) ToQuery() *client.QueryParams {
	return nil
}
