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

type CopyInput struct {
	// Specifies the name of the source blob or file.
	// Beginning with version 2012-02-12, this value may be a URL of up to 2 KB in length that specifies a blob.
	// The value should be URL-encoded as it would appear in a request URI.
	// A source blob in the same storage account can be authenticated via Shared Key.
	// However, if the source is a blob in another account,
	// the source blob must either be public or must be authenticated via a shared access signature.
	// If the source blob is public, no authentication is required to perform the copy operation.
	//
	// Beginning with version 2015-02-21, the source object may be a file in the Azure File service.
	// If the source object is a file that is to be copied to a blob, then the source file must be authenticated
	// using a shared access signature, whether it resides in the same account or in a different account.
	//
	// Only storage accounts created on or after June 7th, 2012 allow the Copy Blob operation to
	// copy from another storage account.
	CopySource string

	// The ID of the Lease
	// Required if the destination blob has an active lease.
	// The lease ID specified for this header must match the lease ID of the destination blob.
	// If the request does not include the lease ID or it is not valid,
	// the operation fails with status code 412 (Precondition Failed).
	//
	// If this header is specified and the destination blob does not currently have an active lease,
	// the operation will also fail with status code 412 (Precondition Failed).
	LeaseID *string

	// The ID of the Lease on the Source Blob
	// Specify to perform the Copy Blob operation only if the lease ID matches the active lease ID of the source blob.
	SourceLeaseID *string

	// For page blobs on a premium account only. Specifies the tier to be set on the target blob
	AccessTier *AccessTier

	// A user-defined name-value pair associated with the blob.
	// If no name-value pairs are specified, the operation will copy the metadata from the source blob or
	// file to the destination blob.
	// If one or more name-value pairs are specified, the destination blob is created with the specified metadata,
	// and metadata is not copied from the source blob or file.
	MetaData map[string]string

	// An ETag value.
	// Specify an ETag value for this conditional header to copy the blob only if the specified
	// ETag value matches the ETag value for an existing destination blob.
	// If the ETag for the destination blob does not match the ETag specified for If-Match,
	// the Blob service returns status code 412 (Precondition Failed).
	IfMatch *string

	// An ETag value, or the wildcard character (*).
	// Specify an ETag value for this conditional header to copy the blob only if the specified
	// ETag value does not match the ETag value for the destination blob.
	// Specify the wildcard character (*) to perform the operation only if the destination blob does not exist.
	// If the specified condition isn't met, the Blob service returns status code 412 (Precondition Failed).
	IfNoneMatch *string

	// A DateTime value.
	// Specify this conditional header to copy the blob only if the destination blob
	// has been modified since the specified date/time.
	// If the destination blob has not been modified, the Blob service returns status code 412 (Precondition Failed).
	IfModifiedSince *string

	// A DateTime value.
	// Specify this conditional header to copy the blob only if the destination blob
	// has not been modified since the specified date/time.
	// If the destination blob has been modified, the Blob service returns status code 412 (Precondition Failed).
	IfUnmodifiedSince *string

	// An ETag value.
	// Specify this conditional header to copy the source blob only if its ETag matches the value specified.
	// If the ETag values do not match, the Blob service returns status code 412 (Precondition Failed).
	// This cannot be specified if the source is an Azure File.
	SourceIfMatch *string

	// An ETag value.
	// Specify this conditional header to copy the blob only if its ETag does not match the value specified.
	// If the values are identical, the Blob service returns status code 412 (Precondition Failed).
	// This cannot be specified if the source is an Azure File.
	SourceIfNoneMatch *string

	// A DateTime value.
	// Specify this conditional header to copy the blob only if the source blob has been modified
	// since the specified date/time.
	// If the source blob has not been modified, the Blob service returns status code 412 (Precondition Failed).
	// This cannot be specified if the source is an Azure File.
	SourceIfModifiedSince *string

	// A DateTime value.
	// Specify this conditional header to copy the blob only if the source blob has not been modified
	// since the specified date/time.
	// If the source blob has been modified, the Blob service returns status code 412 (Precondition Failed).
	// This header cannot be specified if the source is an Azure File.
	SourceIfUnmodifiedSince *string

	// The encryption scope to set for the request content.
	EncryptionScope *string
}

type CopyResponse struct {
	HttpResponse *http.Response

	CopyID     string
	CopyStatus string
}

// Copy copies a blob to a destination within the storage account asynchronously.
func (c Client) Copy(ctx context.Context, containerName, blobName string, input CopyInput) (result CopyResponse, err error) {
	if containerName == "" {
		return result, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return result, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return result, fmt.Errorf("`blobName` cannot be an empty string")
	}

	if input.CopySource == "" {
		return result, fmt.Errorf("`input.CopySource` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: copyOptions{
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
				result.CopyID = resp.Header.Get("x-ms-copy-id")
				result.CopyStatus = resp.Header.Get("x-ms-copy-status")
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

type copyOptions struct {
	input CopyInput
}

func (c copyOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	headers.Append("x-ms-copy-source", c.input.CopySource)

	if c.input.LeaseID != nil {
		headers.Append("x-ms-lease-id", *c.input.LeaseID)
	}

	if c.input.SourceLeaseID != nil {
		headers.Append("x-ms-source-lease-id", *c.input.SourceLeaseID)
	}

	if c.input.AccessTier != nil {
		headers.Append("x-ms-access-tier", string(*c.input.AccessTier))
	}

	if c.input.IfMatch != nil {
		headers.Append("If-Match", *c.input.IfMatch)
	}

	if c.input.IfNoneMatch != nil {
		headers.Append("If-None-Match", *c.input.IfNoneMatch)
	}

	if c.input.IfUnmodifiedSince != nil {
		headers.Append("If-Unmodified-Since", *c.input.IfUnmodifiedSince)
	}

	if c.input.IfModifiedSince != nil {
		headers.Append("If-Modified-Since", *c.input.IfModifiedSince)
	}

	if c.input.SourceIfMatch != nil {
		headers.Append("x-ms-source-if-match", *c.input.SourceIfMatch)
	}

	if c.input.SourceIfNoneMatch != nil {
		headers.Append("x-ms-source-if-none-match", *c.input.SourceIfNoneMatch)
	}

	if c.input.SourceIfModifiedSince != nil {
		headers.Append("x-ms-source-if-modified-since", *c.input.SourceIfModifiedSince)
	}

	if c.input.SourceIfUnmodifiedSince != nil {
		headers.Append("x-ms-source-if-unmodified-since", *c.input.SourceIfUnmodifiedSince)
	}

	if c.input.EncryptionScope != nil {
		headers.Append("x-ms-encryption-scope", *c.input.EncryptionScope)
	}

	headers.Merge(metadata.SetMetaDataHeaders(c.input.MetaData))

	return headers
}

func (c copyOptions) ToOData() *odata.Query {
	return nil
}

func (c copyOptions) ToQuery() *client.QueryParams {
	return nil
}
