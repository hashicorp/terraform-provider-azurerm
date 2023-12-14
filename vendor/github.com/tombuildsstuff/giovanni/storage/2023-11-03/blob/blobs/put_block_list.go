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

type BlockList struct {
	CommittedBlockIDs   []BlockID `xml:"Committed,omitempty"`
	UncommittedBlockIDs []BlockID `xml:"Uncommitted,omitempty"`
	LatestBlockIDs      []BlockID `xml:"Latest,omitempty"`
}

type BlockID struct {
	Value string `xml:",chardata"`
}

type PutBlockListInput struct {
	BlockList          BlockList
	CacheControl       *string
	ContentDisposition *string
	ContentEncoding    *string
	ContentLanguage    *string
	ContentMD5         *string
	ContentType        *string
	MetaData           map[string]string
	LeaseID            *string
}

type PutBlockListResponse struct {
	HttpResponse *client.Response

	ContentMD5   string
	ETag         string
	LastModified string
}

// PutBlockList writes a blob by specifying the list of block IDs that make up the blob.
// In order to be written as part of a blob, a block must have been successfully written
// to the server in a prior Put Block operation.
func (c Client) PutBlockList(ctx context.Context, containerName, blobName string, input PutBlockListInput) (resp PutBlockListResponse, err error) {

	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return resp, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return resp, fmt.Errorf("`blobName` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: putBlockListOptions{
			input: input,
		},
		Path: fmt.Sprintf("/%s/%s", containerName, blobName),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	err = req.Marshal(&input.BlockList)
	if err != nil {
		return resp, fmt.Errorf("marshalling request: %v", err)
	}

	resp.HttpResponse, err = req.Execute(ctx)
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	if resp.HttpResponse != nil {
		if resp.HttpResponse.Header != nil {
			resp.ContentMD5 = resp.HttpResponse.Header.Get("Content-MD5")
			resp.ETag = resp.HttpResponse.Header.Get("ETag")
			resp.LastModified = resp.HttpResponse.Header.Get("Last-Modified")
		}
	}

	return
}

type putBlockListOptions struct {
	input PutBlockListInput
}

func (p putBlockListOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

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

func (p putBlockListOptions) ToOData() *odata.Query {
	return nil
}

func (p putBlockListOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "blocklist")
	return out
}
