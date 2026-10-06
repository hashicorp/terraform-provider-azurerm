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

type SnapshotInput struct {
	// The ID of the Lease
	// This must be specified if a Lease is present on the Blob, else a 403 is returned
	LeaseID *string

	// The encryption scope to set for the request content.
	EncryptionScope *string

	// MetaData is a user-defined name-value pair associated with the blob.
	// If no name-value pairs are specified, the operation will copy the base blob metadata to the snapshot.
	// If one or more name-value pairs are specified, the snapshot is created with the specified metadata,
	// and metadata is not copied from the base blob.
	MetaData map[string]string

	// A DateTime value which will only snapshot the blob if it has been modified since the specified date/time
	// If the base blob has not been modified, the Blob service returns status code 412 (Precondition Failed).
	IfModifiedSince *string

	// A DateTime value which will only snapshot the blob if it has not been modified since the specified date/time
	// If the base blob has been modified, the Blob service returns status code 412 (Precondition Failed).
	IfUnmodifiedSince *string

	// An ETag value to snapshot the blob only if its ETag value matches the value specified.
	// If the values do not match, the Blob service returns status code 412 (Precondition Failed).
	IfMatch *string

	// An ETag value for this conditional header to snapshot the blob only if its ETag value
	// does not match the value specified.
	// If the values are identical, the Blob service returns status code 412 (Precondition Failed).
	IfNoneMatch *string
}

type SnapshotResponse struct {
	HttpResponse *http.Response

	// The ETag of the snapshot
	ETag string

	// A DateTime value that uniquely identifies the snapshot.
	// The value of this header indicates the snapshot version,
	// and may be used in subsequent requests to access the snapshot.
	SnapshotDateTime string
}

// Snapshot captures a Snapshot of a given Blob
func (c Client) Snapshot(ctx context.Context, containerName, blobName string, input SnapshotInput) (result SnapshotResponse, err error) {
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

	if err = metadata.Validate(input.MetaData); err != nil {
		err = fmt.Errorf(fmt.Sprintf("`input.MetaData` is not valid: %s.", err))
		return
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: snapshotOptions{
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
				result.ETag = resp.Header.Get("ETag")
				result.SnapshotDateTime = resp.Header.Get("x-ms-snapshot")
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

type snapshotOptions struct {
	input SnapshotInput
}

func (s snapshotOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

	if s.input.LeaseID != nil {
		headers.Append("x-ms-lease-id", *s.input.LeaseID)
	}

	if s.input.EncryptionScope != nil {
		headers.Append("x-ms-encryption-scope", *s.input.EncryptionScope)
	}

	if s.input.IfModifiedSince != nil {
		headers.Append("If-Modified-Since", *s.input.IfModifiedSince)
	}

	if s.input.IfUnmodifiedSince != nil {
		headers.Append("If-Unmodified-Since", *s.input.IfUnmodifiedSince)
	}

	if s.input.IfMatch != nil {
		headers.Append("If-Match", *s.input.IfMatch)
	}

	if s.input.IfNoneMatch != nil {
		headers.Append("If-None-Match", *s.input.IfNoneMatch)
	}

	headers.Merge(metadata.SetMetaDataHeaders(s.input.MetaData))
	return headers
}

func (s snapshotOptions) ToOData() *odata.Query {
	return nil
}

func (s snapshotOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "snapshot")
	return out
}
