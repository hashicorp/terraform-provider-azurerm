package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/jackofallops/giovanni/storage/internal/metadata"
)

type PutPageBlobInput struct {
	CacheControl       *string
	ContentDisposition *string
	ContentEncoding    *string
	ContentLanguage    *string
	ContentMD5         *string
	ContentType        *string
	LeaseID            *string
	EncryptionScope    *string
	MetaData           map[string]string

	BlobContentLengthBytes int64
	BlobSequenceNumber     *int64
	AccessTier             *AccessTier
}

type PutPageBlobResponse struct {
	HttpResponse *http.Response
}

// PutPageBlob is a wrapper around the Put API call (with a stricter input object)
// which creates a new block blob, or updates the content of an existing page blob.
func (c Client) PutPageBlob(ctx context.Context, containerName, blobName string, input PutPageBlobInput) (result PutPageBlobResponse, err error) {
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

	if input.BlobContentLengthBytes == 0 || input.BlobContentLengthBytes%512 != 0 {
		err = fmt.Errorf("`input.BlobContentLengthBytes` must be aligned to a 512-byte boundary")
		return
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: putPageBlobOptions{
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

type putPageBlobOptions struct {
	input PutPageBlobInput
}

func (p putPageBlobOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

	headers.Append("x-ms-blob-type", string(PageBlob))

	// For a page blob or an page blob, the value of this header must be set to zero,
	// as Put Blob is used only to initialize the blob
	headers.Append("Content-Length", "0")

	// This header specifies the maximum size for the page blob, up to 8 TB.
	// The page blob size must be aligned to a 512-byte boundary.
	headers.Append("x-ms-blob-content-length", strconv.Itoa(int(p.input.BlobContentLengthBytes)))

	if p.input.AccessTier != nil {
		headers.Append("x-ms-access-tier", string(*p.input.AccessTier))
	}

	if p.input.BlobSequenceNumber != nil {
		headers.Append("x-ms-blob-sequence-number", strconv.Itoa(int(*p.input.BlobSequenceNumber)))
	}

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

	headers.Merge(metadata.SetMetaDataHeaders(p.input.MetaData))
	return headers
}

func (p putPageBlobOptions) ToOData() *odata.Query {
	return nil
}

func (p putPageBlobOptions) ToQuery() *client.QueryParams {
	return nil
}
