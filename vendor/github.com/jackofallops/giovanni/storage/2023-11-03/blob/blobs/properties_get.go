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

type GetPropertiesInput struct {
	// The ID of the Lease
	// This must be specified if a Lease is present on the Blob, else a 403 is returned
	LeaseID *string
}

type GetPropertiesResponse struct {
	HttpResponse *http.Response

	// The tier of page blob on a premium storage account or tier of block blob on blob storage or general purpose v2 account.
	AccessTier AccessTier

	// This gives the last time tier was changed on the object.
	// This header is returned only if tier on block blob was ever set.
	// The date format follows RFC 1123
	AccessTierChangeTime string

	// For page blobs on a premium storage account only.
	// If the access tier is not explicitly set on the blob, the tier is inferred based on its content length
	// and this header will be returned with true value.
	// For block blobs on Blob Storage or general purpose v2 account, if the blob does not have the access tier
	// set then we infer the tier from the storage account properties. This header is set only if the block blob
	// tier is inferred
	AccessTierInferred bool

	// For blob storage or general purpose v2 account.
	// If the blob is being rehydrated and is not complete then this header is returned indicating
	// that rehydrate is pending and also tells the destination tier
	ArchiveStatus ArchiveStatus

	// The number of committed blocks present in the blob.
	// This header is returned only for append blobs.
	BlobCommittedBlockCount string

	// The current sequence number for a page blob.
	// This header is not returned for block blobs or append blobs.
	// This header is not returned for block blobs.
	BlobSequenceNumber string

	// The blob type.
	BlobType BlobType

	// If the Cache-Control request header has previously been set for the blob, that value is returned in this header.
	CacheControl string

	// The Content-Disposition response header field conveys additional information about how to process
	// the response payload, and also can be used to attach additional metadata.
	// For example, if set to attachment, it indicates that the user-agent should not display the response,
	// but instead show a Save As dialog.
	ContentDisposition string

	// If the Content-Encoding request header has previously been set for the blob,
	// that value is returned in this header.
	ContentEncoding string

	// If the Content-Language request header has previously been set for the blob,
	// that value is returned in this header.
	ContentLanguage string

	// The size of the blob in bytes.
	// For a page blob, this header returns the value of the x-ms-blob-content-length header stored with the blob.
	ContentLength int64

	// The content type specified for the blob.
	// If no content type was specified, the default content type is `application/octet-stream`.
	ContentType string

	// If the Content-MD5 header has been set for the blob, this response header is returned so that
	// the client can check for message content integrity.
	ContentMD5 string

	// Conclusion time of the last attempted Copy Blob operation where this blob was the destination blob.
	// This value can specify the time of a completed, aborted, or failed copy attempt.
	// This header does not appear if a copy is pending, if this blob has never been the
	// destination in a Copy Blob operation, or if this blob has been modified after a concluded Copy Blob
	// operation using Set Blob Properties, Put Blob, or Put Block List.
	CopyCompletionTime string

	// Included if the blob is incremental copy blob or incremental copy snapshot, if x-ms-copy-status is success.
	// Snapshot time of the last successful incremental copy snapshot for this blob
	CopyDestinationSnapshot string

	// String identifier for the last attempted Copy Blob operation where this blob was the destination blob.
	// This header does not appear if this blob has never been the destination in a Copy Blob operation,
	// or if this blob has been modified after a concluded Copy Blob operation using Set Blob Properties,
	// Put Blob, or Put Block List.
	CopyID string

	// Contains the number of bytes copied and the total bytes in the source in the last attempted
	// Copy Blob operation where this blob was the destination blob.
	// Can show between 0 and Content-Length bytes copied.
	// This header does not appear if this blob has never been the destination in a Copy Blob operation,
	// or if this blob has been modified after a concluded Copy Blob operation using Set Blob Properties,
	// Put Blob, or Put Block List.
	CopyProgress string

	// URL up to 2 KB in length that specifies the source blob used in the last attempted Copy Blob operation
	// where this blob was the destination blob.
	// This header does not appear if this blob has never been the destination in a Copy Blob operation,
	// or if this blob has been modified after a concluded Copy Blob operation using Set Blob Properties,
	// Put Blob, or Put Block List
	CopySource string

	// State of the copy operation identified by x-ms-copy-id, with these values:
	// - success: Copy completed successfully.
	// - pending: Copy is in progress.
	//            Check x-ms-copy-status-description if intermittent, non-fatal errors
	//            impede copy progress but don’t cause failure.
	// - aborted: Copy was ended by Abort Copy Blob.
	// - failed: Copy failed. See x-ms- copy-status-description for failure details.
	// This header does not appear if this blob has never been the destination in a Copy Blob operation,
	// or if this blob has been modified after a completed Copy Blob operation using Set Blob Properties,
	// Put Blob, or Put Block List.
	CopyStatus CopyStatus

	// Describes cause of fatal or non-fatal copy operation failure.
	// This header does not appear if this blob has never been the destination in a Copy Blob operation,
	// or if this blob has been modified after a concluded Copy Blob operation using Set Blob Properties,
	// Put Blob, or Put Block List.
	CopyStatusDescription string

	// The date/time at which the blob was created. The date format follows RFC 1123
	CreationTime string

	// The ETag contains a value that you can use to perform operations conditionally
	ETag string

	// Included if the blob is incremental copy blob.
	IncrementalCopy bool

	// The date/time that the blob was last modified. The date format follows RFC 1123.
	LastModified string

	// When a blob is leased, specifies whether the lease is of infinite or fixed duration
	LeaseDuration LeaseDuration

	// The lease state of the blob
	LeaseState LeaseState

	LeaseStatus LeaseStatus

	// A set of name-value pairs that correspond to the user-defined metadata associated with this blob
	MetaData map[string]string

	// Is the Storage Account encrypted using server-side encryption? This should always return true
	ServerEncrypted bool

	// The encryption scope for the request content.
	EncryptionScope string
}

