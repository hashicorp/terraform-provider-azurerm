package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/tombuildsstuff/giovanni/storage/internal/metadata"
)

type SetMetaDataInput struct {
	// The ID of the Lease
	// This must be specified if a Lease is present on the Blob, else a 403 is returned
	LeaseID *string

	// Any metadata which should be added to this blob
	MetaData map[string]string
}

type SetMetaDataResponse struct {
	HttpResponse *client.Response
}

// SetMetaData marks the specified blob or snapshot for deletion. The blob is later deleted during garbage collection.
func (c Client) SetMetaData(ctx context.Context, containerName, blobName string, input SetMetaDataInput) (resp SetMetaDataResponse, err error) {

	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return resp, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return resp, fmt.Errorf("`blobName` cannot be an empty string")
	}

	if err := metadata.Validate(input.MetaData); err != nil {
		return resp, fmt.Errorf(fmt.Sprintf("`input.MetaData` is not valid: %s.", err))
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: setMetadataOptions{
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

type setMetadataOptions struct {
	input SetMetaDataInput
}

func (s setMetadataOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	if s.input.LeaseID != nil {
		headers.Append("x-ms-lease-id", *s.input.LeaseID)
	}
	headers.Merge(metadata.SetMetaDataHeaders(s.input.MetaData))
	return headers
}

func (s setMetadataOptions) ToOData() *odata.Query {
	return nil
}

func (s setMetadataOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "metadata")
	return out
}
