package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type GetBlockListInput struct {
	BlockListType BlockListType
	LeaseID       *string
}

type GetBlockListResponse struct {
	HttpResponse *http.Response

	// The size of the blob in bytes
	BlobContentLength *int64

	// The Content Type of the blob
	ContentType string

	// The ETag associated with this blob
	ETag string

	// A list of blocks which have been committed
	CommittedBlocks CommittedBlocks `xml:"CommittedBlocks,omitempty"`

	// A list of blocks which have not yet been committed
	UncommittedBlocks UncommittedBlocks `xml:"UncommittedBlocks,omitempty"`
}

// GetBlockList retrieves the list of blocks that have been uploaded as part of a block blob.
func (c Client) GetBlockList(ctx context.Context, containerName, blobName string, input GetBlockListInput) (result GetBlockListResponse, err error) {
	if containerName == "" {
		return result, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return result, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return result, fmt.Errorf("`blobName` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		OptionsObject: getBlockListOptions{
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
				result.ContentType = resp.Header.Get("Content-Type")
				result.ETag = resp.Header.Get("ETag")

				if v := resp.Header.Get("x-ms-blob-content-length"); v != "" {
					i, innerErr := strconv.Atoi(v)
					if innerErr != nil {
						err = fmt.Errorf("parsing `x-ms-blob-content-length` header value %q: %s", v, innerErr)
						return
					}

					i64 := int64(i)
					result.BlobContentLength = &i64
				}
			}

			err = resp.Unmarshal(&result)
			if err != nil {
				err = fmt.Errorf("unmarshalling response: %+v", err)
				return
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

type getBlockListOptions struct {
	input GetBlockListInput
}

func (g getBlockListOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	if g.input.LeaseID != nil {
		headers.Append("x-ms-lease-id", *g.input.LeaseID)
	}
	return headers
}

func (g getBlockListOptions) ToOData() *odata.Query {
	return nil
}

func (g getBlockListOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("blocklisttype", string(g.input.BlockListType))
	out.Append("comp", "blocklist")
	return out
}
