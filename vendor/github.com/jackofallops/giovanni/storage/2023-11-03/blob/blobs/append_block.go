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

	// The encryption scope to set for the request content.
	EncryptionScope *string
}

type AppendBlockResponse struct {
	HttpResponse *http.Response

	BlobAppendOffset        string
	BlobCommittedBlockCount int64
	ContentMD5              string
	ETag                    string
	LastModified            string
}

// AppendBlock commits a new block of data to the end of an existing append blob.
func (c Client) AppendBlock(ctx context.Context, containerName, blobName string, input AppendBlockInput) (result AppendBlockResponse, err error) {
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

	if input.Content != nil && len(*input.Content) > (4*1024*1024) {
		err = fmt.Errorf("`input.Content` must be at most 4MB")
		return
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

	err = req.Marshal(input.Content)
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
				result.BlobAppendOffset = resp.Header.Get("x-ms-blob-append-offset")
				result.ContentMD5 = resp.Header.Get("Content-MD5")
				result.ETag = resp.Header.Get("ETag")
				result.LastModified = resp.Header.Get("Last-Modified")

				if v := resp.Header.Get("x-ms-blob-committed-block-count"); v != "" {
					i, innerErr := strconv.Atoi(v)
					if innerErr != nil {
						err = fmt.Errorf("parsing `x-ms-blob-committed-block-count` header value %q: %+v", v, innerErr)
						return
					}
					result.BlobCommittedBlockCount = int64(i)
				}
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
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
	if a.input.EncryptionScope != nil {
		headers.Append("x-ms-encryption-scope", *a.input.EncryptionScope)
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
