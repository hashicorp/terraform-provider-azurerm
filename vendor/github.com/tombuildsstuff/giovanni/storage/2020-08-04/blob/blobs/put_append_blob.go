package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/tombuildsstuff/giovanni/storage/internal/metadata"
)

type PutAppendBlobInput struct {
	CacheControl       *string
	ContentDisposition *string
	ContentEncoding    *string
	ContentLanguage    *string
	ContentMD5         *string
	ContentType        *string
	LeaseID            *string
	MetaData           map[string]string
}

type PutAppendBlobResponse struct {
	HttpResponse *client.Response
}

// PutAppendBlob is a wrapper around the Put API call (with a stricter input object)
// which creates a new append blob, or updates the content of an existing blob.
func (c Client) PutAppendBlob(ctx context.Context, containerName, blobName string, input PutAppendBlobInput) (resp PutAppendBlobResponse, err error) {

	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return resp, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return resp, fmt.Errorf("`blobName` cannot be an empty string")
	}

	if err := metadata.Validate(input.MetaData); err != nil {
		return resp, fmt.Errorf(fmt.Sprintf("`input.MetaData` is not valid: %s.", err))
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

	resp.HttpResponse, err = req.Execute(ctx)
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

	headers.Merge(metadata.SetMetaDataHeaders(p.input.MetaData))
	return headers
}

func (p putAppendBlobOptions) ToOData() *odata.Query {
	return nil
}

func (p putAppendBlobOptions) ToQuery() *client.QueryParams {
	return nil
}
