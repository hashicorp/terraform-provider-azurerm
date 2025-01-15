package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/jackofallops/giovanni/storage/internal/metadata"
)

type PutAppendBlobInput struct {
	CacheControl       *string
	ContentDisposition *string
	ContentEncoding    *string
	ContentLanguage    *string
	ContentMD5         *string
	ContentType        *string
	LeaseID            *string
	EncryptionScope    *string
	MetaData           map[string]string
}

type PutAppendBlobResponse struct {
	HttpResponse *http.Response
}

// PutAppendBlob is a wrapper around the Put API call (with a stricter input object)
// which creates a new append blob, or updates the content of an existing blob.
func (c Client) PutAppendBlob(ctx context.Context, containerName, blobName string, input PutAppendBlobInput) (result PutAppendBlobResponse, err error) {
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

	if err = metadata.Validate(input.MetaData); err != nil {
		err = fmt.Errorf(fmt.Sprintf("`input.MetaData` is not valid: %s.", err))
		return
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: putAppendBlobOptions{
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

type putAppendBlobOptions struct {
	input PutAppendBlobInput
}

func (p putAppendBlobOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

	headers.Append("x-ms-blob-type", string(AppendBlob))
	headers.Append("Content-Length", "0")

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

func (p putAppendBlobOptions) ToOData() *odata.Query {
	return nil
}

func (p putAppendBlobOptions) ToQuery() *client.QueryParams {
	return nil
}
