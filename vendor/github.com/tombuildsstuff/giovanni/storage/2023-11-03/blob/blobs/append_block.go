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

type AppendBlockInput struct {

	// A number indicating the byte offset to compare.
	// Append Block will succeed only if the append position is equal to this number.
	// If it is not, the request will fail with an AppendPositionConditionNotMet
	// error (HTTP status code 412 – Precondition Failed)
	BlobConditionAppendPosition *int64

	// The max length in bytes permitted for the append blob.
	// If the Append Block operation would cause the blob to exceed that limit or if the blob size
	// is already greater than the value specified in this header, the request will fail with
	// an MaxBlobSizeConditionNotMet error (HTTP status code 412 – Precondition Failed).
	BlobConditionMaxSize *int64

	// The Bytes which should be appended to the end of this Append Blob.
	// This can either be nil, which creates an empty blob, or a byte array
	Content *[]byte

	// An MD5 hash of the block content.
	// This hash is used to verify the integrity of the block during transport.
	// When this header is specified, the storage service compares the hash of the content
	// that has arrived with this header value.
	//
	// Note that this MD5 hash is not stored with the blob.
	// If the two hashes do not match, the operation will fail with error code 400 (Bad Request).
	ContentMD5 *string

	// Required if the blob has an active lease.
	// To perform this operation on a blob with an active lease, specify the valid lease ID for this header.
	LeaseID *string
}

type AppendBlockResponse struct {
	HttpResponse *client.Response

	BlobAppendOffset        string
	BlobCommittedBlockCount int64
	ContentMD5              string
	ETag                    string
	LastModified            string
}

// AppendBlock commits a new block of data to the end of an existing append blob.
func (c Client) AppendBlock(ctx context.Context, containerName, blobName string, input AppendBlockInput) (resp AppendBlockResponse, err error) {

	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return resp, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return resp, fmt.Errorf("`blobName` cannot be an empty string")
	}

	if input.Content != nil && len(*input.Content) > (4*1024*1024) {
		return resp, fmt.Errorf("`input.Content` must be at most 4MB")
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: appendBlockOptions{
			input: input,
		},
		Path: fmt.Sprintf("/%s/%s", containerName, blobName),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	if input.Content != nil {
		req.Body = io.NopCloser(bytes.NewReader(*input.Content))
	}

	req.ContentLength = int64(len(*input.Content))

	resp.HttpResponse, err = req.Execute(ctx)
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	if resp.HttpResponse != nil {
		if resp.HttpResponse.Header != nil {
			resp.BlobAppendOffset = resp.HttpResponse.Header.Get("x-ms-blob-append-offset")
			resp.ContentMD5 = resp.HttpResponse.Header.Get("Content-MD5")
			resp.ETag = resp.HttpResponse.Header.Get("ETag")
			resp.LastModified = resp.HttpResponse.Header.Get("Last-Modified")

			if v := resp.HttpResponse.Header.Get("x-ms-blob-committed-block-count"); v != "" {
				i, innerErr := strconv.Atoi(v)
				if innerErr != nil {
					err = fmt.Errorf("error parsing %q as an integer: %s", v, innerErr)
					return
				}
				resp.BlobCommittedBlockCount = int64(i)
			}

		}
	}

	return
}

type appendBlockOptions struct {
	input AppendBlockInput
}

func (a appendBlockOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	if a.input.BlobConditionAppendPosition != nil {
		headers.Append("x-ms-blob-condition-appendpos", strconv.Itoa(int(*a.input.BlobConditionAppendPosition)))
	}
	if a.input.BlobConditionMaxSize != nil {
		headers.Append("x-ms-blob-condition-maxsize", strconv.Itoa(int(*a.input.BlobConditionMaxSize)))
	}
	if a.input.ContentMD5 != nil {
		headers.Append("x-ms-blob-content-md5", *a.input.ContentMD5)
	}
	if a.input.LeaseID != nil {
		headers.Append("x-ms-lease-id", *a.input.LeaseID)
	}
	if a.input.Content != nil {
		headers.Append("Content-Length", strconv.Itoa(len(*a.input.Content)))
	}
	return headers
}

func (a appendBlockOptions) ToOData() *odata.Query {
	return nil
}

func (a appendBlockOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "appendblock")
	return out
}
