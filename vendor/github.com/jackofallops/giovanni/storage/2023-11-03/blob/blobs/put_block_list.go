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
	LeaseID            *string
	EncryptionScope    *string
	MetaData           map[string]string
}

type PutBlockListResponse struct {
	HttpResponse *http.Response

	ContentMD5   string
	ETag         string
	LastModified string
}

// PutBlockList writes a blob by specifying the list of block IDs that make up the blob.
// In order to be written as part of a blob, a block must have been successfully written
// to the server in a prior Put Block operation.
func (c Client) PutBlockList(ctx context.Context, containerName, blobName string, input PutBlockListInput) (result PutBlockListResponse, err error) {
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
		err = fmt.Errorf("marshalling request: %+v", err)
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil && resp.Response != nil {
		result.HttpResponse = resp.Response

		if err == nil {
			if resp.Header != nil {
				result.ContentMD5 = resp.Header.Get("Content-MD5")
				result.ETag = resp.Header.Get("ETag")
				result.LastModified = resp.Header.Get("Last-Modified")
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
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
	if p.input.EncryptionScope != nil {
		headers.Append("x-ms-encryption-scope", *p.input.EncryptionScope)
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