// GetProperties returns all user-defined metadata, standard HTTP properties, and system properties for the blob
func (c Client) GetProperties(ctx context.Context, containerName, blobName string, input GetPropertiesInput) (result GetPropertiesResponse, err error) {
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
			http.StatusOK,
		},
		HttpMethod: http.MethodHead,
		OptionsObject: getPropertiesOptions{
			leaseID: input.LeaseID,
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
				result.AccessTier = AccessTier(resp.Header.Get("x-ms-access-tier"))
				result.AccessTierChangeTime = resp.Header.Get("x-ms-access-tier-change-time")
				result.ArchiveStatus = ArchiveStatus(resp.Header.Get("x-ms-archive-status"))
				result.BlobCommittedBlockCount = resp.Header.Get("x-ms-blob-committed-block-count")
				result.BlobSequenceNumber = resp.Header.Get("x-ms-blob-sequence-number")
				result.BlobType = BlobType(resp.Header.Get("x-ms-blob-type"))
				result.CacheControl = resp.Header.Get("Cache-Control")
				result.ContentDisposition = resp.Header.Get("Content-Disposition")
				result.ContentEncoding = resp.Header.Get("Content-Encoding")
				result.ContentLanguage = resp.Header.Get("Content-Language")
				result.ContentMD5 = resp.Header.Get("Content-MD5")
				result.ContentType = resp.Header.Get("Content-Type")
				result.CopyCompletionTime = resp.Header.Get("x-ms-copy-completion-time")
				result.CopyDestinationSnapshot = resp.Header.Get("x-ms-copy-destination-snapshot")
				result.CopyID = resp.Header.Get("x-ms-copy-id")
				result.CopyProgress = resp.Header.Get("x-ms-copy-progress")
				result.CopySource = resp.Header.Get("x-ms-copy-source")
				result.CopyStatus = CopyStatus(resp.Header.Get("x-ms-copy-status"))
				result.CopyStatusDescription = resp.Header.Get("x-ms-copy-status-description")
				result.CreationTime = resp.Header.Get("x-ms-creation-time")
				result.ETag = resp.Header.Get("Etag")
				result.LastModified = resp.Header.Get("Last-Modified")
				result.LeaseDuration = LeaseDuration(resp.Header.Get("x-ms-lease-duration"))
				result.LeaseState = LeaseState(resp.Header.Get("x-ms-lease-state"))
				result.LeaseStatus = LeaseStatus(resp.Header.Get("x-ms-lease-status"))
				result.EncryptionScope = resp.Header.Get("x-ms-encryption-scope")
				result.MetaData = metadata.ParseFromHeaders(resp.Header)

				if v := resp.Header.Get("x-ms-access-tier-inferred"); v != "" {
					b, innerErr := strconv.ParseBool(v)
					if innerErr != nil {
						err = fmt.Errorf("parsing `x-ms-access-tier-inferred` header value %q: %s", v, innerErr)
						return
					}
					result.AccessTierInferred = b
				}

				if v := resp.Header.Get("Content-Length"); v != "" {
					i, innerErr := strconv.Atoi(v)
					if innerErr != nil {
						err = fmt.Errorf("parsing `Content-Length` header value %q: %s", v, innerErr)
					}
					result.ContentLength = int64(i)
				}

				if v := resp.Header.Get("x-ms-incremental-copy"); v != "" {
					b, innerErr := strconv.ParseBool(v)
					if innerErr != nil {
						err = fmt.Errorf("parsing `x-ms-incremental-copy` header value %q: %s", v, innerErr)
						return
					}
					result.IncrementalCopy = b
				}

				if v := resp.Header.Get("x-ms-server-encrypted"); v != "" {
					b, innerErr := strconv.ParseBool(v)
					if innerErr != nil {
						err = fmt.Errorf("parsing `x-ms-server-encrypted` header value %q: %s", v, innerErr)
						return
					}
					result.ServerEncrypted = b
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

type getPropertiesOptions struct {
	leaseID *string
}

func (g getPropertiesOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

	if g.leaseID != nil {
		headers.Append("x-ms-lease-id", *g.leaseID)
	}
	return headers
}

func (g getPropertiesOptions) ToOData() *odata.Query {
	return nil
}

func (g getPropertiesOptions) ToQuery() *client.QueryParams {
	return nil
}
