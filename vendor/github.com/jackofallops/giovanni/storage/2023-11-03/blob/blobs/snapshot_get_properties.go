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

type GetSnapshotPropertiesInput struct {
	// The ID of the Lease
	// This must be specified if a Lease is present on the Blob, else a 403 is returned
	LeaseID *string

	// The ID of the Snapshot which should be retrieved
	SnapshotID string
}

// GetSnapshotProperties returns all user-defined metadata, standard HTTP properties, and system properties for
// the specified snapshot of a blob
func (c Client) GetSnapshotProperties(ctx context.Context, containerName, blobName string, input GetSnapshotPropertiesInput) (result GetPropertiesResponse, err error) {
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

	if input.SnapshotID == "" {
		err = fmt.Errorf("`input.SnapshotID` cannot be an empty string")
		return
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodHead,
		OptionsObject: snapshotGetPropertiesOptions{
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
				result.MetaData = metadata.ParseFromHeaders(resp.Header)

				if v := resp.Header.Get("Content-Length"); v != "" {
					i, innerErr := strconv.Atoi(v)
					if innerErr != nil {
						err = fmt.Errorf("parsing `Content-Length` header value %q: %+v", v, innerErr)
					}

					result.ContentLength = int64(i)
				}

				if v := resp.Header.Get("x-ms-access-tier-inferred"); v != "" {
					b, innerErr := strconv.ParseBool(v)
					if innerErr != nil {
						err = fmt.Errorf("parsing `x-ms-access-tier-inferred` header value %q: %+v", v, innerErr)
						return
					}

					result.AccessTierInferred = b
				}

				if v := resp.Header.Get("x-ms-incremental-copy"); v != "" {
					b, innerErr := strconv.ParseBool(v)
					if innerErr != nil {
						err = fmt.Errorf("parsing `x-ms-incremental-copy` header value %q: %+v", v, innerErr)
						return
					}

					result.IncrementalCopy = b
				}

				if v := resp.Header.Get("x-ms-server-encrypted"); v != "" {
					b, innerErr := strconv.ParseBool(v)
					if innerErr != nil {
						err = fmt.Errorf("parsing `\"x-ms-server-encrypted` header value %q: %+v", v, innerErr)
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

	if result.HttpResponse != nil {
	}

	return
}

type snapshotGetPropertiesOptions struct {
	input GetSnapshotPropertiesInput
}

func (s snapshotGetPropertiesOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	if s.input.LeaseID != nil {
		headers.Append("x-ms-lease-id", *s.input.LeaseID)
	}
	return headers
}

func (s snapshotGetPropertiesOptions) ToOData() *odata.Query {
	return nil
}

func (s snapshotGetPropertiesOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("snapshot", s.input.SnapshotID)
	return out
}
