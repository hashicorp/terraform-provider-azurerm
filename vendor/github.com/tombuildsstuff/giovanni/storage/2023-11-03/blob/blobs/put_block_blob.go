package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/tombuildsstuff/giovanni/storage/internal/metadata"
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
	MetaData           map[string]string
}

type PutBlockBlobResponse struct {
	HttpResponse *client.Response
}

// PutBlockBlob is a wrapper around the Put API call (with a stricter input object)
// which creates a new block append blob, or updates the content of an existing block blob.
func (c Client) PutBlockBlob(ctx context.Context, containerName, blobName string, input PutBlockBlobInput) (resp PutBlockBlobResponse, err error) {

	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return resp, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return resp, fmt.Errorf("`blobName` cannot be an empty string")
	}

	if input.Content != nil && len(*input.Content) == 0 {
		return resp, fmt.Errorf("`input.Content` must either be nil or not empty")
	}

	if err := metadata.Validate(input.MetaData); err != nil {
		return resp, fmt.Errorf(fmt.Sprintf("`input.MetaData` is not valid: %s.", err))
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
