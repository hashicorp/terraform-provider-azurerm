package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type IncrementalCopyBlobInput struct {
	CopySource        string
	IfModifiedSince   *string
	IfUnmodifiedSince *string
	IfMatch           *string
	IfNoneMatch       *string
}

type IncrementalCopyBlob struct {
	HttpResponse *client.Response
}

// IncrementalCopyBlob copies a snapshot of the source page blob to a destination page blob.
// The snapshot is copied such that only the differential changes between the previously copied
// snapshot are transferred to the destination.
// The copied snapshots are complete copies of the original snapshot and can be read or copied from as usual.
func (c Client) IncrementalCopyBlob(ctx context.Context, containerName, blobName string, input IncrementalCopyBlobInput) (resp IncrementalCopyBlob, err error) {

	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return resp, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return resp, fmt.Errorf("`blobName` cannot be an empty string")
	}

	if input.CopySource == "" {
		return resp, fmt.Errorf("`input.CopySource` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: incrementalCopyBlobOptions{
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

	return
}

type incrementalCopyBlobOptions struct {
	input IncrementalCopyBlobInput
}

func (i incrementalCopyBlobOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

	headers.Append("x-ms-copy-source", i.input.CopySource)

	if i.input.IfModifiedSince != nil {
		headers.Append("If-Modified-Since", *i.input.IfModifiedSince)
	}
	if i.input.IfUnmodifiedSince != nil {
		headers.Append("If-Unmodified-Since", *i.input.IfUnmodifiedSince)
	}
	if i.input.IfMatch != nil {
		headers.Append("If-Match", *i.input.IfMatch)
	}
	if i.input.IfNoneMatch != nil {
		headers.Append("If-None-Match", *i.input.IfNoneMatch)
	}
	return headers
}

func (i incrementalCopyBlobOptions) ToOData() *odata.Query {
	return nil
}

func (i incrementalCopyBlobOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "incrementalcopy")
	return out
}
