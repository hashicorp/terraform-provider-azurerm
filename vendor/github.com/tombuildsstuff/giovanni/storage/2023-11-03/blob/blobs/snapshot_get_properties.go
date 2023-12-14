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

type GetSnapshotPropertiesInput struct {
	// The ID of the Lease
	// This must be specified if a Lease is present on the Blob, else a 403 is returned
	LeaseID *string

	// The ID of the Snapshot which should be retrieved
	SnapshotID string
}

// GetSnapshotProperties returns all user-defined metadata, standard HTTP properties, and system properties for
// the specified snapshot of a blob
func (c Client) GetSnapshotProperties(ctx context.Context, containerName, blobName string, input GetSnapshotPropertiesInput) (resp GetPropertiesResponse, err error) {

	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return resp, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return resp, fmt.Errorf("`blobName` cannot be an empty string")
	}

	if input.SnapshotID == "" {
		return resp, fmt.Errorf("`input.SnapshotID` cannot be an empty string")
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

	resp.HttpResponse, err = req.Execute(ctx)
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	if resp.HttpResponse != nil {
		if resp.HttpResponse.Header != nil {
			resp.AccessTier = AccessTier(resp.HttpResponse.Header.Get("x-ms-access-tier"))
			resp.AccessTierChangeTime = resp.HttpResponse.Header.Get("x-ms-access-tier-change-time")
			resp.ArchiveStatus = ArchiveStatus(resp.HttpResponse.Header.Get("x-ms-archive-status"))
			resp.BlobCommittedBlockCount = resp.HttpResponse.Header.Get("x-ms-blob-committed-block-count")
			resp.BlobSequenceNumber = resp.HttpResponse.Header.Get("x-ms-blob-sequence-number")
			resp.BlobType = BlobType(resp.HttpResponse.Header.Get("x-ms-blob-type"))
			resp.CacheControl = resp.HttpResponse.Header.Get("Cache-Control")
			resp.ContentDisposition = resp.HttpResponse.Header.Get("Content-Disposition")
			resp.ContentEncoding = resp.HttpResponse.Header.Get("Content-Encoding")
			resp.ContentLanguage = resp.HttpResponse.Header.Get("Content-Language")
			resp.ContentMD5 = resp.HttpResponse.Header.Get("Content-MD5")
			resp.ContentType = resp.HttpResponse.Header.Get("Content-Type")
			resp.CopyCompletionTime = resp.HttpResponse.Header.Get("x-ms-copy-completion-time")
			resp.CopyDestinationSnapshot = resp.HttpResponse.Header.Get("x-ms-copy-destination-snapshot")
			resp.CopyID = resp.HttpResponse.Header.Get("x-ms-copy-id")
			resp.CopyProgress = resp.HttpResponse.Header.Get("x-ms-copy-progress")
			resp.CopySource = resp.HttpResponse.Header.Get("x-ms-copy-source")
			resp.CopyStatus = CopyStatus(resp.HttpResponse.Header.Get("x-ms-copy-status"))
			resp.CopyStatusDescription = resp.HttpResponse.Header.Get("x-ms-copy-status-description")
			resp.CreationTime = resp.HttpResponse.Header.Get("x-ms-creation-time")
			resp.ETag = resp.HttpResponse.Header.Get("Etag")
			resp.LastModified = resp.HttpResponse.Header.Get("Last-Modified")
			resp.LeaseDuration = LeaseDuration(resp.HttpResponse.Header.Get("x-ms-lease-duration"))
			resp.LeaseState = LeaseState(resp.HttpResponse.Header.Get("x-ms-lease-state"))
			resp.LeaseStatus = LeaseStatus(resp.HttpResponse.Header.Get("x-ms-lease-status"))
			resp.MetaData = metadata.ParseFromHeaders(resp.HttpResponse.Header)

			if v := resp.HttpResponse.Header.Get("x-ms-access-tier-inferred"); v != "" {
				b, innerErr := strconv.ParseBool(v)
				if innerErr != nil {
					err = fmt.Errorf("error parsing %q as a bool: %s", v, innerErr)
					return
				}

				resp.AccessTierInferred = b
			}

			if v := resp.HttpResponse.Header.Get("Content-Length"); v != "" {
				i, innerErr := strconv.Atoi(v)
				if innerErr != nil {
					err = fmt.Errorf("error parsing %q as an integer: %s", v, innerErr)
				}

				resp.ContentLength = int64(i)
			}

			if v := resp.HttpResponse.Header.Get("x-ms-incremental-copy"); v != "" {
				b, innerErr := strconv.ParseBool(v)
				if innerErr != nil {
					err = fmt.Errorf("error parsing %q as a bool: %s", v, innerErr)
					return
				}

				resp.IncrementalCopy = b
			}

			if v := resp.HttpResponse.Header.Get("x-ms-server-encrypted"); v != "" {
				b, innerErr := strconv.ParseBool(v)
				if innerErr != nil {
					err = fmt.Errorf("error parsing %q as a bool: %s", v, innerErr)
					return
				}

				resp.ServerEncrypted = b
			}
		}
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
